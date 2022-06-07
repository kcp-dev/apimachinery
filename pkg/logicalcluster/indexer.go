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

package logicalcluster

import (
	"fmt"

	"github.com/kcp-dev/logicalcluster"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
)

const (
	ByWorkspaceIndex = "kcp-byWorkspace"
)

// IndexByWorkspace returns an index using the underlying logical cluster name of the given object.
// It returns an error if the given object is not of meta/v1#Object type.
// It is meant to be consumed as an indexer function in SharedIndexInformer#AddIndexers.
func IndexByWorkspace(obj interface{}) ([]string, error) {
	metaObj, ok := obj.(metav1.Object)
	if !ok {
		return []string{}, fmt.Errorf("obj is supposed to be a metav1.Object, but is %T", obj)
	}

	lcluster := logicalcluster.From(metaObj)
	return []string{lcluster.String()}, nil
}

// AddByWorkspaceIndexer registers an indexer which indexes objects by their associated logical cluster.
// If the indexer is already registered, the function will be a no-op and return immediately without error.
// It returns an error if adding the indexer failed.
//
// The indexer can be referenced using the ByWorkspaceIndex const,
// i.e. informer.GetIndexer().ByIndex(logicalcluster.ByWorkspaceIndex, "clusterFoo")
func AddByWorkspaceIndexer(informer cache.SharedIndexInformer) error {
	if _, found := informer.GetIndexer().GetIndexers()[ByWorkspaceIndex]; !found {
		err := informer.AddIndexers(cache.Indexers{
			ByWorkspaceIndex: IndexByWorkspace,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
