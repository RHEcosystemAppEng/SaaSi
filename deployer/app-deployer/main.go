package main

import (
	"fmt"

	"github.com/RHEcosystemAppEng/SaaSi/deployer/app-deployer/pkg/api"
	"github.com/RHEcosystemAppEng/SaaSi/deployer/deployer-lib/config"
)

func main() {
	args := config.ParseEnvs()
	fmt.Print(args) //PLACEHOLDER

	// ADD LOGGER

	api.HandleRequests(args)
}
