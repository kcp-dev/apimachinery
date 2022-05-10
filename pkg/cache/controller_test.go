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

	"github.com/kcp-dev/logicalcluster"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
)

func newUnstructured(cluster, namespace, name string, labels labels.Set) *unstructured.Unstructured {
	u := new(unstructured.Unstructured)
	u.SetClusterName(cluster)
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
