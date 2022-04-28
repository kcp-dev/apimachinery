package cache

import (
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/api/meta"
)

const (
	ClusterIndexName             = "cluster"
	ClusterAndNamespaceIndexName = "cluster-and-namespace"
)

func ClusterIndexFunc(obj interface{}) ([]string, error) {
	meta, err := meta.Accessor(obj)
	if err != nil {
		return []string{""}, fmt.Errorf("object has no meta: %v", err)
	}
	// return []string{meta.GetZZZ_DeprecatedClusterName()}, nil
	index := []string{meta.GetZZZ_DeprecatedClusterName()}
	return index, nil
}

func ClusterAndNamespaceIndexFunc(obj interface{}) ([]string, error) {
	meta, err := meta.Accessor(obj)
	if err != nil {
		return []string{""}, fmt.Errorf("object has no meta: %v", err)
	}
	// return []string{meta.GetZZZ_DeprecatedClusterName() + "/" + meta.GetNamespace()}, nil
	index := []string{meta.GetZZZ_DeprecatedClusterName() + "/" + meta.GetNamespace()}
	return index, nil

}

func ClusterAwareKeyFunc(obj interface{}) (string, error) {
	meta, err := meta.Accessor(obj)
	if err != nil {
		return "", fmt.Errorf("object has no meta: %v", err)
	}
	clusterName := meta.GetZZZ_DeprecatedClusterName()
	namespace := meta.GetNamespace()
	name := meta.GetName()

	return strings.Join([]string{clusterName, namespace, name}, "/"), nil
}
