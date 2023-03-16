package main

import (
	"fmt"

	"github.com/RHEcosystemAppEng/SaaSi/deployer/app-deployer/pkg/api"
	"github.com/RHEcosystemAppEng/SaaSi/deployer/deployer-lib/config"
)

func main() {
	envConfig := config.ParseEnvs()
	fmt.Print(envConfig) //PLACEHOLDER

	// ADD LOGGER

	api.HandleRequests()
}
