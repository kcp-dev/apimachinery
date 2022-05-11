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

package discovery

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/kcp-dev/logicalcluster"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

func TestClusterDiscoveryClient(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		var obj interface{}
		switch req.URL.Path {
		case "/clusters/cluster1/api":
			obj = &metav1.APIVersions{
				Versions: []string{
					"v1",
				},
			}
		case "/clusters/cluster1/apis":
			obj = &metav1.APIGroupList{
				Groups: []metav1.APIGroup{
					{
						Name: "extensions",
						Versions: []metav1.GroupVersionForDiscovery{
							{GroupVersion: "extensions/v1beta1"},
						},
					},
				},
			}
		case "/clusters/cluster2/api":
			obj = &metav1.APIVersions{
				Versions: []string{
					"v2",
				},
			}
		case "/clusters/cluster2/apis":
			obj = &metav1.APIGroupList{
				Groups: []metav1.APIGroup{
					{
						Name: "extensions",
						Versions: []metav1.GroupVersionForDiscovery{
							{GroupVersion: "extensions/v2beta2"},
						},
					},
				},
			}
		default:
			w.WriteHeader(http.StatusNotFound)
			fmt.Println(req.URL)
			return
		}
		output, err := json.Marshal(obj)
		if err != nil {
			t.Fatalf("unexpected encoding error: %v", err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(output)
		require.NoError(t, err)
	}))

	defer server.Close()

	client, err := NewClusterDiscoveryClientForConfig(&rest.Config{Host: server.URL + "/"})
	require.NoError(t, err)

	cluster1Client := client.Cluster(logicalcluster.New("cluster1"))
	apiGroupList, err := cluster1Client.ServerGroups()
	require.NoError(t, err)

	groupVersions := metav1.ExtractGroupVersions(apiGroupList)
	if !reflect.DeepEqual(groupVersions, []string{"v1", "extensions/v1beta1"}) {
		t.Errorf("expected: %q, got: %q", []string{"v1", "extensions/v1beta1"}, groupVersions)
	}

	cluster2Client := client.Cluster(logicalcluster.New("cluster2"))
	apiGroupList, err = cluster2Client.ServerGroups()
	require.NoError(t, err)

	groupVersions = metav1.ExtractGroupVersions(apiGroupList)
	if !reflect.DeepEqual(groupVersions, []string{"v2", "extensions/v2beta2"}) {
		t.Errorf("expected: %q, got: %q", []string{"v2", "extensions/v2beta2"}, groupVersions)
	}
}
