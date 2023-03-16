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

type InfraExporterOutput struct {
	ClusterId    string `json:"clusterId"`
	Status       string `json:"status"`
	ErrorMessage string `json:"errorMessage"`
	Location     string `json:"location"`
}

func NewInfraExporterFromConfig(config *config.Config, clusterConfig *config.ClusterConfig, connectionStatus *connect.ConnectionStatus, logger *logrus.Logger) *InfraExporter {
	exporter := InfraExporter{infraContext: NewInfraContextFromConfig(config, clusterConfig, connectionStatus, logger)}

	return &exporter
}

func (e *InfraExporter) Export() InfraExporterOutput {
	output := InfraExporterOutput{ClusterId: e.infraContext.clusterConfig.ClusterId}
	err := utils.RunCommandAndLog(e.infraContext.Logger(), e.infraContext.ExportScript, "-k", e.infraContext.KubeConfigPath(),
		"-i", e.infraContext.clusterConfig.ClusterId, "-r", e.infraContext.ClusterFolder)
	if err != nil {
		output.ErrorMessage = err.Error()
		output.Status = utils.Failed.String()
	} else {
		output.Status = utils.Ok.String()
		output.Location = e.infraContext.ClusterFolder
	}
	return output
}
