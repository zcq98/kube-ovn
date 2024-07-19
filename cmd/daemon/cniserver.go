package daemon

import (
	"fmt"
	"net/http"
	"net/http/pprof"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/klog/v2"
	"k8s.io/sample-controller/pkg/signals"

	kubeovninformer "github.com/kubeovn/kube-ovn/pkg/client/informers/externalversions"
	"github.com/kubeovn/kube-ovn/pkg/daemon"
	"github.com/kubeovn/kube-ovn/pkg/ovs"
	"github.com/kubeovn/kube-ovn/pkg/server"
	"github.com/kubeovn/kube-ovn/pkg/util"
	"github.com/kubeovn/kube-ovn/versions"
)

const svcName = "kube-ovn-cni"

func CmdMain() {
	defer klog.Flush()

	daemon.InitMetrics()
	util.InitKlogMetrics()

	config := daemon.ParseFlags()
	klog.Infof(versions.String())

	ovs.UpdateOVSVsctlLimiter(config.OVSVsctlConcurrency)

	nicBridgeMappings, err := daemon.InitOVSBridges()
	if err != nil {
		util.LogFatalAndExit(err, "failed to initialize OVS bridges")
	}

	if err = config.Init(nicBridgeMappings); err != nil {
		util.LogFatalAndExit(err, "failed to initialize config")
	}

	if err := Retry(util.ChassisRetryMaxTimes, util.ChassisCniDaemonRetryInterval, initChassisAnno, config); err != nil {
		util.LogFatalAndExit(err, "failed to initialize ovn chassis annotation")
	}

	if err = daemon.InitMirror(config); err != nil {
		util.LogFatalAndExit(err, "failed to initialize ovs mirror")
	}
	klog.Info("init node gw")
	if err = daemon.InitNodeGateway(config); err != nil {
		util.LogFatalAndExit(err, "failed to initialize node gateway")
	}

	if err := initForOS(); err != nil {
		util.LogFatalAndExit(err, "failed to do the OS initialization")
	}

	stopCh := signals.SetupSignalHandler().Done()
	podInformerFactory := kubeinformers.NewSharedInformerFactoryWithOptions(config.KubeClient, 0,
		kubeinformers.WithTweakListOptions(func(listOption *v1.ListOptions) {
			listOption.FieldSelector = fmt.Sprintf("spec.nodeName=%s", config.NodeName)
			listOption.AllowWatchBookmarks = true
		}))
	nodeInformerFactory := kubeinformers.NewSharedInformerFactoryWithOptions(config.KubeClient, 0,
		kubeinformers.WithTweakListOptions(func(listOption *v1.ListOptions) {
			listOption.AllowWatchBookmarks = true
		}))
	kubeovnInformerFactory := kubeovninformer.NewSharedInformerFactoryWithOptions(config.KubeOvnClient, 0,
		kubeovninformer.WithTweakListOptions(func(listOption *v1.ListOptions) {
			listOption.AllowWatchBookmarks = true
		}))
	ctl, err := daemon.NewController(config, stopCh, podInformerFactory, nodeInformerFactory, kubeovnInformerFactory)
	if err != nil {
		util.LogFatalAndExit(err, "failed to create controller")
	}
	klog.Info("start daemon controller")
	go ctl.Run(stopCh)
	go daemon.RunServer(config, ctl)
	if err := mvCNIConf(config.CniConfDir, config.CniConfFile, config.CniConfName); err != nil {
		util.LogFatalAndExit(err, "failed to mv cni config file")
	}

	mux := http.NewServeMux()
	if config.EnableMetrics {
		mux.Handle("/metrics", promhttp.Handler())
	}
	if config.EnablePprof {
		mux.HandleFunc("/debug/pprof/", pprof.Index)
		mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
		mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	}

	addr := util.GetDefaultListenAddr()
	if config.EnableVerboseConnCheck {
		go func() {
			connListenaddr := util.JoinHostPort(addr, config.TCPConnCheckPort)
			if err := util.TCPConnectivityListen(connListenaddr); err != nil {
				util.LogFatalAndExit(err, "failed to start TCP listen on addr %s ", addr)
			}
		}()

		go func() {
			connListenaddr := util.JoinHostPort(addr, config.UDPConnCheckPort)
			if err := util.UDPConnectivityListen(connListenaddr); err != nil {
				util.LogFatalAndExit(err, "failed to start UDP listen on addr %s ", addr)
			}
		}()
	}

	listenAddr := util.JoinHostPort(addr, config.PprofPort)
	if !config.SecureServing {
		server := &http.Server{
			Addr:              listenAddr,
			ReadHeaderTimeout: 3 * time.Second,
			Handler:           mux,
		}
		util.LogFatalAndExit(server.ListenAndServe(), "failed to listen and server on %s", server.Addr)
	} else {
		ch, err := server.SecureServing(listenAddr, svcName, mux)
		if err != nil {
			util.LogFatalAndExit(err, "failed to serve on %s", listenAddr)
		}
		<-ch
	}
}

func mvCNIConf(configDir, configFile, confName string) error {
	// #nosec
	data, err := os.ReadFile(configFile)
	if err != nil {
		klog.Errorf("failed to read cni config file %s, %v", configFile, err)
		return err
	}

	cniConfPath := filepath.Join(configDir, confName)
	return os.WriteFile(cniConfPath, data, 0o644)
}

func Retry(attempts, sleep int, f func(configuration *daemon.Configuration) error, ctrl *daemon.Configuration) (err error) {
	for i := 0; ; i++ {
		err = f(ctrl)
		if err == nil {
			return
		}
		if i >= (attempts - 1) {
			break
		}
		time.Sleep(time.Duration(sleep) * time.Second)
	}
	return err
}

func initChassisAnno(cfg *daemon.Configuration) error {
	chassisID, err := os.ReadFile(util.ChassisLoc)
	if err != nil {
		klog.Errorf("read chassis file failed, %v", err)
		return err
	}

	chassesName := strings.TrimSpace(string(chassisID))
	if chassesName == "" {
		// not ready yet
		err = fmt.Errorf("chassis id is empty")
		klog.Error(err)
		return err
	}
	annotations := map[string]any{util.ChassisAnnotation: chassesName}
	if err = util.UpdateNodeAnnotations(cfg.KubeClient.CoreV1().Nodes(), cfg.NodeName, annotations); err != nil {
		klog.Errorf("failed to update chassis annotation of node %s: %v", cfg.NodeName, err)
		return err
	}

	return nil
}
