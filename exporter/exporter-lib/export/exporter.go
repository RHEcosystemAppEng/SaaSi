package export

import (
	"fmt"

	"github.com/RHEcosystemAppEng/SaaSi/exporter/exporter-lib/config"
	"github.com/RHEcosystemAppEng/SaaSi/exporter/exporter-lib/connect"
	"github.com/RHEcosystemAppEng/SaaSi/exporter/exporter-lib/export/app"
	"github.com/RHEcosystemAppEng/SaaSi/exporter/exporter-lib/export/infra"
	"github.com/RHEcosystemAppEng/SaaSi/exporter/exporter-lib/export/utils"
	"github.com/sirupsen/logrus"
)

type Exporter struct {
	config        *config.Config
	appExporter   *app.AppExporter
	infraExporter *infra.InfraExporter
	logger        *logrus.Logger
}
type ExporterOutput struct {
	Status        string                    `json:"status"`
	ErrorMessage  string                    `json:"errorMessage"`
	AppExporter   app.AppExporterOutput     `json:"appExporter"`
	InfraExporter infra.InfraExporterOutput `json:"infraExporter"`
}

func NewExporterFromConfig(config *config.Config) *Exporter {
	exporter := Exporter{config: config}
	exporter.logger = utils.GetLogger(config.Debug)
	return &exporter
}

func (e *Exporter) Export(exporterConfig *config.ExporterConfig) ExporterOutput {
	output := ExporterOutput{}
	err := exporterConfig.Validate()
	if err != nil {
		output.Status = utils.Failed.String()
		output.ErrorMessage = err.Error()
		return output
	}
	connectionStatus := connect.ConnectCluster(&exporterConfig.Cluster, e.logger)
	if connectionStatus.Error != nil {
		message := fmt.Sprintf("Cannot connect to given cluster: %s", connectionStatus.Error)
		e.logger.Errorf(message)
		output.Status = utils.Failed.String()
		output.ErrorMessage = message
		return output
	}

	e.infraExporter = infra.NewInfraExporterFromConfig(e.config, exporterConfig, connectionStatus, e.logger)
	e.appExporter = app.NewAppExporterFromConfig(e.config, exporterConfig, connectionStatus, e.logger)

	infraExporterOutput := e.infraExporter.Export()
	output.InfraExporter = infraExporterOutput
	if infraExporterOutput.Status == utils.Failed.String() {
		output.Status = utils.Failed.String()
		return output
	}

	appExporterOutput := e.appExporter.Export()
	output.AppExporter = appExporterOutput
	output.Status = appExporterOutput.Status
	return output
}
