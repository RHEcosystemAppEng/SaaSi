package main

import (
	"github.com/RHEcosystemAppEng/SaaSi/deployer/deployer-lib/config"
	"github.com/RHEcosystemAppEng/SaaSi/deployer/deployer-lib/utils"
	"github.com/RHEcosystemAppEng/SaaSi/deployer/infra-deployer/pkg/api"
)

func main() {

	// get environment variables
	args := config.ParseEnvs()

	// set up logger based on debug parameter
	logger := utils.GetLogger(args.Debug)

	// print runtime configuration
	utils.PrettyPrint(logger, "Runtime configuration: %s", args)

	// handle HTTP requests
	api.HandleRequests(args, logger)
}
