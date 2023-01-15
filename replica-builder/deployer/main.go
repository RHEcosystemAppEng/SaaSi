package main

import (
	"os"
	"log"

	"github.com/kr/pretty"
	"github.com/RHEcosystemAppEng/SaaSi/replica-builder/deployer/pkg/packager"
	"github.com/RHEcosystemAppEng/SaaSi/replica-builder/deployer/pkg/deployer"
	"github.com/RHEcosystemAppEng/SaaSi/replica-builder/deployer/pkg/config"
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
	applicationPkg := packager.NewApplicationPkg(applicationConfig)

	// deploy deployment package
	deployer.NewDeployment(applicationPkg)
}
