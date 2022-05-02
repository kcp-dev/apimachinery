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
	clusterName := logicalcluster.From(meta).String()
	return []string{ToClusterAwareKey(clusterName, "", "")}, nil
}

// ClusterAndNamespaceIndexFunc indexes by cluster and namespace name
func ClusterAndNamespaceIndexFunc(obj interface{}) ([]string, error) {
	meta, err := meta.Accessor(obj)
	if err != nil {
		return []string{}, fmt.Errorf("object has no meta: %v", err)
	}
	clusterName := logicalcluster.From(meta).String()
	return []string{ToClusterAwareKey(clusterName, meta.GetNamespace(), "")}, nil

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

	return ToClusterAwareKey(clusterName, namespace, name), nil
}

func ToClusterAwareKey(cluster, namespace, name string) string {
	return strings.Join([]string{cluster, namespace, name}, "/")
}
