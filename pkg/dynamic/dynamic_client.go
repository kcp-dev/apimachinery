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

package dynamic

import (
	"github.com/kcp-dev/logicalcluster"

	"k8s.io/client-go/dynamic"
	restclient "k8s.io/client-go/rest"
)

// ClusterDynamicClient holds an unscoped REST config and client
type ClusterDynamicClient struct {
	client *restclient.RESTClient
	config *restclient.Config
}

// NewClusterDynamicClientForConfig takes a rest config and returns a ClusterDynamicClient
// that can create scoped dynamic clients for logical clusters
func NewClusterDynamicClientForConfig(config *restclient.Config) (*ClusterDynamicClient, error) {

	config = dynamic.ConfigFor(config)

	httpClient, err := restclient.HTTPClientFor(config)
	if err != nil {
		return nil, err
	}
	client, err := restclient.UnversionedRESTClientForConfigAndClient(config, httpClient)
	if err != nil {
		return nil, err
	}

	return &ClusterDynamicClient{config: config, client: client}, nil
}

// Cluster creates a dynamic client scoped to the specified logical cluster
// A cross-cluster client will be returned if logicalcluster.Wildcard is provided
func (c *ClusterDynamicClient) Cluster(cluster logicalcluster.Name) dynamic.Interface {
	scopedConfig := restclient.CopyConfig(c.config)

	scopedConfig.Host = scopedConfig.Host + "/" + cluster.Path()

	// The original rest config has already parsed as valid. Modifying the host to scope it to a valid
	// cluster should never result in an invalid URL, so there should be no real possibility of a panic here.
	return dynamic.NewForConfigOrDie(scopedConfig)
}
