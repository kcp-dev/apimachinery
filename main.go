package main

import (
	"context"

	"github.com/kcp-dev/apimachinery/pkg/client"
	"github.com/kcp-dev/logicalcluster/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

func main() {
	// create an un-scoped client
	var cfg *rest.Config // loaded somehow
	unscopedClient, err := client.NewForConfig(cfg)
	if err != nil {
		panic(err)
	}

	// scope the whole group to a cluster, so it can be used in a single-cluster context
	rbacClient := unscopedClient.Cluster(logicalcluster.New("legacy"))
	// now, use it as normal
	rbacClient.RoleBindings("default").Get(context.TODO(), "system", metav1.GetOptions{})

	// defer scoping until later, keeping un-scoped but strongly typed clients for specific resources
	clusterRoleClient := unscopedClient.ClusterRoles()
	roleClient := unscopedClient.Roles()
	// use them *almost* as normal
	clusterRoleClient.Cluster(logicalcluster.New("whoa")).ClusterRoles().Get(context.TODO(), "foo", metav1.GetOptions{})
	// NOTE: for cluster-scoped resources we could technically make the above call not stutter on e.g. ClusterRoles():
	// clusterRoleClient.Cluster(logicalcluster.New("whoa")).Get(context.TODO(), "foo", metav1.GetOptions{})

	roleClient.Cluster(logicalcluster.New("dang")).Roles("default").Get(context.TODO(), "bar", metav1.GetOptions{})
	// NOTE: for namespace-scoped resources we could technically make the above call not stutter on e.g. Roles():
	// roleClient.Cluster(logicalcluster.New("dang")).Namespace("default").Get(context.TODO(), "foo", metav1.GetOptions{})
}
