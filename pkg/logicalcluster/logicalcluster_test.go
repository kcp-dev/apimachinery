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
	"encoding/json"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogicalCluster_Split(t *testing.T) {
	tests := []struct {
		cn     LogicalCluster
		parent LogicalCluster
		name   string
	}{
		{New(""), New(""), ""},
		{New("foo"), New(""), "foo"},
		{New("foo:bar"), New("foo"), "bar"},
		{New("foo:bar:baz"), New("foo:bar"), "baz"},
		{New("foo::baz"), New("foo:"), "baz"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotParent, gotName := tt.cn.Split()
			if gotParent != tt.parent {
				t.Errorf("Split() gotParent = %v, want %v", gotParent, tt.parent)
			}
			if gotName != tt.name {
				t.Errorf("Split() gotName = %v, want %v", gotName, tt.name)
			}
		})
	}
}

func TestJSON(t *testing.T) {
	type JWT struct {
		I  int            `json:"i"`
		CN LogicalCluster `json:"cn"`
	}

	jwt := JWT{
		I:  1,
		CN: New("foo:bar"),
	}

	bs, err := json.Marshal(jwt)
	require.NoError(t, err)
	require.Equal(t, `{"i":1,"cn":"foo:bar"}`, string(bs))

	var jwt2 JWT
	err = json.Unmarshal(bs, &jwt2)
	require.NoError(t, err)
	require.Equal(t, jwt, jwt2)
}

func TestNew(t *testing.T) {
	tests := []struct {
		value string
		want  LogicalCluster
	}{
		{"", LogicalCluster{}},
		{"root", LogicalCluster{value: "root"}},
		{"root:foo:bar", LogicalCluster{value: "root:foo:bar"}},
		{"system:admin", LogicalCluster{value: "system:admin"}},

		{":foo:bar", LogicalCluster{value: "root:foo:bar"}},
	}
	for _, tt := range tests {
		t.Run(tt.value, func(t *testing.T) {
			if got := New(tt.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrettyString(t *testing.T) {
	tests := []struct {
		cluster string
		want    string
	}{
		{"", ""},
		{"root", "root"},
		{"root:foo:bar", ":foo:bar"},
		{"system:admin", "system:admin"},
	}
	for _, tt := range tests {
		t.Run(tt.cluster, func(t *testing.T) {
			lc := LogicalCluster{value: tt.cluster}
			if got := lc.PrettyString(); got != tt.want {
				t.Errorf("PrettyString() = %v, want %v", got, tt.want)
			}
		})
	}
}
