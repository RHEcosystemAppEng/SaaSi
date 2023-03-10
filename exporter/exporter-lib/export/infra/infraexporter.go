package infra

import (
	"github.com/RHEcosystemAppEng/SaaSi/exporter/exporter-lib/config"
	"github.com/RHEcosystemAppEng/SaaSi/exporter/exporter-lib/connect"
	"github.com/RHEcosystemAppEng/SaaSi/exporter/exporter-lib/export/utils"
	"github.com/sirupsen/logrus"
)

type InfraExporter struct {
	infraContext *InfraContext
}

func NewInfraExporterFromConfig(config *config.Config, exporterConfig *config.ExporterConfig, connectionStatus *connect.ConnectionStatus, logger *logrus.Logger) *InfraExporter {
	exporter := InfraExporter{infraContext: NewInfraContextFromConfig(config, exporterConfig, connectionStatus, logger)}

	return &exporter
}

func (e *InfraExporter) Export() error {
	return utils.RunCommandAndLog(e.infraContext.Logger(), e.infraContext.ExportScript, "-k", e.infraContext.KubeConfigPath(),
		"-i", e.infraContext.clusterConfig.ClusterId, "-r", e.infraContext.ClusterFolder)
}
