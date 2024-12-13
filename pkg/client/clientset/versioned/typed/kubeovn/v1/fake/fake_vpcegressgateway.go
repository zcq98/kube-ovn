/*
Copyright The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	v1 "github.com/kubeovn/kube-ovn/pkg/apis/kubeovn/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeVpcEgressGateways implements VpcEgressGatewayInterface
type FakeVpcEgressGateways struct {
	Fake *FakeKubeovnV1
	ns   string
}

var vpcegressgatewaysResource = v1.SchemeGroupVersion.WithResource("vpc-egress-gateways")

var vpcegressgatewaysKind = v1.SchemeGroupVersion.WithKind("VpcEgressGateway")

// Get takes name of the vpcEgressGateway, and returns the corresponding vpcEgressGateway object, and an error if there is any.
func (c *FakeVpcEgressGateways) Get(ctx context.Context, name string, options metav1.GetOptions) (result *v1.VpcEgressGateway, err error) {
	emptyResult := &v1.VpcEgressGateway{}
	obj, err := c.Fake.
		Invokes(testing.NewGetActionWithOptions(vpcegressgatewaysResource, c.ns, name, options), emptyResult)

	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1.VpcEgressGateway), err
}

// List takes label and field selectors, and returns the list of VpcEgressGateways that match those selectors.
func (c *FakeVpcEgressGateways) List(ctx context.Context, opts metav1.ListOptions) (result *v1.VpcEgressGatewayList, err error) {
	emptyResult := &v1.VpcEgressGatewayList{}
	obj, err := c.Fake.
		Invokes(testing.NewListActionWithOptions(vpcegressgatewaysResource, vpcegressgatewaysKind, c.ns, opts), emptyResult)

	if obj == nil {
		return emptyResult, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1.VpcEgressGatewayList{ListMeta: obj.(*v1.VpcEgressGatewayList).ListMeta}
	for _, item := range obj.(*v1.VpcEgressGatewayList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested vpcEgressGateways.
func (c *FakeVpcEgressGateways) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchActionWithOptions(vpcegressgatewaysResource, c.ns, opts))

}

// Create takes the representation of a vpcEgressGateway and creates it.  Returns the server's representation of the vpcEgressGateway, and an error, if there is any.
func (c *FakeVpcEgressGateways) Create(ctx context.Context, vpcEgressGateway *v1.VpcEgressGateway, opts metav1.CreateOptions) (result *v1.VpcEgressGateway, err error) {
	emptyResult := &v1.VpcEgressGateway{}
	obj, err := c.Fake.
		Invokes(testing.NewCreateActionWithOptions(vpcegressgatewaysResource, c.ns, vpcEgressGateway, opts), emptyResult)

	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1.VpcEgressGateway), err
}

// Update takes the representation of a vpcEgressGateway and updates it. Returns the server's representation of the vpcEgressGateway, and an error, if there is any.
func (c *FakeVpcEgressGateways) Update(ctx context.Context, vpcEgressGateway *v1.VpcEgressGateway, opts metav1.UpdateOptions) (result *v1.VpcEgressGateway, err error) {
	emptyResult := &v1.VpcEgressGateway{}
	obj, err := c.Fake.
		Invokes(testing.NewUpdateActionWithOptions(vpcegressgatewaysResource, c.ns, vpcEgressGateway, opts), emptyResult)

	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1.VpcEgressGateway), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeVpcEgressGateways) UpdateStatus(ctx context.Context, vpcEgressGateway *v1.VpcEgressGateway, opts metav1.UpdateOptions) (result *v1.VpcEgressGateway, err error) {
	emptyResult := &v1.VpcEgressGateway{}
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceActionWithOptions(vpcegressgatewaysResource, "status", c.ns, vpcEgressGateway, opts), emptyResult)

	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1.VpcEgressGateway), err
}

// Delete takes name of the vpcEgressGateway and deletes it. Returns an error if one occurs.
func (c *FakeVpcEgressGateways) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteActionWithOptions(vpcegressgatewaysResource, c.ns, name, opts), &v1.VpcEgressGateway{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeVpcEgressGateways) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	action := testing.NewDeleteCollectionActionWithOptions(vpcegressgatewaysResource, c.ns, opts, listOpts)

	_, err := c.Fake.Invokes(action, &v1.VpcEgressGatewayList{})
	return err
}

// Patch applies the patch and returns the patched vpcEgressGateway.
func (c *FakeVpcEgressGateways) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.VpcEgressGateway, err error) {
	emptyResult := &v1.VpcEgressGateway{}
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceActionWithOptions(vpcegressgatewaysResource, c.ns, name, pt, data, opts, subresources...), emptyResult)

	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1.VpcEgressGateway), err
}