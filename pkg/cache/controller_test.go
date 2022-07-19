/*
Copyright 2022 The KCP Authors.

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

package cache

import (
	"testing"

	"github.com/kcp-dev/logicalcluster/v2"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
)

func newUnstructured(cluster, namespace, name string, labels labels.Set) *unstructured.Unstructured {
	u := new(unstructured.Unstructured)
	u.SetZZZ_DeprecatedClusterName(cluster)
	u.SetNamespace(namespace)
	u.SetName(name)
	u.SetLabels(labels)
	return u
}

func TestClusterIndexFunc(t *testing.T) {
	tests := map[string]struct {
		obj     *unstructured.Unstructured
		desired string
	}{
		"bare cluster":             {obj: newUnstructured("test", "", "name", nil), desired: "test//"},
		"cluster and namespace":    {obj: newUnstructured("test", "namespace", "name", nil), desired: "test//"},
		"bare cluster with dashes": {obj: newUnstructured("test-with-dashes", "", "name", nil), desired: "test-with-dashes//"},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result, err := ClusterIndexFunc(tt.obj)
			require.NoError(t, err, "unexpected error calling ClusterIndexFunc")
			require.Len(t, result, 1, "ClusterIndexFunc should return one result")
			require.Equal(t, result[0], tt.desired)

			clusterName := logicalcluster.From(tt.obj).String()
			key := ToClusterAwareKey(clusterName, "", "")

			require.Equal(t, result[0], key, "ClusterIndexFunc and ToClusterAwareKey functions have diverged")
		})
	}
}

func TestClusterAndNamespaceIndexFunc(t *testing.T) {
	tests := map[string]struct {
		obj     *unstructured.Unstructured
		desired string
	}{
		"bare cluster":          {obj: newUnstructured("test", "", "name", nil), desired: "test//"},
		"cluster and namespace": {obj: newUnstructured("test", "testing", "name", nil), desired: "test/testing/"},
	}
	for _, tt := range tests {
		t.Run(tt.desired, func(t *testing.T) {
			result, err := ClusterAndNamespaceIndexFunc(tt.obj)
			require.NoError(t, err, "unexpected error calling ClusterAndNamespaceIndexFunc")
			require.Len(t, result, 1, "ClusterIndexFunc should return one result")
			require.Equal(t, result[0], tt.desired)

			clusterName := logicalcluster.From(tt.obj).String()
			namespace := tt.obj.GetNamespace()
			key := ToClusterAwareKey(clusterName, namespace, "")

			require.Equal(t, result[0], key, "ClusterAndNamespaceIndexFunc and ToClusterAwareKey functions have diverged")
		})
	}
}

func TestClusterAwareKeyFunc(t *testing.T) {
	tests := map[string]struct {
		obj     *unstructured.Unstructured
		desired string
	}{
		"cluster, namespace and name": {obj: newUnstructured("cluster", "namespace", "name", nil), desired: "cluster/namespace/name"},
		"cluster and name":            {obj: newUnstructured("cluster", "", "name", nil), desired: "cluster//name"},
	}
	for _, tt := range tests {
		t.Run(tt.desired, func(t *testing.T) {
			keyFuncResult, err := ClusterAwareKeyFunc(tt.obj)
			require.NoError(t, err, "unexpected error calling ClusterAwareKeyFunc")
			require.Equal(t, keyFuncResult, tt.desired)

			clusterName := logicalcluster.From(tt.obj).String()
			namespace := tt.obj.GetNamespace()
			name := tt.obj.GetName()

			key := ToClusterAwareKey(clusterName, namespace, name)
			require.Equal(t, key, keyFuncResult, "ClusterAwareKeyFunc and ToClusterAwareKey functions have diverged")
		})
	}
}

func TestSplitClusterAwareKey(t *testing.T) {
	tests := map[string]struct {
		clusterKey  string
		cluster     string
		namespace   string
		name        string
		shouldError bool
	}{
		"cluster + namespace + name": {clusterKey: "c1/ns1/n1", cluster: "c1", namespace: "ns1", name: "n1"},
		"cluster + name":             {clusterKey: "c1//n1", cluster: "c1", name: "n1"},
		"namespace + name":           {clusterKey: "/ns1/n1", namespace: "ns1", name: "n1"},
		"invalid":                    {clusterKey: "ns1/n1", namespace: "ns1", name: "n1", shouldError: true},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if tt.shouldError {
				_, _, _, err := SplitClusterAwareKey(tt.clusterKey)
				require.Error(t, err, "An invalid cluster key should result in an error")
				return
			}
			clusterAwareKey := ToClusterAwareKey(tt.cluster, tt.namespace, tt.name)
			require.Equal(t, clusterAwareKey, tt.clusterKey, "ToClusterAwareKey has changed, these tests may need to be updated")

			cluster, namespace, name, err := SplitClusterAwareKey(tt.clusterKey)
			require.NoError(t, err, "Splitting a valid cluster-aware key should not result in an error")

			require.Equal(t, cluster, tt.cluster, "cluster did not match after split")
			require.Equal(t, namespace, tt.namespace, "namespace did not match after split")
			require.Equal(t, name, tt.name, "name did not match after split")

		})
	}
}
