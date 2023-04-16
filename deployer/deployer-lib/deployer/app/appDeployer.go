package app

import (
	"fmt"

	"github.com/RHEcosystemAppEng/SaaSi/deployer/deployer-lib/config"
	"github.com/RHEcosystemAppEng/SaaSi/deployer/deployer-lib/connect"
	"github.com/RHEcosystemAppEng/SaaSi/deployer/deployer-lib/context"
	"github.com/RHEcosystemAppEng/SaaSi/deployer/deployer-lib/deployer/app/deployer"
	"github.com/RHEcosystemAppEng/SaaSi/deployer/deployer-lib/deployer/app/packager"
	"github.com/RHEcosystemAppEng/SaaSi/deployer/deployer-lib/utils"
	"github.com/sirupsen/logrus"
)

var err error

type ApplicationOutput struct {
	ApplicationName string `json:"applicationName"`
	Status          string `json:"status"`
	ErrorMessage    string `json:"errorMessage"`
	Location        string `json:"location"`
}

func Deploy(componentConfig config.ComponentConfig, args *config.Args, logger *logrus.Logger) *ApplicationOutput {
	applicationName := componentConfig.ApplicationConfig.Name

	// validate deployemnt data
	err = componentConfig.Validate()
	if err != nil {
		return &ApplicationOutput{
			ApplicationName: applicationName,
			Status:          utils.Failed.String(),
			ErrorMessage:    fmt.Sprintf("Invalid configuration: %s", err.Error()),
			Location:        "",
		}
	}

	// connect to cluster
	kubeConnection := connect.ConnectToCluster(componentConfig.ClusterConfig, logger)
	if kubeConnection.Error != nil {
		return &ApplicationOutput{
			ApplicationName: applicationName,
			Status:          utils.Failed.String(),
			ErrorMessage:    fmt.Sprintf("Cannot connect to given cluster: %s", kubeConnection.Error.Error()),
			Location:        "",
		}
	}

	// create deployer context to hold global deployment parameters
	deployerContext := context.InitDeployerContext(args, kubeConnection, logger)
	if componentConfig.ClusterConfig.Provision.Provisioned {
		deployerContext.KubeConnection.KubeConfigPath = componentConfig.ClusterConfig.Provision.KubeConfigPath
	}

	// create application deployment package
	applicationPkg := packager.NewApplicationPkg(componentConfig.ApplicationConfig, deployerContext)
	if applicationPkg.Error != nil {
		return &ApplicationOutput{
			ApplicationName: applicationName,
			Status:          utils.Failed.String(),
			ErrorMessage:    fmt.Sprintf("Failed to create application deployment package: %s", applicationPkg.Error.Error()),
			Location:        "",
		}
	}

	// check if all mandatory variables have been set, else list unset vars and throw exception
	if len(applicationPkg.UnsetMandatoryParams) > 0 {
		UnsetMandatoryParamsMsg := fmt.Sprintf("Missing configuration for the following mandatory parameters (<FILEPATH>: <MANDATORY_PARAMETERS>.)\n%s", utils.StringifyMap(applicationPkg.UnsetMandatoryParams))
		return &ApplicationOutput{
			ApplicationName: applicationName,
			Status:          utils.Failed.String(),
			ErrorMessage:    UnsetMandatoryParamsMsg,
			Location:        "",
		}
	}

	// deploy application deployment package
	err = deployer.DeployApplication(applicationPkg)
	if err != nil {
		return &ApplicationOutput{
			ApplicationName: applicationName,
			Status:          utils.Failed.String(),
			ErrorMessage:    fmt.Sprintf("Failed to deploy application deployment package: %s", err.Error()),
			Location:        "",
		}
	}

	return &ApplicationOutput{
		ApplicationName: applicationName,
		Status:          utils.Ok.String(),
		ErrorMessage:    "",
		Location:        applicationPkg.UuidDir,
	}
}
