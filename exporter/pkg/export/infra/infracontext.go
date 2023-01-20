package infra

import (
	"github.com/RHEcosystemAppEng/SaaSi/exporter/pkg/config"
	"github.com/RHEcosystemAppEng/SaaSi/exporter/pkg/connect"
	"k8s.io/client-go/rest"
)

type InfraContext struct {
	clusterConfig    *config.ClusterConfig
	connectionStatus *connect.ConnectionStatus
	scriptFolder     string
}

func NewInfraContextFromConfig(config *config.Config, connectionStatus *connect.ConnectionStatus) *InfraContext {
	context := InfraContext{clusterConfig: &config.Exporter.Cluster, connectionStatus: connectionStatus}
	return &context
}

func (c *InfraContext) KubeConfig() *rest.Config {
	return c.connectionStatus.KubeConfig
}

func (c *InfraContext) KubeConfigPath() string {
	return c.connectionStatus.KubeConfigPath
}
