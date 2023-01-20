package export

import (
	"log"

	"github.com/RHEcosystemAppEng/SaaSi/exporter/pkg/config"
	"github.com/RHEcosystemAppEng/SaaSi/exporter/pkg/connect"
	"github.com/RHEcosystemAppEng/SaaSi/exporter/pkg/export/app"
	"github.com/RHEcosystemAppEng/SaaSi/exporter/pkg/export/infra"
)

type Exporter struct {
	config        *config.Config
	appExporter   *app.AppExporter
	infraExporter *infra.InfraExporter
}

func NewExporterFromConfig(config *config.Config) *Exporter {
	exporter := Exporter{config: config}

	exporter.infraExporter = infra.NewInfraExporterFromConfig(config)
	exporter.appExporter = app.NewAppExporterFromConfig(config)

	return &exporter
}

func (e *Exporter) Export() {
	connectionStatus := connect.ConnectCluster(&e.config.Exporter.Cluster)
	if connectionStatus.Error != nil {
		log.Fatalf("Cannot connect to given cluster: %s", connectionStatus.Error)
	}

	e.infraExporter.Export(connectionStatus)
	e.appExporter.Export(connectionStatus)
}
