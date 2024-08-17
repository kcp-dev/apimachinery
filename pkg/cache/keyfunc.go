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
	"fmt"
	"strings"

	"github.com/kcp-dev/logicalcluster/v3"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/client-go/tools/cache"
)

// Key identifies an object of some implict kind, possibly where
// cluster or namespace are not issues.
// Put another way, Key identifies an instance of some implicit resource.
// Put another way, Key is what SplitMetaClusterNamespaceKey wants to return,
// but passing one of these instead of a string means that the receiver never
// has to worry about syntax errors.
type Key struct {
	Cluster   logicalcluster.Name
	Namespace string
	Name      string
}

// NewKey constructs a Key
func NewKey(cluster logicalcluster.Name, namespace, name string) Key {
	return Key{Cluster: cluster, Namespace: namespace, Name: name}
}

func (key Key) Parts() (cluster logicalcluster.Name, namespace, name string) {
	return key.Cluster, key.Namespace, key.Name
}

// String returns the standard string representation of the given Key
func (key Key) String() string {
	var ans string
	if key.Cluster != "" {
		ans += key.Cluster.String() + "|"
	}
	if key.Namespace != "" {
		ans += key.Namespace + "/"
	}
	ans += key.Name
	return ans

}

// ParseKey inverts the usual encoding, complaining on syntax error
func ParseKey(encoded string) (Key, error) {
	var key Key
	var err error
	key.Cluster, key.Namespace, key.Name, err = SplitMetaClusterNamespaceKey(encoded)
	return key, err
}

// DeletionHandlingMetaClusterNamespaceKeyFunc checks for
// DeletedFinalStateUnknown objects before calling
// MetaClusterNamespaceKeyFunc.
func DeletionHandlingMetaClusterNamespaceKeyFunc(obj interface{}) (string, error) {
	if d, ok := obj.(cache.DeletedFinalStateUnknown); ok {
		return d.Key, nil
	}
	return MetaClusterNamespaceKeyFunc(obj)
}

// ObjMetaClusterNamespaceKey is a convenient default KeyFunc which knows how to make
// structured keys for API objects which implement meta.Interface.
// This is the structured alternative to MetaClusterNamespaceKeyFunc;
// putting such an object reference in a queue means that no parsing errors are possible downstream.
func ObjMetaClusterNamespaceKey(obj interface{}) (Key, error) {
	if key, ok := obj.(cache.ExplicitKey); ok {
		return ParseKey(string(key))
	}
	meta, err := meta.Accessor(obj)
	if err != nil {
		return Key{}, fmt.Errorf("object has no meta: %v", err)
	}
	return Key{Cluster: logicalcluster.From(meta), Namespace: meta.GetNamespace(), Name: meta.GetName()}, nil
}

// MetaClusterNamespaceKeyFunc is a convenient default KeyFunc which knows how to make
// keys for API objects which implement meta.Interface.
// The key uses the format <clusterName>|<namespace>/<name> unless <namespace> is empty, then
// it's just <clusterName>|<name>, and if running in a single-cluster context where no explicit
// cluster name is given, it's just <name>.
func MetaClusterNamespaceKeyFunc(obj interface{}) (string, error) {
	key, err := ObjMetaClusterNamespaceKey(obj)
	return key.String(), err
}

// ToClusterAwareKey formats a cluster, namespace, and name as a key.
func ToClusterAwareKey(cluster, namespace, name string) string {
	return Key{Cluster: logicalcluster.Name(cluster), Namespace: namespace, Name: name}.String()
}

// SplitMetaClusterNamespaceKey returns the cluster, namespace, and name that
// MetaClusterNamespaceKeyFunc encoded into key.
func SplitMetaClusterNamespaceKey(key string) (clusterName logicalcluster.Name, namespace, name string, err error) {
	invalidKey := fmt.Errorf("unexpected key format: %q", key)
	outerParts := strings.Split(key, "|")
	switch len(outerParts) {
	case 1:
		namespace, name, err := cache.SplitMetaNamespaceKey(outerParts[0])
		if err != nil {
			err = invalidKey
		}
		return "", namespace, name, err
	case 2:
		namespace, name, err := cache.SplitMetaNamespaceKey(outerParts[1])
		if err != nil {
			err = invalidKey
		}
		return logicalcluster.Name(outerParts[0]), namespace, name, err
	default:
		return "", "", "", invalidKey
	}
}
