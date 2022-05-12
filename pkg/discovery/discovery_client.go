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
	"time"

	"github.com/kcp-dev/logicalcluster"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/kubernetes/scheme"
	restclient "k8s.io/client-go/rest"
)

const (
	// defaultTimeout is the maximum amount of time per request when no timeout has been set on a RESTClient.
	// Defaults to 32s in order to have a distinguishable length of time, relative to other timeouts that exist.
	defaultTimeout = 32 * time.Second
)

// ClusterDiscoveryClient holds an unscoped REST config and client
type ClusterDiscoveryClient struct {
	restClient *restclient.RESTClient
	config     *restclient.Config
}

// NewClusterDiscoveryClientForConfig takes a rest config and returns a ClusterDiscoveryClient
// that can create scoped discovery clients for logical clusters
func NewClusterDiscoveryClientForConfig(c *restclient.Config) (*ClusterDiscoveryClient, error) {
	err := setDiscoveryDefaults(c)
	if err != nil {
		return nil, err
	}

	httpClient, err := restclient.HTTPClientFor(c)
	if err != nil {
		return nil, err
	}
	client, err := restclient.UnversionedRESTClientForConfigAndClient(c, httpClient)
	if err != nil {
		return nil, err
	}

	return &ClusterDiscoveryClient{config: c, restClient: client}, nil
}

// Cluster creates a discovery client scoped to the specified logical cluster
func (c *ClusterDiscoveryClient) Cluster(cluster logicalcluster.Name) discovery.DiscoveryInterface {
	scopedConfig := restclient.CopyConfig(c.config)

	scopedConfig.Host = scopedConfig.Host + "/" + cluster.Path()

	// This shouldn't be able to panic
	return discovery.NewDiscoveryClientForConfigOrDie(scopedConfig)
}

// setDiscoveryDefaults sets sane defaults in the provided REST config to allow API discovery to succeed
func setDiscoveryDefaults(config *restclient.Config) error {
	config.APIPath = ""
	config.GroupVersion = nil
	if config.Timeout == 0 {
		config.Timeout = defaultTimeout
	}
	if config.Burst == 0 && config.QPS < 100 {
		// discovery is expected to be bursty, increase the default burst
		// to accommodate looking up resource info for many API groups.
		// matches burst set by ConfigFlags#ToDiscoveryClient().
		// see https://issue.k8s.io/86149
		config.Burst = 100
	}
	codec := runtime.NoopEncoder{Decoder: scheme.Codecs.UniversalDecoder()}
	config.NegotiatedSerializer = serializer.NegotiatedSerializerWrapper(runtime.SerializerInfo{Serializer: codec})
	if len(config.UserAgent) == 0 {
		config.UserAgent = restclient.DefaultKubernetesUserAgent()
	}
	return nil
}
