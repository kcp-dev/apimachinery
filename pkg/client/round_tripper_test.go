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

package client

import (
	"testing"

	"github.com/kcp-dev/logicalcluster"
)

func TestRoundTripper_generatePath(t *testing.T) {
	tests := []struct {
		originalPath string
		cluster      logicalcluster.Name
		desired      string
	}{
		{"", logicalcluster.New("test"), "/clusters/test"},
		{"/prefix/", logicalcluster.New("test"), "/clusters/test/prefix/"},
		{"/several/divisions/of/prefix", logicalcluster.New("test"), "/clusters/test/several/divisions/of/prefix"},
		{"/prefix", logicalcluster.New("test"), "/clusters/test/prefix"},
		{"prefix", logicalcluster.New("test"), "/clusters/test/prefix"},
	}
	for _, tt := range tests {
		t.Run(tt.desired, func(t *testing.T) {
			result := generatePath(tt.originalPath, tt.cluster)
			if result != tt.desired {
				t.Errorf("got %v, want %v", result, tt.desired)
			}
		})
	}
}
