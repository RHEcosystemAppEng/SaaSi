package connect

import (
	"io/ioutil"
	"log"

	"github.com/RHEcosystemAppEng/SaaSi/deployer/deployer-lib/config"
	"k8s.io/apimachinery/pkg/version"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

const (
	DEFAULT_NAMESPACE_NAME    = "default"
	DEFAULT_AUTH_INFO         = "default"
	DEFAULT_CLUSTER_NAME      = "default-cluster"
	DEFAULT_CONTEXT_NAME      = "default-context"
	CLIENT_CONFIG_KIND        = "Config"
	CLIENT_CONFIG_API_VERSION = "v1"
	TEMPORARY_FILE_DIR        = "/tmp"
	KUBE_CONFIG_FILE_NAME     = "config"
)

var (
	err             error
	discoveryClient *discovery.DiscoveryClient
	versionInfo     *version.Info
)

type KubeConnection struct {
	KubeConfig     *rest.Config
	KubeConfigPath string
}

func ConnectToCluster(clusterConfig config.ClusterConfig, authByToken bool) *KubeConnection {
	// init kube connection
	conn := KubeConnection{}

	// bind kubeconfig to kube connection
	conn.KubeConfig = &rest.Config{}

	// set credentials in kube config
	conn.KubeConfig.Host = clusterConfig.Server
	if authByToken {
		conn.KubeConfig.BearerToken = clusterConfig.Token
	} else {
		conn.KubeConfig.Password = clusterConfig.Token
		conn.KubeConfig.Username = clusterConfig.User
	}
	conn.KubeConfig.Insecure = true

	// generate kube config
	conn.generateKubeConfiguration(authByToken)

	// discover supported resources in the api server
	discoveryClient, err = discovery.NewDiscoveryClientForConfig(conn.KubeConfig)
	if err != nil {
		log.Fatalf("Cannot connect to given cluster: %s", err)
	}

	// retrieve and parse the servers version
	versionInfo, err = discoveryClient.ServerVersion()
	if err != nil {
		log.Fatalf("Cannot connect to given cluster: %s", err)
	}

	log.Printf("Connected to cluster %s at server %s with version %s", clusterConfig.ClusterId, clusterConfig.Server, versionInfo)

	return &conn
}

func (conn *KubeConnection) generateKubeConfiguration(authByToken bool) {

	// define cluster configuration
	clusters := make(map[string]*api.Cluster)
	clusters[DEFAULT_CLUSTER_NAME] = &api.Cluster{
		Server:                conn.KubeConfig.Host,
		InsecureSkipTLSVerify: true,
		// CertificateAuthorityData: secret.Data["ca.crt"],
	}

	// define context configuration
	contexts := make(map[string]*api.Context)
	contexts[DEFAULT_CONTEXT_NAME] = &api.Context{
		Cluster:   DEFAULT_CLUSTER_NAME,
		Namespace: DEFAULT_NAMESPACE_NAME,
		AuthInfo:  DEFAULT_AUTH_INFO,
	}

	// define auth info configuration
	var authinfos map[string]*api.AuthInfo
	if authByToken {
		authinfos = make(map[string]*api.AuthInfo)
		authinfos[DEFAULT_AUTH_INFO] = &api.AuthInfo{
			Token: conn.KubeConfig.BearerToken,
		}
		//otherwise, auth by basic authentication
	} else {
		authinfos = make(map[string]*api.AuthInfo)
		authinfos[DEFAULT_AUTH_INFO] = &api.AuthInfo{
			Username: conn.KubeConfig.Username,
			Password: conn.KubeConfig.Password,
		}
	}

	// define client config configuration
	clientConfig := api.Config{
		Kind:           CLIENT_CONFIG_KIND,
		APIVersion:     CLIENT_CONFIG_API_VERSION,
		Clusters:       clusters,
		Contexts:       contexts,
		CurrentContext: DEFAULT_CONTEXT_NAME,
		AuthInfos:      authinfos,
	}

	kubeconfig, err := ioutil.TempFile(TEMPORARY_FILE_DIR, KUBE_CONFIG_FILE_NAME)
	if err != nil {
		log.Fatalf("Cannot connect to given cluster: %s", err)
	}

	clientcmd.WriteToFile(clientConfig, kubeconfig.Name())
	log.Printf("Saved kubeconfig to %s", kubeconfig.Name())

	conn.KubeConfigPath = kubeconfig.Name()
}
