package cache

import (
	"github.com/kcp-dev/apimachinery/pkg/logicalcluster"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/tools/cache"
)

// GenericLister is a lister skin on a generic Indexer
type GenericClusterLister interface {
	// List will return all objects across clusters
	List(selector labels.Selector) (ret []runtime.Object, err error)
	// ByCluster will give you a GenericLister for one namespace
	ByCluster(cluster logicalcluster.LogicalCluster) cache.GenericLister
}

// NewGenericClusterLister creates a new instance for the genericClusterLister.
func NewGenericClusterLister(indexer cache.Indexer, resource schema.GroupResource) *genericClusterLister {
	return &genericClusterLister{indexer: indexer, resource: resource}
}

type genericClusterLister struct {
	indexer  cache.Indexer
	resource schema.GroupResource
}

func (s *genericClusterLister) List(selector labels.Selector) (ret []runtime.Object, err error) {
	if selector == nil {
		selector = labels.NewSelector()
	}
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(runtime.Object))
	})
	return ret, err
}

func (s *genericClusterLister) ByCluster(cluster logicalcluster.LogicalCluster) cache.GenericLister {
	return &genericLister{indexer: s.indexer, resource: s.resource, cluster: cluster}
}

type genericLister struct {
	indexer  cache.Indexer
	cluster  logicalcluster.LogicalCluster
	resource schema.GroupResource
}

func (s *genericLister) List(selector labels.Selector) (ret []runtime.Object, err error) {
	selectAll := selector == nil || selector.Empty()
	list, err := s.indexer.ByIndex(ClusterIndexName, s.cluster.String())
	if err != nil {
		return nil, err
	}

	if selector == nil {
		selector = labels.Everything()
	}
	for i := range list {
		item := list[i].(runtime.Object)
		if selectAll {
			ret = append(ret, item)
		} else {
			metadata, err := meta.Accessor(item)
			if err != nil {
				return nil, err
			}
			if selector.Matches(labels.Set(metadata.GetLabels())) {
				ret = append(ret, item)
			}
		}
	}

	return ret, err
}

func (s *genericLister) Get(name string) (runtime.Object, error) {
	metadata := &metav1.ObjectMeta{
		ClusterName: s.cluster.String(),
		Name:        name,
	}
	obj, exists, err := s.indexer.Get(metadata)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(s.resource, name)
	}
	return obj.(runtime.Object), nil
}

func (s *genericLister) ByNamespace(namespace string) cache.GenericNamespaceLister {
	return &genericNamespaceLister{indexer: s.indexer, namespace: namespace, resource: s.resource, cluster: s.cluster}
}

type genericNamespaceLister struct {
	indexer   cache.Indexer
	cluster   logicalcluster.LogicalCluster
	namespace string
	resource  schema.GroupResource
}

func (s *genericNamespaceLister) List(selector labels.Selector) (ret []runtime.Object, err error) {
	selectAll := selector == nil || selector.Empty()
	list, err := s.indexer.Index(ClusterAndNamespaceIndexName, &metav1.ObjectMeta{
		ClusterName: s.cluster.String(),
		Namespace:   s.namespace,
	})
	if err != nil {
		return nil, err
	}

	for i := range list {
		item := list[i].(runtime.Object)
		if selectAll {
			ret = append(ret, item)
		} else {
			metadata, err := meta.Accessor(item)
			if err != nil {
				return nil, err
			}
			if selector.Matches(labels.Set(metadata.GetLabels())) {
				ret = append(ret, item)
			}
		}
	}
	return ret, err
}

func (s *genericNamespaceLister) Get(name string) (runtime.Object, error) {
	metadata := &metav1.ObjectMeta{
		ClusterName: s.cluster.String(),
		Namespace:   s.namespace,
		Name:        name,
	}
	obj, exists, err := s.indexer.Get(metadata)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(s.resource, name)
	}
	return obj.(runtime.Object), nil
}
