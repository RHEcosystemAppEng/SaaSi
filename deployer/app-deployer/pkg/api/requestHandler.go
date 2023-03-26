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

func handleError(message string, err error, rw http.ResponseWriter, applicationName string, logger *logrus.Logger) {
	message = fmt.Sprintf(message, err.Error())
	logger.Errorf(message)
	rw.WriteHeader(http.StatusBadRequest)
	rw.Header().Set("Content-Type", "application/json")
	output := ApplicationOutput{ApplicationName: applicationName, Status: utils.Failed.String(), ErrorMessage: message}
	jsonOutput, _ := json.Marshal(output)
	rw.Write([]byte(jsonOutput))
}
