package main

import (
	"encoding/json"
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

	authByToken := true
 	componentConfig := config.InitDeployerConfig()
	pretty.Printf("Deploying the following configuration: \n%# v", componentConfig)
	clusterConfig := componentConfig.ClusterConfig
    var clusterProvisioned = false
	var kubeConfigPath string = ""
	//Check if a cluster has been requested
	if !reflect.ValueOf(clusterConfig).IsZero(){
		// If there is no existing cluster, need to provision a new one, and postpone the deployment of the application to when the cluster will be ready.
        if reflect.ValueOf(clusterConfig.Server).IsZero() &&
			reflect.ValueOf(clusterConfig.Token).IsZero() &&
			reflect.ValueOf(clusterConfig.User).IsZero(){
			//deployApp = false
			infraContext := context.InitInfraContext()
			beautifiedConfig, err := json.MarshalIndent(clusterConfig.Params, "", "   ")
			if err != nil {
				return
			}

			log.Printf("About to deploy a cluster, clusterId:  %s , with following configuration (Every field that is not populated will be defaulted from source cluster): \n" , clusterConfig.ClusterId)
			pretty.Printf("%s \n",string(beautifiedConfig))
			newClusterDetails := provisioner.ProvisionCluster(infraContext, &clusterConfig.Params,clusterConfig.Aws, componentConfig.FlagArgs.RootSourceDir)


			log.Printf("Successfully deployed a cluster, clusterId:  %s ",clusterConfig.ClusterId)
			log.Printf("returned details of provisioned cluster with id : %s,  %+v\n" ,clusterConfig.ClusterId, newClusterDetails)
			// If requested to deploy also the application on new cluster, then need to update kubeconfig with new details

			if !reflect.ValueOf(componentConfig.ApplicationConfig).IsZero(){
				log.Print("Applying new Cluster address and Credentials , and kubeconfig file to be used by application deployment")
				clusterProvisioned = true
				clusterConfig.Server = newClusterDetails.ApiServer
				clusterConfig.User = newClusterDetails.User
				clusterConfig.Token = newClusterDetails.Password
				kubeConfigPath = newClusterDetails.KubeConfigPath
				authByToken = false

			}

		}
	}
	 // connect to cluster
	 kubeConnection := connect.ConnectToCluster(clusterConfig, authByToken)
	 // create deployer context to hold global variables
	 deployerContext := context.InitDeployerContext(componentConfig.FlagArgs, kubeConnection)
	 if clusterProvisioned {
		 deployerContext.KubeConnection.KubeConfigPath = kubeConfigPath
	 }
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

}