package discovery

import (
	"path"
	"time"

	"github.com/kcp-dev/apimachinery/pkg/logicalcluster"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/kubernetes/scheme"
	restclient "k8s.io/client-go/rest"
)

type ClusterDiscoveryClient interface {
	Cluster(cluster logicalcluster.LogicalCluster) discovery.DiscoveryInterface
}

func NewClusterDiscoveryClientForConfig(c *restclient.Config) (ClusterDiscoveryClient, error) {

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

	return &clusterDiscoveryClient{config: c, restClient: client}, nil
}

type clusterDiscoveryClient struct {
	restClient *restclient.RESTClient
	config     *restclient.Config
}

func (c *clusterDiscoveryClient) Cluster(cluster logicalcluster.LogicalCluster) discovery.DiscoveryInterface {
	scopedConfig := restclient.CopyConfig(c.config)

	// This does not parse as a URL in the rest.DefaultServerURL function
	scopedConfig.Host = path.Join(scopedConfig.Host, cluster.Path())

	// This does nothing AFAICT
	// scopedConfig.APIPath = cluster.Path()

	return discovery.NewDiscoveryClientForConfigOrDie(scopedConfig)
}

func setDiscoveryDefaults(config *restclient.Config) error {
	//TODO
	defaultTimeout := 32 * time.Second
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
