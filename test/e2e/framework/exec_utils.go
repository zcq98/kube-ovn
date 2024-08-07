package framework

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kubernetes/test/e2e/framework"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"

	"github.com/kubeovn/kube-ovn/pkg/util"
)

// ExecCommandInContainer executes a command in the specified container.
func ExecCommandInContainer(f *Framework, namespace, pod, container string, cmd ...string) (string, string, error) {
	framework.Logf("Executing command %q in container %s/%s/%s", cmd, namespace, pod, container)
	return util.ExecuteCommandInContainer(f.ClientSet, f.ClientConfig(), namespace, pod, container, cmd...)
}

// ExecShellInContainer executes the specified command on the pod's container.
func ExecShellInContainer(f *Framework, namespace, pod, container, cmd string) (string, string, error) {
	return ExecCommandInContainer(f, namespace, pod, container, "/bin/sh", "-c", cmd)
}

func execCommandInPod(ctx context.Context, f *Framework, namespace, pod string, cmd ...string) (string, string, error) {
	ginkgo.GinkgoHelper()

	p, err := f.PodClientNS(namespace).Get(ctx, pod, metav1.GetOptions{})
	framework.ExpectNoError(err, "failed to get pod %s/%s", namespace, pod)
	gomega.Expect(p.Spec.Containers).NotTo(gomega.BeEmpty())
	return ExecCommandInContainer(f, namespace, pod, p.Spec.Containers[0].Name, cmd...)
}

// ExecShellInPod executes the specified command on the pod.
func ExecShellInPod(ctx context.Context, f *Framework, namespace, pod, cmd string) (string, string, error) {
	ginkgo.GinkgoHelper()
	return execCommandInPod(ctx, f, namespace, pod, "/bin/sh", "-c", cmd)
}
