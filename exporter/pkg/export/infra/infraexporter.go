package infra

import (
	"github.com/RHEcosystemAppEng/SaaSi/exporter/pkg/config"
	"github.com/RHEcosystemAppEng/SaaSi/exporter/pkg/connect"
)

type InfraExporter struct {
	infraContext *InfraContext
}

func NewInfraExporterFromConfig(config *config.Config, connectionStatus *connect.ConnectionStatus) *InfraExporter {
	exporter := InfraExporter{infraContext: NewInfraContextFromConfig(config, connectionStatus)}

	return &exporter
}

func (i *InfraExporter) Export() {
}
