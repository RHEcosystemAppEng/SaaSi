package main

import (
	config "github.com/RHEcosystemAppEng/SaaSi/exporter/pkg/config"
	export "github.com/RHEcosystemAppEng/SaaSi/exporter/pkg/export"
	"github.com/kr/pretty"
)

func main() {
	config := config.ReadConfig()
	pretty.Printf("Export configuration %# v", config)

	exporter := export.NewExporterFromConfig(config)
	exporter.Export()
}
