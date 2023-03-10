package connect

import (
	"io/ioutil"

	"github.com/RHEcosystemAppEng/SaaSi/exporter/exporter-lib/config"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/version"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

type ConnectionStatus struct {
	KubeConfigPath string
	KubeConfig     *rest.Config
	Error          error
}

func ConnectCluster(clusterConfig *config.ClusterConfig, logger *logrus.Logger) *ConnectionStatus {
	status := ConnectionStatus{}

	status.KubeConfig = &rest.Config{}
	// kubeConfig.Username = i.clusterConfig.User
	status.KubeConfig.BearerToken = clusterConfig.Token
	status.KubeConfig.Host = clusterConfig.Server
	// kubeConfig.APIPath = i.clusterConfig.Server
	status.KubeConfig.Insecure = true

	status.KubeConfigPath, status.Error = generateKubeConfiguration(clusterConfig, logger)
	if status.Error == nil {
		logger.Infof("Connected to cluster %s at server %s", clusterConfig.ClusterId, clusterConfig.Server)

		var discoveryClient *discovery.DiscoveryClient
		discoveryClient, status.Error = discovery.NewDiscoveryClientForConfig(status.KubeConfig)
		if status.Error == nil {
			var version *version.Info
			version, status.Error = discoveryClient.ServerVersion()
			if status.Error == nil {
				logger.Infof("Connected to cluster with version: %s", version)
			}
		}
	}
	return &status
}

func generateKubeConfiguration(clusterConfig *config.ClusterConfig, logger *logrus.Logger) (string, error) {
	namespace := "default"
	clusters := make(map[string]*api.Cluster)
	clusters["default-cluster"] = &api.Cluster{
		Server:                clusterConfig.Server,
		InsecureSkipTLSVerify: true,
		// CertificateAuthorityData: secret.Data["ca.crt"],
	}

	contexts := make(map[string]*api.Context)
	contexts["default-context"] = &api.Context{
		Cluster:   "default-cluster",
		Namespace: namespace,
		AuthInfo:  namespace,
	}

	authinfos := make(map[string]*api.AuthInfo)
	authinfos[namespace] = &api.AuthInfo{
		Token: clusterConfig.Token,
	}

	clientConfig := api.Config{
		Kind:           "Config",
		APIVersion:     "v1",
		Clusters:       clusters,
		Contexts:       contexts,
		CurrentContext: "default-context",
		AuthInfos:      authinfos,
	}

	kubeconfig, err := ioutil.TempFile("/tmp", "config")
	if err == nil {
		clientcmd.WriteToFile(clientConfig, kubeconfig.Name())
		logger.Infof("Saved kubeconfig to %s", kubeconfig.Name())
	}
	return kubeconfig.Name(), err
}
