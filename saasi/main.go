package main

import (
	config "github.com/RHEcosystemAppEng/SaaSi/exporter/exporter-lib/config"
	export "github.com/RHEcosystemAppEng/SaaSi/exporter/exporter-lib/export"
	"github.com/RHEcosystemAppEng/SaaSi/exporter/exporter-lib/export/utils"
)

func main() {
	config := config.ReadConfigFromFlags()
	utils.PrettyPrint(config.Logger, "Runtime configuration: %s", config)
	exporterConfig := config.ReadExporterConfig()
	utils.PrettyPrint(config.Logger, "Export configuration: %s", exporterConfig)

	exporter := export.NewExporterFromConfig(config)
	output := exporter.Export(exporterConfig)
	utils.PrettyPrint(config.Logger, "Output: %s", output)
}
