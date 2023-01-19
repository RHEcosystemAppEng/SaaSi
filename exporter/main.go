package main

import (
	"log"
	"os"

	config "github.com/RHEcosystemAppEng/SaaSi/exporter/pkg/config"
	"github.com/RHEcosystemAppEng/SaaSi/exporter/pkg/installer"
	"github.com/kr/pretty"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Expected 1 argument, got ", len(os.Args)-1)
	}

	appConfig := config.ReadConfig(os.Args[1])
	pretty.Printf("Exporting application %# v", appConfig)
	exporterConfig := config.NewInstallerConfigFromApplicationConfig(appConfig)

	clusterRolesInspector := installer.NewClusterRolesInspector()
	clusterRolesInspector.LoadClusterRoles()

	exporter := installer.NewExporterFromConfig(appConfig, exporterConfig)
	exporter.PrepareOutput()
	exporter.ExportWithCrane()

	parametrizer := installer.NewParametrizerFromConfig(appConfig, exporterConfig)
	parametrizer.ExposeParameters()

	installer := installer.NewInstallerFromConfig(appConfig, exporterConfig, clusterRolesInspector)
	installer.BuildKustomizeInstaller()
}
