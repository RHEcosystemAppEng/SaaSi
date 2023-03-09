package main

import (
	config "github.com/RHEcosystemAppEng/SaaSi/exporter/exporter-lib/config"
	export "github.com/RHEcosystemAppEng/SaaSi/exporter/exporter-lib/export"
	"github.com/kr/pretty"
)

func main() {
	config := config.ReadConfig()
	pretty.Printf("Export configuration %# v", config)

	exporter := export.NewExporterFromConfig(config)
	exporter.Export()
}
