package main

import (
	config "github.com/RHEcosystemAppEng/SaaSi/exporter/exporter-lib/config"
	export "github.com/RHEcosystemAppEng/SaaSi/exporter/exporter-lib/export"
	"github.com/RHEcosystemAppEng/SaaSi/exporter/exporter-lib/export/utils"
)

func main() {
	config := config.ReadConfigFromFlags()
	logger := utils.GetLogger(config.Debug)
	utils.PrettyPrint(logger, "Runtime configuration: %s", config)
	exporterConfig := config.ReadExporterConfig()
	exporterConfig.InitializeForCLI()
	err := exporterConfig.Validate()
	if err != nil {
		logger.Fatalf("Invalid configuration: %s", err)
	}
	utils.PrettyPrint(logger, "Export configuration: %s", exporterConfig)

	exporter := export.NewExporterFromConfig(config)
	output := exporter.Export(exporterConfig)
	utils.PrettyPrint(logger, "Output: %s", output)
}
