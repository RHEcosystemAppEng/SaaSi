package main

import (
	"log"
	"os"

	config "github.com/RHEcosystemAppEng/SaaSi/exporter/pkg/config"
	export "github.com/RHEcosystemAppEng/SaaSi/exporter/pkg/export"
	"github.com/kr/pretty"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Expected 1 argument, got ", len(os.Args)-1)
	}

	config := config.ReadConfig(os.Args[1])
	pretty.Printf("Export configuration %# v", config)

	exporter := export.NewExporterFromConfig(config)
	exporter.Export()
}
