package main

import (
	config "github.com/RHEcosystemAppEng/SaaSi/exporter/exporter-lib/config"
	export "github.com/RHEcosystemAppEng/SaaSi/exporter/exporter-lib/export"
	"github.com/kr/pretty"
)

func main() {
	config := config.ReadConfig()
	pretty.Printf("Runtime configuration %# v", config)
	exporterConfig := config.ReadExporterConfig()
	pretty.Printf("Export configuration %# v", exporterConfig)

	exporter := export.NewExporterFromConfig(config)
	exporter.Export(exporterConfig)
}
