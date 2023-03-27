package api

import (
	"fmt"
	"io"
	"net/http"

	"github.com/RHEcosystemAppEng/SaaSi/deployer/deployer-lib/config"
	"github.com/RHEcosystemAppEng/SaaSi/deployer/deployer-lib/deployer/app"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

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
		logger.Infof("Running export request: %# v", string(reqBody))

		output := app.Deploy(deployerConfig.Deployer, args, logger)

		handleResponse(rw, logger, output)
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

	handleInfo(rw)
}
