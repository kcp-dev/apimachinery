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
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

// LogicalCluster is the name of a logical cluster. A logical cluster is
// 1. a (part of) etcd prefix to store objects in that cluster
// 2. a (part of) a http path which serves a Kubernetes-cluster-like API with
//    discovery, OpenAPI and the actual API groups.
// 3. a value in metadata.clusterName in objects from cross-workspace list/watches,
//    which is used to identify the logical cluster.
//
// A logical cluster is of type string, and is either "root" or a logical cluster
// (called parent) followed by the separator ":" and a word. In other words, it is
// like a path, but with colons instead of slashes.
type LogicalCluster string

const seperator = ":"
const Root = LogicalCluster("root")

// Path returns a path segment for the logical cluster to access its API.
func (cn LogicalCluster) Path() string {
	return path.Join("/clusters", string(cn))
}

// String returns the string representation of the logical cluster name.
func (cn LogicalCluster) String() string {
	return string(cn)
}

// From returns a logical cluster name from an Object's
// metadata.clusterName.
func From(obj v1.Object) LogicalCluster {
	return LogicalCluster(obj.GetClusterName())
}

// Parent returns the parent logical cluster name of the given logical cluster name.
func (cn LogicalCluster) Parent() (LogicalCluster, bool) {
	if cn == Root {
		return Root, false
	}
	parent, _ := cn.Split()
	return parent, true
}

// Split splits logical cluster immediately following the final colon,
// separating it into a parent logical cluster and name component.
// If there is no colon in path, Split returns an empty logical cluster name
// and name set to path.
// The returned values have the property that lcn = dir+file.
func (cn LogicalCluster) Split() (parent LogicalCluster, name string) {
	i := strings.LastIndex(string(cn), seperator)
	return cn[:i+1], string(cn)[i+1:]
}

// IsRoot returns true if the logical cluster name is the root logical cluster.
func (cn LogicalCluster) IsRoot() bool {
	return cn == Root
}

// Join joins a parent logical cluster name and a name component.
func (cn LogicalCluster) Join(name string) LogicalCluster {
	return LogicalCluster(string(cn) + seperator + name)
}

// IsValid returns true if the logical cluster name is valid.
func (cn LogicalCluster) IsValid() bool {
	return cn == Root || (strings.Contains(string(cn), seperator) && strings.HasPrefix(string(cn), Root.String()+seperator))
}
