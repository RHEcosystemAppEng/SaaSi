package infra

import (
	"github.com/RHEcosystemAppEng/SaaSi/exporter/pkg/config"
	"github.com/RHEcosystemAppEng/SaaSi/exporter/pkg/connect"
)

type InfraExporter struct {
	clusterConfig *config.ClusterConfig
}

func NewInfraExporterFromConfig(config *config.Config) *InfraExporter {
	exporter := InfraExporter{clusterConfig: &config.Exporter.Cluster}

	return &exporter
}

func (i *InfraExporter) Export(connectionStatus *connect.ConnectionStatus) {
}
