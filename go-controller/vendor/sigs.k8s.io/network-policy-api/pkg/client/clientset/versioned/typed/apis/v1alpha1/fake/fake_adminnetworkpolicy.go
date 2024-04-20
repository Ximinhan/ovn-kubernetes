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
	json "encoding/json"
	"fmt"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
	v1alpha1 "sigs.k8s.io/network-policy-api/apis/v1alpha1"
	apisv1alpha1 "sigs.k8s.io/network-policy-api/pkg/client/applyconfiguration/apis/v1alpha1"
)

// FakeAdminNetworkPolicies implements AdminNetworkPolicyInterface
type FakeAdminNetworkPolicies struct {
	Fake *FakePolicyV1alpha1
}

var adminnetworkpoliciesResource = v1alpha1.SchemeGroupVersion.WithResource("adminnetworkpolicies")

var adminnetworkpoliciesKind = v1alpha1.SchemeGroupVersion.WithKind("AdminNetworkPolicy")

// Get takes name of the adminNetworkPolicy, and returns the corresponding adminNetworkPolicy object, and an error if there is any.
func (c *FakeAdminNetworkPolicies) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.AdminNetworkPolicy, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(adminnetworkpoliciesResource, name), &v1alpha1.AdminNetworkPolicy{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.AdminNetworkPolicy), err
}

// List takes label and field selectors, and returns the list of AdminNetworkPolicies that match those selectors.
func (c *FakeAdminNetworkPolicies) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.AdminNetworkPolicyList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(adminnetworkpoliciesResource, adminnetworkpoliciesKind, opts), &v1alpha1.AdminNetworkPolicyList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.AdminNetworkPolicyList{ListMeta: obj.(*v1alpha1.AdminNetworkPolicyList).ListMeta}
	for _, item := range obj.(*v1alpha1.AdminNetworkPolicyList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested adminNetworkPolicies.
func (c *FakeAdminNetworkPolicies) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(adminnetworkpoliciesResource, opts))
}

// Create takes the representation of a adminNetworkPolicy and creates it.  Returns the server's representation of the adminNetworkPolicy, and an error, if there is any.
func (c *FakeAdminNetworkPolicies) Create(ctx context.Context, adminNetworkPolicy *v1alpha1.AdminNetworkPolicy, opts v1.CreateOptions) (result *v1alpha1.AdminNetworkPolicy, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(adminnetworkpoliciesResource, adminNetworkPolicy), &v1alpha1.AdminNetworkPolicy{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.AdminNetworkPolicy), err
}

// Update takes the representation of a adminNetworkPolicy and updates it. Returns the server's representation of the adminNetworkPolicy, and an error, if there is any.
func (c *FakeAdminNetworkPolicies) Update(ctx context.Context, adminNetworkPolicy *v1alpha1.AdminNetworkPolicy, opts v1.UpdateOptions) (result *v1alpha1.AdminNetworkPolicy, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(adminnetworkpoliciesResource, adminNetworkPolicy), &v1alpha1.AdminNetworkPolicy{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.AdminNetworkPolicy), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeAdminNetworkPolicies) UpdateStatus(ctx context.Context, adminNetworkPolicy *v1alpha1.AdminNetworkPolicy, opts v1.UpdateOptions) (*v1alpha1.AdminNetworkPolicy, error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateSubresourceAction(adminnetworkpoliciesResource, "status", adminNetworkPolicy), &v1alpha1.AdminNetworkPolicy{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.AdminNetworkPolicy), err
}

// Delete takes name of the adminNetworkPolicy and deletes it. Returns an error if one occurs.
func (c *FakeAdminNetworkPolicies) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteActionWithOptions(adminnetworkpoliciesResource, name, opts), &v1alpha1.AdminNetworkPolicy{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeAdminNetworkPolicies) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(adminnetworkpoliciesResource, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.AdminNetworkPolicyList{})
	return err
}

// Patch applies the patch and returns the patched adminNetworkPolicy.
func (c *FakeAdminNetworkPolicies) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.AdminNetworkPolicy, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(adminnetworkpoliciesResource, name, pt, data, subresources...), &v1alpha1.AdminNetworkPolicy{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.AdminNetworkPolicy), err
}

// Apply takes the given apply declarative configuration, applies it and returns the applied adminNetworkPolicy.
func (c *FakeAdminNetworkPolicies) Apply(ctx context.Context, adminNetworkPolicy *apisv1alpha1.AdminNetworkPolicyApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.AdminNetworkPolicy, err error) {
	if adminNetworkPolicy == nil {
		return nil, fmt.Errorf("adminNetworkPolicy provided to Apply must not be nil")
	}
	data, err := json.Marshal(adminNetworkPolicy)
	if err != nil {
		return nil, err
	}
	name := adminNetworkPolicy.Name
	if name == nil {
		return nil, fmt.Errorf("adminNetworkPolicy.Name must be provided to Apply")
	}
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(adminnetworkpoliciesResource, *name, types.ApplyPatchType, data), &v1alpha1.AdminNetworkPolicy{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.AdminNetworkPolicy), err
}

// ApplyStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating ApplyStatus().
func (c *FakeAdminNetworkPolicies) ApplyStatus(ctx context.Context, adminNetworkPolicy *apisv1alpha1.AdminNetworkPolicyApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.AdminNetworkPolicy, err error) {
	if adminNetworkPolicy == nil {
		return nil, fmt.Errorf("adminNetworkPolicy provided to Apply must not be nil")
	}
	data, err := json.Marshal(adminNetworkPolicy)
	if err != nil {
		return nil, err
	}
	name := adminNetworkPolicy.Name
	if name == nil {
		return nil, fmt.Errorf("adminNetworkPolicy.Name must be provided to Apply")
	}
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(adminnetworkpoliciesResource, *name, types.ApplyPatchType, data, "status"), &v1alpha1.AdminNetworkPolicy{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.AdminNetworkPolicy), err
}
