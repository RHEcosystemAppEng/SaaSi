package main

import (
	"os"
	"log"

	"github.com/RHEcosystemAppEng/SaaSi/replica-builder/replica-installer/pkg/packager"
	"github.com/RHEcosystemAppEng/SaaSi/replica-builder/replica-installer/pkg/config"
	"github.com/kr/pretty"
)

var (
	err error
)

func main() {

	// get application config yaml as input
	if len(os.Args) != 2 {
		log.Fatal("Expected 1 argument, got ", len(os.Args)-1)
	}

	// init ApplicationConfig object
	applicationConfig := config.ReadApplicationConfig(os.Args[1])
	pretty.Printf("Exporting application %# v", applicationConfig)

	// create deployment package
	_ = packager.NewApplicationPkg(applicationConfig)
}
