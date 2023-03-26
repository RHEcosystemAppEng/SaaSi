package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/RHEcosystemAppEng/SaaSi/deployer/deployer-lib/config"
	"github.com/RHEcosystemAppEng/SaaSi/deployer/deployer-lib/utils"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

const (
	APP_DEPLOYER_PATH = "/deploy/application"
	HOST              = "0.0.0.0"
	PORT              = 8080
	POST              = "POST"
	GET               = "GET"
	CONTENT_TYPE      = "application/json"
	APPLICATION_NAME  = "app-deployer"
	BUILD_VERSION     = "dev"
	STATUS            = "up"
)

var (
	err    error
	router = mux.NewRouter()
)

type ApplicationOutput struct {
	ApplicationName string `json:"applicationName"`
	Status          string `json:"status"`
	ErrorMessage    string `json:"errorMessage"`
	Location        string `json:"location"`
}

type applicationInfo struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
	Status  string `yaml:"status"`
}

func HandleRequests(args *config.Args, logger *logrus.Logger) {

	// define routes
	router.Path(APP_DEPLOYER_PATH).HandlerFunc(deploy(args, logger)).Methods(POST)
	router.Path(APP_DEPLOYER_PATH).HandlerFunc(info).Methods(GET)

	// init hosting URL
	url := fmt.Sprintf("%s:%d", HOST, PORT)
	logger.Infof("Starting listener as %s", url)
	if err = http.ListenAndServe(url, router); err != nil {
		logger.Fatal(err)
	}
}

func handleError(rw http.ResponseWriter, logger *logrus.Logger, message string, applicationName string) {

	logger.Errorf(message)

	// set output parameters and marshal to json format
	output := ApplicationOutput{
		ApplicationName: applicationName,
		Status:          utils.Failed.String(),
		ErrorMessage:    message,
		Location:        "",
	}
	jsonOutput, _ := json.Marshal(output)

	// set http parameters and produce response
	rw.WriteHeader(http.StatusBadRequest)
	rw.Header().Set("Content-Type", CONTENT_TYPE)
	rw.Write([]byte(jsonOutput))
}

func handleOk(rw http.ResponseWriter, logger *logrus.Logger, applicationName string, outputDir string) {

	logger.Infof("Application %s deployer successfully")

	// set output parameters and marshal to json format
	output := ApplicationOutput{
		ApplicationName: applicationName,
		Status:          STATUS,
		ErrorMessage:    "",
		Location:        outputDir,
	}
	jsonOutput, _ := json.Marshal(output)

	// set http parameters and produce response
	rw.WriteHeader(http.StatusOK)
	rw.Header().Set("Content-Type", CONTENT_TYPE)
	rw.Write([]byte(jsonOutput))
}

func handleInfo(rw http.ResponseWriter) {

	// set output parameters and marshal to json format
	output := applicationInfo{
		Name:    APPLICATION_NAME,
		Version: BUILD_VERSION,
		Status:  STATUS,
	}
	jsonOutput, _ := json.Marshal(output)

	// set http parameters and produce response
	rw.WriteHeader(http.StatusOK)
	rw.Header().Set("Content-Type", CONTENT_TYPE)
	rw.Write([]byte(jsonOutput))
}
