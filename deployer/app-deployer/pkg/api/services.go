package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/RHEcosystemAppEng/SaaSi/deployer/deployer-lib/config"
	"github.com/RHEcosystemAppEng/SaaSi/deployer/deployer-lib/connect"
	"github.com/RHEcosystemAppEng/SaaSi/deployer/deployer-lib/context"
	"github.com/RHEcosystemAppEng/SaaSi/deployer/deployer-lib/deployer/app/deployer"
	"github.com/RHEcosystemAppEng/SaaSi/deployer/deployer-lib/deployer/app/packager"
	"github.com/RHEcosystemAppEng/SaaSi/deployer/deployer-lib/utils"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

const (
	CONTENT_TYPE     = "application/json"
	APPLICATION_NAME = "app-deployer"
	BUILD_VERSION    = "dev"
	STATUS           = "up"
)

type applicationInfo struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
	Status  string `yaml:"status"`
}

func deploy(args *config.Args, logger *logrus.Logger) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		// validate requested path
		if req.URL.Path != APP_DEPLOYER_PATH {
			http.Error(rw, fmt.Sprintf("Unmanaged path %s", req.URL.Path), http.StatusNotFound)
			http.NotFound(rw, req)
			return
		}

		// validate requested method
		if req.Method != POST {
			http.Error(rw, fmt.Sprintf("Expect method %s at %s, got %v", POST, APP_DEPLOYER_PATH, req.Method), http.StatusMethodNotAllowed)
			http.NotFound(rw, req)
			return
		}

		// get deployment data from request body
		reqBody, err := io.ReadAll(req.Body)
		if err != nil {
			message := fmt.Sprintf("Cannot execute deploy service, IO error while reading data file: %s", err.Error())
			logger.Errorf(message)
			http.Error(rw, message, http.StatusUnprocessableEntity)
			return
		}

		// convert deployment data to deployer config
		deployerConfig := &config.DeployerConfig{}
		err = yaml.Unmarshal(reqBody, deployerConfig)
		if err != nil {
			message := fmt.Sprintf("Cannot unmarshal request body to expected model: %s", err.Error())
			logger.Errorf(message)
			logger.Errorf("Request body: %# v", string(reqBody))
			http.Error(rw, message, http.StatusUnprocessableEntity)
			return
		}

		// validate deployemnt data
		err = deployerConfig.Deployer.Validate()
		if err != nil {
			handleError("Invalid configuration: %s", err, rw, deployerConfig.Deployer.ApplicationConfig.Name, logger)
			return
		}
		logger.Infof("Running deploy request on: %# v", string(reqBody))

		// connect to cluster
		kubeConnection := connect.ConnectToCluster(deployerConfig.Deployer.ClusterConfig, false, logger)
		if kubeConnection.Error != nil {
			handleError("Cannot connect to given cluster: %s", kubeConnection.Error, rw, deployerConfig.Deployer.ApplicationConfig.Name, logger)
			return
		}

		// create deployer context to hold global variables
		deployerContext := context.InitDeployerContext(args, kubeConnection)

		// create application deployment package
		applicationPkg := packager.NewApplicationPkg(deployerConfig.Deployer.ApplicationConfig, deployerContext)
		// ADD HANDLE ERROR

		// check if all mandatory variables have been set, else list unset vars and throw exception
		if len(applicationPkg.UnsetMandatoryParams) > 0 {
			UnsetMandatoryParamsMsg := fmt.Sprintf("Missing configuration for the following mandatory parameters (<FILEPATH>: <MANDATORY_PARAMETERS>.)\n%s", utils.StringifyMap(applicationPkg.UnsetMandatoryParams))
			handleError(UnsetMandatoryParamsMsg, nil, rw, deployerConfig.Deployer.ApplicationConfig.Name, logger)
			return
		}

		// deploy application deployment package
		deployer.DeployApplication(applicationPkg)
		// ADD HANDLE ERROR

	}
}

func info(rw http.ResponseWriter, req *http.Request) {

	// validate requested path
	if req.URL.Path != APP_DEPLOYER_PATH {
		http.Error(rw, fmt.Sprintf("Unmanaged path %s", req.URL.Path), http.StatusNotFound)
		http.NotFound(rw, req)
		return
	}

	// validate requested method
	if req.Method != GET {
		http.Error(rw, fmt.Sprintf("Expect method %s at %s, got %v", GET, APP_DEPLOYER_PATH, req.Method), http.StatusMethodNotAllowed)
		http.NotFound(rw, req)
		return
	}

	// init application info
	rw.WriteHeader(http.StatusOK)
	rw.Header().Set("Content-Type", CONTENT_TYPE)
	applicationInfo := applicationInfo{
		Name:    APPLICATION_NAME,
		Version: BUILD_VERSION,
		Status:  STATUS,
	}

	// transmit application info
	applicationInfoYaml, _ := json.Marshal(applicationInfo)
	rw.Write([]byte(applicationInfoYaml))
}
