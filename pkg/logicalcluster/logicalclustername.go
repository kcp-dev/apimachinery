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
	"path"

	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

// LogicalClusterName is the name of the logical cluster. A logical cluster is
// 1. a (part of) etcd prefix to store objects in that cluster
// 2. a (part of) an http path which serves a Kubernetes-cluster-like API with
//    discovery, OpenAPI and the actual API groups.
// 3. a value in metadata.clusterName in objects from cross-workspace list/watches,
//    which is used to identify the logical cluster.
type LogicalClusterName string

// LogicalClusterPath returns a path segment for the logical cluster to access its API.
func (cn LogicalClusterName) LogicalClusterPath() string {
	return path.Join("/clusters", string(cn))
}

// String returns the string representation of the logical cluster name.
func (cn LogicalClusterName) String() string {
	return string(cn)
}

// From returns a logical cluster name from an Object's
// metadata.clusterName.
func From(obj v1.Object) LogicalClusterName {
	return LogicalClusterName(obj.GetClusterName())
}
