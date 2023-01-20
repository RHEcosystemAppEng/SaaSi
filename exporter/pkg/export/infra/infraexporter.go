package infra

import (
	"log"

	"github.com/RHEcosystemAppEng/SaaSi/exporter/pkg/config"
	"github.com/RHEcosystemAppEng/SaaSi/exporter/pkg/connect"
	"github.com/RHEcosystemAppEng/SaaSi/exporter/pkg/export/utils"
)

type InfraExporter struct {
	infraContext *InfraContext
}

func NewInfraExporterFromConfig(config *config.Config, connectionStatus *connect.ConnectionStatus) *InfraExporter {
	exporter := InfraExporter{infraContext: NewInfraContextFromConfig(config, connectionStatus)}

	return &exporter
}

func (e *InfraExporter) Export() {
	log.Printf("Running infra exporter with context: %v", e.infraContext)
	utils.RunCommandAndLog(e.infraContext.ExportScript, "-k", e.infraContext.KubeConfigPath(),
		"-i", e.infraContext.clusterConfig.ClusterId, "-r", e.infraContext.ClusterFolder)
}
