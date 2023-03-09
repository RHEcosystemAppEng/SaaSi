package infra

import (
	"github.com/RHEcosystemAppEng/SaaSi/exporter/exporter-lib/config"
	"github.com/RHEcosystemAppEng/SaaSi/exporter/exporter-lib/connect"
	"github.com/RHEcosystemAppEng/SaaSi/exporter/exporter-lib/export/utils"
)

type InfraExporter struct {
	infraContext *InfraContext
}

func NewInfraExporterFromConfig(config *config.Config, connectionStatus *connect.ConnectionStatus) *InfraExporter {
	exporter := InfraExporter{infraContext: NewInfraContextFromConfig(config, connectionStatus)}

	return &exporter
}

func (e *InfraExporter) Export() {
	utils.RunCommandAndLog(e.infraContext.ExportScript, "-k", e.infraContext.KubeConfigPath(),
		"-i", e.infraContext.clusterConfig.ClusterId, "-r", e.infraContext.ClusterFolder)
}
