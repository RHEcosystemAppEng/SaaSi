package main

import (
	"log"
	"os"

	config "github.com/RHEcosystemAppEng/SaaSi/exporter/pkg/config"
	"github.com/RHEcosystemAppEng/SaaSi/exporter/pkg/export"
	"github.com/kr/pretty"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Expected 1 argument, got ", len(os.Args)-1)
	}

	exporterConfig := config.ReadConfig(os.Args[1])
	pretty.Printf("Exporting application %# v", exporterConfig)
	context := export.NewContextFromConfig(exporterConfig)

	clusterRolesInspector := export.NewClusterRolesInspector()
	clusterRolesInspector.LoadClusterRoles()

	exporter := export.NewExporterFromConfig(&exporterConfig.Exporter.Application, context)
	exporter.PrepareOutput()
	exporter.ExportWithCrane()

	parametrizer := export.NewParametrizerFromConfig(&exporterConfig.Exporter.Application, context)
	parametrizer.ExposeParameters()

	installer := export.NewInstallerFromConfig(&exporterConfig.Exporter.Application, context, clusterRolesInspector)
	installer.BuildKustomizeInstaller()
}
