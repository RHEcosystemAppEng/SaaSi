package infra

import (
	"fmt"
	"strings"

	"github.com/RHEcosystemAppEng/SaaSi/deployer/deployer-lib/config"
	"github.com/RHEcosystemAppEng/SaaSi/deployer/deployer-lib/context"
	"github.com/RHEcosystemAppEng/SaaSi/deployer/deployer-lib/deployer/infra/provisioner"
	"github.com/RHEcosystemAppEng/SaaSi/deployer/deployer-lib/utils"
	"github.com/sirupsen/logrus"
)

var err error

const DEFAULT_CLUSTER_ID = "unknown"

type InfrastructureOutput struct {
	ClusterId    string `json:"clusterId"`
	Status       string `json:"status"`
	ErrorMessage string `json:"errorMessage"`
	Location     string `json:"location"`
}

func Deploy(componentConfig config.ComponentConfig, args *config.Args, logger *logrus.Logger) *InfrastructureOutput {
	clusterId := DEFAULT_CLUSTER_ID

	// validate deployemnt data
	err = componentConfig.ValidateForInfraDeployment()
	if err != nil {
		if !strings.Contains(err.Error(), "missing clusterId configuration") {
			clusterId = componentConfig.ClusterConfig.ClusterId
		}

		return &InfrastructureOutput{
			ClusterId:    clusterId,
			Status:       utils.Failed.String(),
			ErrorMessage: fmt.Sprintf("Invalid configuration: %s", err.Error()),
			Location:     "",
		}
	}

	clusterId = componentConfig.ClusterConfig.ClusterId

	// create infra context to hold global build parameters
	infraContext := context.InitInfraContext(args, logger)
	logger.Infof("DEBUG - clusterId: %v", infraContext)

	// provision cluster
	newClusterDetails := provisioner.ProvisionCluster(infraContext, &componentConfig.ClusterConfig.Params, componentConfig.ClusterConfig.Aws)
	if newClusterDetails.Error != nil {
		return &InfrastructureOutput{
			ClusterId:    clusterId,
			Status:       utils.Failed.String(),
			ErrorMessage: fmt.Sprintf("Failed to provision cluster %s: %s", clusterId, err.Error()),
			Location:     "",
		}
	}

	return &InfrastructureOutput{
		ClusterId:    clusterId,
		Status:       utils.Ok.String(),
		ErrorMessage: "",
		Location:     "",
	}

}
