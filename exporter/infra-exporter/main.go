package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/RHEcosystemAppEng/SaaSi/exporter/exporter-lib/config"
	"github.com/RHEcosystemAppEng/SaaSi/exporter/exporter-lib/connect"
	"github.com/RHEcosystemAppEng/SaaSi/exporter/exporter-lib/export/infra"
	"github.com/RHEcosystemAppEng/SaaSi/exporter/exporter-lib/export/utils"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type InfraExporterService struct {
	config *config.Config
	logger *logrus.Logger
}

type applicationInfo struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
	Status  string `yaml:"status"`
}

var BuildVersion = "development"
var router = mux.NewRouter()

func main() {
	infraExporterService := InfraExporterService{}
	infraExporterService.config = config.ReadConfigFromEnvVars()
	infraExporterService.logger = utils.GetLogger(infraExporterService.config.Debug)
	utils.PrettyPrint(infraExporterService.logger, "Runtime configuration: %s", infraExporterService.config)

	infraExporterService.logger.Infof("Running %s with version %s", os.Args[0], BuildVersion)

	router.Path("/export/infra").HandlerFunc(infraExporterService.export).Methods("POST")
	router.Path("/export/infra").HandlerFunc(infraExporterService.info).Methods("GET")

	host := "0.0.0.0"
	url := fmt.Sprintf("%s:%d", host, 8080)
	infraExporterService.logger.Infof("Starting listener as %s", url)
	if err := http.ListenAndServe(url, router); err != nil {
		infraExporterService.logger.Fatal(err)
	}
}

func (e *InfraExporterService) export(rw http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/export/infra" {
		if req.Method == "POST" {
			reqBody, err := io.ReadAll(req.Body)
			if err != nil {
				message := fmt.Sprintf("Cannot execute export service, IO error whiole reading the data file: %s", err.Error())
				e.logger.Errorf(message)
				http.Error(rw, message, http.StatusUnprocessableEntity)
				return
			}

			clusterConfig := &config.ClusterConfig{}
			err = yaml.Unmarshal(reqBody, clusterConfig)
			if err != nil {
				message := fmt.Sprintf("Cannot unmarshal request body to expected model: %s", err.Error())
				e.logger.Errorf(message)
				e.logger.Errorf("Request body is: %# v", string(reqBody))
				http.Error(rw, message, http.StatusUnprocessableEntity)
				return
			}

			err = clusterConfig.ValidateClusterConfig()
			if err != nil {
				e.handleError("Invalid configuration: %s", err, rw, clusterConfig)
				return
			}
			e.logger.Infof("Running export request: %# v", string(reqBody))
			connectionStatus := connect.ConnectCluster(clusterConfig, e.logger)
			if connectionStatus.Error != nil {
				e.handleError("Cannot connect to given cluster: %s", connectionStatus.Error, rw, clusterConfig)
				return
			}

			infraExporter := infra.NewInfraExporterFromConfig(e.config, clusterConfig, connectionStatus, e.logger)
			output := infraExporter.Export()
			yamlOutput, err := json.Marshal(output)
			if err != nil {
				message := fmt.Sprintf("Cannot marshal response output to expected model: %s", err.Error())
				e.logger.Errorf("Output is: %# v", output)
				http.Error(rw, message, http.StatusUnprocessableEntity)
				return
			}
			rw.WriteHeader(http.StatusOK)
			rw.Header().Set("Content-Type", "application/json")
			rw.Write([]byte(yamlOutput))
		} else {
			http.Error(rw, fmt.Sprintf("Expect method POST at /export/infra, got %v", req.Method), http.StatusMethodNotAllowed)
		}
	} else {
		http.Error(rw, fmt.Sprintf("Unmanaged path %s", req.URL.Path), http.StatusNotFound)
		http.NotFound(rw, req)
		return
	}
}

func (e *InfraExporterService) handleError(message string, err error, rw http.ResponseWriter, clusterConfig *config.ClusterConfig) {
	message = fmt.Sprintf(message, err.Error())
	e.logger.Errorf(message)
	rw.WriteHeader(http.StatusBadRequest)
	rw.Header().Set("Content-Type", "application/json")
	output := infra.InfraExporterOutput{ClusterId: clusterConfig.ClusterId, Status: utils.Failed.String(),
		ErrorMessage: message}
	yamlOutput, _ := json.Marshal(output)
	rw.Write([]byte(yamlOutput))

}
func (e *InfraExporterService) info(rw http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/export/infra" {
		if req.Method == "GET" {
			rw.WriteHeader(http.StatusOK)
			rw.Header().Set("Content-Type", "application/json")
			applicationInfo := applicationInfo{Name: "infra-exporter", Version: BuildVersion, Status: "up"}
			yamlOutput, _ := json.Marshal(applicationInfo)

			rw.Write([]byte(yamlOutput))
		} else {
			http.Error(rw, fmt.Sprintf("Expect method GET at /export/infra, got %v", req.Method), http.StatusMethodNotAllowed)
		}
	} else {
		http.Error(rw, fmt.Sprintf("Unmanaged path %s", req.URL.Path), http.StatusNotFound)
		http.NotFound(rw, req)
		return
	}
}
