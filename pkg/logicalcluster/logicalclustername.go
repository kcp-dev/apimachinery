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

type LogicalClusterName string

func (cn LogicalClusterName) LogicalClusterPath() string {
	return path.Join("/clusters", string(cn))
}

func (cn LogicalClusterName) String() string {
	return string(cn)
}

func GetLogicalClusterName(obj v1.Object) LogicalClusterName {
	return LogicalClusterName(obj.GetClusterName())
}
