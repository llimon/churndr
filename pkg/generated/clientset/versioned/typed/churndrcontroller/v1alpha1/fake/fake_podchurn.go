/*
Copyright 2019, Luis E Limon

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
	v1alpha1 "github.com/llimon/churndr/pkg/apis/churndrcontroller/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakePodchurns implements PodchurnInterface
type FakePodchurns struct {
	Fake *FakeChurndrcontrollerV1alpha1
	ns   string
}

var podchurnsResource = schema.GroupVersionResource{Group: "churndrcontroller.churndr.com", Version: "v1alpha1", Resource: "podchurns"}

var podchurnsKind = schema.GroupVersionKind{Group: "churndrcontroller.churndr.com", Version: "v1alpha1", Kind: "Podchurn"}

// Get takes name of the podchurn, and returns the corresponding podchurn object, and an error if there is any.
func (c *FakePodchurns) Get(name string, options v1.GetOptions) (result *v1alpha1.Podchurn, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(podchurnsResource, c.ns, name), &v1alpha1.Podchurn{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Podchurn), err
}

// List takes label and field selectors, and returns the list of Podchurns that match those selectors.
func (c *FakePodchurns) List(opts v1.ListOptions) (result *v1alpha1.PodchurnList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(podchurnsResource, podchurnsKind, c.ns, opts), &v1alpha1.PodchurnList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.PodchurnList{ListMeta: obj.(*v1alpha1.PodchurnList).ListMeta}
	for _, item := range obj.(*v1alpha1.PodchurnList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested podchurns.
func (c *FakePodchurns) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(podchurnsResource, c.ns, opts))

}

// Create takes the representation of a podchurn and creates it.  Returns the server's representation of the podchurn, and an error, if there is any.
func (c *FakePodchurns) Create(podchurn *v1alpha1.Podchurn) (result *v1alpha1.Podchurn, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(podchurnsResource, c.ns, podchurn), &v1alpha1.Podchurn{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Podchurn), err
}

// Update takes the representation of a podchurn and updates it. Returns the server's representation of the podchurn, and an error, if there is any.
func (c *FakePodchurns) Update(podchurn *v1alpha1.Podchurn) (result *v1alpha1.Podchurn, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(podchurnsResource, c.ns, podchurn), &v1alpha1.Podchurn{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Podchurn), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakePodchurns) UpdateStatus(podchurn *v1alpha1.Podchurn) (*v1alpha1.Podchurn, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(podchurnsResource, "status", c.ns, podchurn), &v1alpha1.Podchurn{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Podchurn), err
}

// Delete takes name of the podchurn and deletes it. Returns an error if one occurs.
func (c *FakePodchurns) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(podchurnsResource, c.ns, name), &v1alpha1.Podchurn{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakePodchurns) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(podchurnsResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &v1alpha1.PodchurnList{})
	return err
}

// Patch applies the patch and returns the patched podchurn.
func (c *FakePodchurns) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.Podchurn, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(podchurnsResource, c.ns, name, pt, data, subresources...), &v1alpha1.Podchurn{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Podchurn), err
}