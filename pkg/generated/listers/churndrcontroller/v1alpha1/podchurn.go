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

// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "github.com/llimon/churndr/pkg/apis/churndrcontroller/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// PodchurnLister helps list Podchurns.
// All objects returned here must be treated as read-only.
type PodchurnLister interface {
	// List lists all Podchurns in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.Podchurn, err error)
	// Podchurns returns an object that can list and get Podchurns.
	Podchurns(namespace string) PodchurnNamespaceLister
	PodchurnListerExpansion
}

// podchurnLister implements the PodchurnLister interface.
type podchurnLister struct {
	indexer cache.Indexer
}

// NewPodchurnLister returns a new PodchurnLister.
func NewPodchurnLister(indexer cache.Indexer) PodchurnLister {
	return &podchurnLister{indexer: indexer}
}

// List lists all Podchurns in the indexer.
func (s *podchurnLister) List(selector labels.Selector) (ret []*v1alpha1.Podchurn, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.Podchurn))
	})
	return ret, err
}

// Podchurns returns an object that can list and get Podchurns.
func (s *podchurnLister) Podchurns(namespace string) PodchurnNamespaceLister {
	return podchurnNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// PodchurnNamespaceLister helps list and get Podchurns.
// All objects returned here must be treated as read-only.
type PodchurnNamespaceLister interface {
	// List lists all Podchurns in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.Podchurn, err error)
	// Get retrieves the Podchurn from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1alpha1.Podchurn, error)
	PodchurnNamespaceListerExpansion
}

// podchurnNamespaceLister implements the PodchurnNamespaceLister
// interface.
type podchurnNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all Podchurns in the indexer for a given namespace.
func (s podchurnNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.Podchurn, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.Podchurn))
	})
	return ret, err
}

// Get retrieves the Podchurn from the indexer for a given namespace and name.
func (s podchurnNamespaceLister) Get(name string) (*v1alpha1.Podchurn, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("podchurn"), name)
	}
	return obj.(*v1alpha1.Podchurn), nil
}
