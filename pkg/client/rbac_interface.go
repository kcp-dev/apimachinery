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

	rbacv1 "k8s.io/client-go/kubernetes/typed/rbac/v1"
)

// ========== THE FOLLOWING CODE WOULD LIVE IN kcp-dev/apimachinery ==========

type client interface {
	RESTClient() rest.Interface
}

type ClientConstructor[R client] struct {
	NewForConfigAndClient func(*rest.Config, *http.Client) (R, error)
	New                   func(rest.Interface) R
}

type ClientCache[R client] interface {
	ClusterOrDie(name logicalcluster.Name) R
	Cluster(name logicalcluster.Name) (R, error)
}

func NewClientCache[R client](cfg *rest.Config, client *http.Client, constructor *ClientConstructor[R]) ClientCache[R] {
	return &clientCache[R]{
		cfg:         cfg,
		client:      client,
		constructor: constructor,

		RWMutex:          &sync.RWMutex{},
		clientsByCluster: map[logicalcluster.Name]rest.Interface{},
	}
}

type clientCache[R client] struct {
	cfg         *rest.Config
	client      *http.Client
	constructor *ClientConstructor[R]

	*sync.RWMutex
	clientsByCluster map[logicalcluster.Name]rest.Interface
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
		return c.constructor.New(cachedClient), nil
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
		return c.constructor.New(cachedClient), nil
	}

	c.clientsByCluster[name] = instance.RESTClient()

	return instance, nil
}

// ========== THE FOLLOWING CODE WOULD BE GENERATED FOR OTHERS' TYPES ==========

type RbacV1ClusterInterface interface {
	RbacV1ClusterScoper
	ClusterRolesClusterScoperGetter
	ClusterRoleBindingsClusterScoperGetter
	RoleClusterScoperGetter
	RoleBindingsClusterScoperGetter
}

type RbacV1ClusterScoper interface {
	Cluster(logicalcluster.Name) rbacv1.RbacV1Interface
}

type ClusterRolesClusterScoperGetter interface {
	ClusterRoles() ClusterRolesClusterScoper
}

type ClusterRolesClusterScoper interface {
	Cluster(logicalcluster.Name) rbacv1.ClusterRolesGetter
}

type RoleClusterScoperGetter interface {
	Roles() RoleClusterScoper
}

type RoleClusterScoper interface {
	Cluster(logicalcluster.Name) rbacv1.RolesGetter
}

type ClusterRoleBindingsClusterScoperGetter interface {
	ClusterRoleBindings() ClusterRoleBindingClusterScoper
}

type ClusterRoleBindingClusterScoper interface {
	Cluster(logicalcluster.Name) rbacv1.ClusterRoleBindingsGetter
}

type RoleBindingsClusterScoperGetter interface {
	RoleBindings() RoleBindingClusterScoper
}

type RoleBindingClusterScoper interface {
	Cluster(logicalcluster.Name) rbacv1.RoleBindingsGetter
}

// NewForConfig creates a new RbacV1Client for the given config.
// NewForConfig is equivalent to NewForConfigAndClient(c, httpClient),
// where httpClient was generated with rest.HTTPClientFor(c).
func NewForConfig(c *rest.Config) (*RbacV1ClusterClient, error) {
	client, err := rest.HTTPClientFor(c)
	if err != nil {
		return nil, err
	}
	return NewForConfigAndClient(c, client)
}

// NewForConfigAndClient creates a new RbacV1Client for the given config and http client.
// Note the http client provided takes precedence over the configured transport values.
func NewForConfigAndClient(c *rest.Config, h *http.Client) (*RbacV1ClusterClient, error) {
	cache := NewClientCache(c, h, &ClientConstructor[*rbacv1.RbacV1Client]{
		New:                   rbacv1.New,
		NewForConfigAndClient: rbacv1.NewForConfigAndClient,
	})
	if _, err := cache.Cluster(logicalcluster.New("root")); err != nil {
		return nil, err
	}
	return &RbacV1ClusterClient{clientCache: cache}, nil
}

// NewForConfigOrDie creates a new RbacV1ClusterClient for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *RbacV1ClusterClient {
	client, err := NewForConfig(c)
	if err != nil {
		panic(err)
	}
	return client
}

type RbacV1ClusterClient struct {
	clientCache ClientCache[*rbacv1.RbacV1Client]
}

func (c *RbacV1ClusterClient) Cluster(name logicalcluster.Name) rbacv1.RbacV1Interface {
	return c.clientCache.ClusterOrDie(name)
}

func (c *RbacV1ClusterClient) ClusterRoles() ClusterRolesClusterScoper {
	return &clusterRolesClusterClient{ClientCache: c.clientCache}
}

type clusterRolesClusterClient struct {
	ClientCache[*rbacv1.RbacV1Client]
}

func (c *clusterRolesClusterClient) Cluster(name logicalcluster.Name) rbacv1.ClusterRolesGetter {
	return c.ClientCache.ClusterOrDie(name)
}

func (c *RbacV1ClusterClient) ClusterRoleBindings() ClusterRoleBindingClusterScoper {
	return &clusterRoleBindingsClusterClient{ClientCache: c.clientCache}
}

type clusterRoleBindingsClusterClient struct {
	ClientCache[*rbacv1.RbacV1Client]
}

func (c *clusterRoleBindingsClusterClient) Cluster(name logicalcluster.Name) rbacv1.ClusterRoleBindingsGetter {
	return c.ClientCache.ClusterOrDie(name)
}

func (c *RbacV1ClusterClient) Roles() RoleClusterScoper {
	return &rolesClusterClient{ClientCache: c.clientCache}
}

type rolesClusterClient struct {
	ClientCache[*rbacv1.RbacV1Client]
}

func (c *rolesClusterClient) Cluster(name logicalcluster.Name) rbacv1.RolesGetter {
	return c.ClientCache.ClusterOrDie(name)
}

func (c *RbacV1ClusterClient) RoleBindings() RoleBindingClusterScoper {
	return &roleBindingsClusterClient{ClientCache: c.clientCache}
}

type roleBindingsClusterClient struct {
	ClientCache[*rbacv1.RbacV1Client]
}

func (c *roleBindingsClusterClient) Cluster(name logicalcluster.Name) rbacv1.RoleBindingsGetter {
	return c.ClientCache.ClusterOrDie(name)
}

var _ RbacV1ClusterInterface = (*RbacV1ClusterClient)(nil)
