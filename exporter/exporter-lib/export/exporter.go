package export

import (
	"log"

	"github.com/RHEcosystemAppEng/SaaSi/exporter/exporter-lib/config"
	"github.com/RHEcosystemAppEng/SaaSi/exporter/exporter-lib/connect"
	"github.com/RHEcosystemAppEng/SaaSi/exporter/exporter-lib/export/app"
	"github.com/RHEcosystemAppEng/SaaSi/exporter/exporter-lib/export/infra"
)

type Exporter struct {
	config        *config.Config
	appExporter   *app.AppExporter
	infraExporter *infra.InfraExporter
}

func NewExporterFromConfig(config *config.Config) *Exporter {
	exporter := Exporter{config: config}
	return &exporter
}

func (e *Exporter) Export() {
	connectionStatus := connect.ConnectCluster(&e.config.Exporter.Cluster)
	if connectionStatus.Error != nil {
		log.Fatalf("Cannot connect to given cluster: %s", connectionStatus.Error)
	}

	e.infraExporter = infra.NewInfraExporterFromConfig(e.config, connectionStatus)
	e.appExporter = app.NewAppExporterFromConfig(e.config, connectionStatus)

	e.infraExporter.Export()
	e.appExporter.Export()
}
