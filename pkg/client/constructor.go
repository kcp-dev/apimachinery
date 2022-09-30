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
	"net/http"
	"sync"

	"github.com/kcp-dev/logicalcluster/v2"
	"k8s.io/client-go/rest"
)

type Constructor[R any] struct {
	NewForConfigAndClient func(*rest.Config, *http.Client) (R, error)
}

type Cache[R any] interface {
	ClusterOrDie(name logicalcluster.Name) R
	Cluster(name logicalcluster.Name) (R, error)
}

func NewCache[R any](cfg *rest.Config, client *http.Client, constructor *Constructor[R]) Cache[R] {
	return &clientCache[R]{
		cfg:         cfg,
		client:      client,
		constructor: constructor,

		RWMutex:          &sync.RWMutex{},
		clientsByCluster: map[logicalcluster.Name]R{},
	}
}

type clientCache[R any] struct {
	cfg         *rest.Config
	client      *http.Client
	constructor *Constructor[R]

	*sync.RWMutex
	clientsByCluster map[logicalcluster.Name]R
}

func (c *clientCache[R]) ClusterOrDie(name logicalcluster.Name) R {
	client, err := c.Cluster(name)
	if err != nil {
		// we ensure that the config is valid in the constructor, and we assume that any changes
		// we make to it during scoping will not make it invalid, in order to hide the error from
		// downstream callers (as it should forever be nil); this is slightly risky
		panic(err)
	}
	return client
}

func (c *clientCache[R]) Cluster(name logicalcluster.Name) (R, error) {
	c.RLock()
	if cachedClient, exists := c.clientsByCluster[name]; exists {
		return cachedClient, nil
	}
	c.RUnlock()

	cfg := SetCluster(rest.CopyConfig(c.cfg), name)
	instance, err := c.constructor.NewForConfigAndClient(cfg, c.client)
	if err != nil {
		var result R
		return result, err
	}

	c.Lock()
	defer c.Unlock()
	if cachedClient, exists := c.clientsByCluster[name]; exists {
		return cachedClient, nil
	}

	c.clientsByCluster[name] = instance

	return instance, nil
}
