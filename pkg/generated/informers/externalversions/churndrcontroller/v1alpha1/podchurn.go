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

// Code generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	time "time"

	churndrcontrollerv1alpha1 "github.com/llimon/churndr/pkg/apis/churndrcontroller/v1alpha1"
	versioned "github.com/llimon/churndr/pkg/generated/clientset/versioned"
	internalinterfaces "github.com/llimon/churndr/pkg/generated/informers/externalversions/internalinterfaces"
	v1alpha1 "github.com/llimon/churndr/pkg/generated/listers/churndrcontroller/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// PodchurnInformer provides access to a shared informer and lister for
// Podchurns.
type PodchurnInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha1.PodchurnLister
}

type podchurnInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewPodchurnInformer constructs a new informer for Podchurn type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewPodchurnInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredPodchurnInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredPodchurnInformer constructs a new informer for Podchurn type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredPodchurnInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.ChurndrcontrollerV1alpha1().Podchurns(namespace).List(options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.ChurndrcontrollerV1alpha1().Podchurns(namespace).Watch(options)
			},
		},
		&churndrcontrollerv1alpha1.Podchurn{},
		resyncPeriod,
		indexers,
	)
}

func (f *podchurnInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredPodchurnInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *podchurnInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&churndrcontrollerv1alpha1.Podchurn{}, f.defaultInformer)
}

func (f *podchurnInformer) Lister() v1alpha1.PodchurnLister {
	return v1alpha1.NewPodchurnLister(f.Informer().GetIndexer())
}