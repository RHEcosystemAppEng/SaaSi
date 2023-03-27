package main

import (
	"encoding/json"
	"reflect"

	"github.com/RHEcosystemAppEng/SaaSi/deployer/deployer-lib/config"
	"github.com/RHEcosystemAppEng/SaaSi/deployer/deployer-lib/context"
	"github.com/RHEcosystemAppEng/SaaSi/deployer/deployer-lib/deployer/app"
	"github.com/RHEcosystemAppEng/SaaSi/deployer/deployer-lib/deployer/infra/provisioner"
	"github.com/RHEcosystemAppEng/SaaSi/deployer/deployer-lib/utils"
	"github.com/kr/pretty"
)

var (
	err error
)

func main() {

	// get flag variables
	args := config.ParseFlags()

	// set up logger based on debug parameter
	logger := utils.GetLogger(args.Debug)

	// print runtime configuration
	utils.PrettyPrint(logger, "Runtime configuration: %s", args)

	// unmarshal deployer config and get cluster and application configs
	componentConfig := config.InitDeployerConfig(args.ConfigFile)

	// print deployment configuration
	utils.PrettyPrint(logger, "Deployment configuration: %s", componentConfig)

	clusterConfig := componentConfig.ClusterConfig
	clusterConfig.Provision.AuthByCreds = false
	clusterConfig.Provision.Provisioned = false
	clusterConfig.Provision.KubeConfigPath = ""
	//Check if a cluster has been requested
	if !reflect.ValueOf(clusterConfig).IsZero() {
		// If there is no existing cluster, need to provision a new one, and postpone the deployment of the application to when the cluster will be ready.
		if reflect.ValueOf(clusterConfig.Server).IsZero() &&
			reflect.ValueOf(clusterConfig.Token).IsZero() &&
			reflect.ValueOf(clusterConfig.User).IsZero() {
			//deployApp = false
			infraContext := context.InitInfraContext(args)
			beautifiedConfig, err := json.MarshalIndent(clusterConfig.Params, "", "   ")
			if err != nil {
				return
			}

			logger.Infof("About to deploy a cluster, clusterId:  %s , with following configuration (Every field that is not populated will be defaulted from source cluster): \n", clusterConfig.ClusterId)
			pretty.Printf("%s \n", string(beautifiedConfig))
			newClusterDetails := provisioner.ProvisionCluster(infraContext, &clusterConfig.Params, clusterConfig.Aws, args.RootSourceDir)

			logger.Infof("Successfully deployed a cluster, clusterId:  %s ", clusterConfig.ClusterId)
			logger.Infof("returned details of provisioned cluster with id : %s,  %+v\n", clusterConfig.ClusterId, newClusterDetails)
			// If requested to deploy also the application on new cluster, then need to update kubeconfig with new details

			if !reflect.ValueOf(componentConfig.ApplicationConfig).IsZero() {
				logger.Info("Applying new Cluster address and Credentials , and kubeconfig file to be used by application deployment")
				clusterConfig.Provision.Provisioned = true
				clusterConfig.Server = newClusterDetails.ApiServer
				clusterConfig.User = newClusterDetails.User
				clusterConfig.Token = newClusterDetails.Password
				clusterConfig.Provision.KubeConfigPath = newClusterDetails.KubeConfigPath
				clusterConfig.Provision.AuthByCreds = true

			}

		}
	} else {
		logger.Info("No cluster to deploy")
	}

	// check if application deployment has been requested
	if !reflect.ValueOf(componentConfig.ApplicationConfig).IsZero() {

		// package and deploy application
		output := app.Deploy(*componentConfig, args, logger)

		// report process result
		utils.PrettyPrint(logger, "Output: %s", output)

	} else {
		logger.Info("No application to deploy")
	}
}
