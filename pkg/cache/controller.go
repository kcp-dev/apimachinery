package cache

import (
	"fmt"
	"strings"

	"github.com/kcp-dev/apimachinery/pkg/logicalcluster"
	"k8s.io/apimachinery/pkg/api/meta"
)

const (
	// ClusterIndexName is the name of the index that allows you to filter by cluster
	ClusterIndexName = "cluster"
	// ClusterAndNamespaceIndexName is the name of index that allows you to filter by cluster and namespace
	ClusterAndNamespaceIndexName = "cluster-and-namespace"
)

// ClusterIndexFunc indexes by cluster name
func ClusterIndexFunc(obj interface{}) ([]string, error) {
	meta, err := meta.Accessor(obj)
	if err != nil {
		return []string{}, fmt.Errorf("object has no meta: %v", err)
	}
	return []string{logicalcluster.From(meta).String()}, nil
}

// ClusterAndNamespaceIndexFunc indexes by cluster and namespace name
func ClusterAndNamespaceIndexFunc(obj interface{}) ([]string, error) {
	meta, err := meta.Accessor(obj)
	if err != nil {
		return []string{}, fmt.Errorf("object has no meta: %v", err)
	}
	// TODO(fabianvf): Should I call ClusterAwareKeyFunc on this to ensure the formatting will always match?
	return []string{logicalcluster.From(meta).String() + "/" + meta.GetNamespace()}, nil

}

// ClusterAwareKeyFunc keys on cluster, namespace and name
func ClusterAwareKeyFunc(obj interface{}) (string, error) {
	meta, err := meta.Accessor(obj)
	if err != nil {
		return "", fmt.Errorf("object has no meta: %v", err)
	}
	clusterName := logicalcluster.From(meta).String()
	namespace := meta.GetNamespace()
	name := meta.GetName()

	return strings.Join([]string{clusterName, namespace, name}, "/"), nil
}
