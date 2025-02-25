/*
Copyright 2025 Rancher Labs, Inc.

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

// Code generated by main. DO NOT EDIT.

package fake

import (
	"context"

	v1 "github.com/rancher/rancher/pkg/apis/rke.cattle.io/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeRKEBootstrapTemplates implements RKEBootstrapTemplateInterface
type FakeRKEBootstrapTemplates struct {
	Fake *FakeRkeV1
	ns   string
}

var rkebootstraptemplatesResource = v1.SchemeGroupVersion.WithResource("rkebootstraptemplates")

var rkebootstraptemplatesKind = v1.SchemeGroupVersion.WithKind("RKEBootstrapTemplate")

// Get takes name of the rKEBootstrapTemplate, and returns the corresponding rKEBootstrapTemplate object, and an error if there is any.
func (c *FakeRKEBootstrapTemplates) Get(ctx context.Context, name string, options metav1.GetOptions) (result *v1.RKEBootstrapTemplate, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(rkebootstraptemplatesResource, c.ns, name), &v1.RKEBootstrapTemplate{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1.RKEBootstrapTemplate), err
}

// List takes label and field selectors, and returns the list of RKEBootstrapTemplates that match those selectors.
func (c *FakeRKEBootstrapTemplates) List(ctx context.Context, opts metav1.ListOptions) (result *v1.RKEBootstrapTemplateList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(rkebootstraptemplatesResource, rkebootstraptemplatesKind, c.ns, opts), &v1.RKEBootstrapTemplateList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1.RKEBootstrapTemplateList{ListMeta: obj.(*v1.RKEBootstrapTemplateList).ListMeta}
	for _, item := range obj.(*v1.RKEBootstrapTemplateList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested rKEBootstrapTemplates.
func (c *FakeRKEBootstrapTemplates) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(rkebootstraptemplatesResource, c.ns, opts))

}

// Create takes the representation of a rKEBootstrapTemplate and creates it.  Returns the server's representation of the rKEBootstrapTemplate, and an error, if there is any.
func (c *FakeRKEBootstrapTemplates) Create(ctx context.Context, rKEBootstrapTemplate *v1.RKEBootstrapTemplate, opts metav1.CreateOptions) (result *v1.RKEBootstrapTemplate, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(rkebootstraptemplatesResource, c.ns, rKEBootstrapTemplate), &v1.RKEBootstrapTemplate{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1.RKEBootstrapTemplate), err
}

// Update takes the representation of a rKEBootstrapTemplate and updates it. Returns the server's representation of the rKEBootstrapTemplate, and an error, if there is any.
func (c *FakeRKEBootstrapTemplates) Update(ctx context.Context, rKEBootstrapTemplate *v1.RKEBootstrapTemplate, opts metav1.UpdateOptions) (result *v1.RKEBootstrapTemplate, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(rkebootstraptemplatesResource, c.ns, rKEBootstrapTemplate), &v1.RKEBootstrapTemplate{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1.RKEBootstrapTemplate), err
}

// Delete takes name of the rKEBootstrapTemplate and deletes it. Returns an error if one occurs.
func (c *FakeRKEBootstrapTemplates) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteActionWithOptions(rkebootstraptemplatesResource, c.ns, name, opts), &v1.RKEBootstrapTemplate{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeRKEBootstrapTemplates) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(rkebootstraptemplatesResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1.RKEBootstrapTemplateList{})
	return err
}

// Patch applies the patch and returns the patched rKEBootstrapTemplate.
func (c *FakeRKEBootstrapTemplates) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.RKEBootstrapTemplate, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(rkebootstraptemplatesResource, c.ns, name, pt, data, subresources...), &v1.RKEBootstrapTemplate{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1.RKEBootstrapTemplate), err
}
