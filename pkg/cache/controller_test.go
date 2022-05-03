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

	"github.com/kcp-dev/apimachinery/pkg/logicalcluster"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestClusterIndexFunc(t *testing.T) {
	tests := []struct {
		obj     *corev1.ConfigMap
		desired string
	}{
		{obj: &corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{ClusterName: "test"}}, desired: "test//"},
		{obj: &corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{ClusterName: "test-with-dashes"}}, desired: "test-with-dashes//"},
	}
	for _, tt := range tests {
		t.Run(tt.desired, func(t *testing.T) {
			result, err := ClusterIndexFunc(tt.obj)
			if err != nil {
				t.Error(err)
			}
			if result[0] != tt.desired {
				t.Errorf("got %v, want %v", result[0], tt.desired)
			}
			clusterName := logicalcluster.From(tt.obj).String()
			namespace := tt.obj.GetNamespace()
			name := tt.obj.GetName()
			key := ToClusterAwareKey(clusterName, namespace, name)

			if result[0] != key {
				t.Errorf("Index and Keyfunc have diverged, got %v, want %v", result[0], key)
			}
		})
	}
}

func TestClusterAndNamespaceIndexFunc(t *testing.T) {
	tests := []struct {
		obj     *corev1.ConfigMap
		desired string
	}{
		{obj: &corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{ClusterName: "test"}}, desired: "test//"},
		{obj: &corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{ClusterName: "test", Namespace: "testing"}}, desired: "test/testing/"},
	}
	for _, tt := range tests {
		t.Run(tt.desired, func(t *testing.T) {
			result, err := ClusterAndNamespaceIndexFunc(tt.obj)
			if err != nil {
				t.Error(err)
			}
			if result[0] != tt.desired {
				t.Errorf("got %v, want %v", result[0], tt.desired)
			}

			clusterName := logicalcluster.From(tt.obj).String()
			namespace := tt.obj.GetNamespace()
			name := tt.obj.GetName()
			key := ToClusterAwareKey(clusterName, namespace, name)

			if result[0] != key {
				t.Errorf("Index and Keyfunc have diverged, got %v, want %v", result[0], key)
			}
		})
	}
}

func TestClusterAwareKeyFunc(t *testing.T) {
	tests := []struct {
		obj     *corev1.ConfigMap
		desired string
	}{
		{obj: &corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{ClusterName: "cluster", Namespace: "namespace"}}, desired: "cluster/namespace/"},
		{obj: &corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{ClusterName: "cluster", Namespace: "namespace", Name: "name"}}, desired: "cluster/namespace/name"},
		{obj: &corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{ClusterName: "cluster", Name: "name"}}, desired: "cluster//name"},
	}
	for _, tt := range tests {
		t.Run(tt.desired, func(t *testing.T) {
			keyFuncResult, err := ClusterAwareKeyFunc(tt.obj)
			if err != nil {
				t.Error(err)
			}
			if keyFuncResult != tt.desired {
				t.Errorf("got %v, want %v", keyFuncResult, tt.desired)
			}
			clusterName := logicalcluster.From(tt.obj).String()
			namespace := tt.obj.GetNamespace()
			name := tt.obj.GetName()

			key := ToClusterAwareKey(clusterName, namespace, name)
			if key != keyFuncResult {
				t.Errorf("got %v, want %v", key, tt.desired)
			}
		})
	}
}
