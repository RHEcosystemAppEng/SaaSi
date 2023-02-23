package main

import (
	"github.com/RHEcosystemAppEng/SaaSi/deployer/pkg/config"
	"github.com/RHEcosystemAppEng/SaaSi/deployer/pkg/connect"
	"github.com/RHEcosystemAppEng/SaaSi/deployer/pkg/context"
	"github.com/RHEcosystemAppEng/SaaSi/deployer/pkg/deployer/app/deployer"
	"github.com/RHEcosystemAppEng/SaaSi/deployer/pkg/deployer/app/packager"
	"github.com/RHEcosystemAppEng/SaaSi/deployer/pkg/deployer/infra/provisioner"
	"github.com/RHEcosystemAppEng/SaaSi/deployer/pkg/utils"
	"github.com/kr/pretty"
	"log"
	"reflect"
)

func main() {

	// Unmarshal deployer config and get cluster and application configs


	var deployApp bool = true
 	componentConfig := config.InitDeployerConfig()
	clusterConfig := componentConfig.ClusterConfig
	// connect to cluster
	kubeConnection := connect.ConnectToCluster(clusterConfig,true)

	// create deployer context to hold global variables
	deployerContext := context.InitDeployerContext(componentConfig.FlagArgs, kubeConnection)

	//Check if a cluster has been requested
	if !reflect.ValueOf(clusterConfig).IsZero(){
		// If there is no existing cluster, need to provision a new one, and postpone the deployment of the application to when the cluster will be ready.
        if reflect.ValueOf(clusterConfig.Server).IsZero() &&
			reflect.ValueOf(clusterConfig.Token).IsZero() &&
			reflect.ValueOf(clusterConfig.User).IsZero(){
			deployApp = false
			infraContext := context.InitInfraContext()
			NewClusterDetails := provisioner.ProvisionCluster(infraContext, &clusterConfig.Params,clusterConfig.Aws, deployerContext.GetRootSourceDir())


			log.Printf("Successfully deployed a cluster, clusterId:  %s ",clusterConfig.ClusterId)
			log.Printf("returned details of provisioned cluster with id : %s,  %s" ,clusterConfig.ClusterId, NewClusterDetails)
			// If requested to deploy also the application on new cluster, then need to update kubeconfig with new details
			if !reflect.ValueOf(componentConfig.ApplicationConfig).IsZero(){
				
				newClusterConfig := config.ClusterConfig{
					Server:        "",
					User:          "",
					Token:         "",
					FromClusterId: "",
					ClusterId:     "",
					Aws:           config.AwsSettings{},
					Params:        config.ClusterParams{},
				}
				newClusterKubeConnection := connect.ConnectToCluster(newClusterConfig,false)
				deployerContext.KubeConnection = newClusterKubeConnection
			}

		}
	}
	 if deployApp {
		pretty.Printf("Deploying the following configuration: \n%# v", componentConfig)


		// check if application deployment has been requested
		if !reflect.ValueOf(componentConfig.ApplicationConfig).IsZero() {

			// create application deployment package
			applicationPkg := packager.NewApplicationPkg(componentConfig.ApplicationConfig, deployerContext)

			// check if all mandatory variables have been set, else list unset vars and throw exception
			if len(applicationPkg.UnsetMandatoryParams) > 0 {
				log.Fatalf("ERROR: Please complete missing configuration for the following mandatory parameters (<FILEPATH>: <MANDATORY_PARAMETERS>.)\n%s", utils.StringifyMap(applicationPkg.UnsetMandatoryParams))
			}

			// deploy application deployment package
			deployer.DeployApplication(applicationPkg)

		} else {
			log.Println("No application to deploy")
		}
		//Need to wait for Cluster to be provisioned
	 } else {
		 log.Println("Mock application deployment after cluster provisioning")
	 }

}