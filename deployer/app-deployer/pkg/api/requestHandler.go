package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/RHEcosystemAppEng/SaaSi/deployer/deployer-lib/config"
	"github.com/RHEcosystemAppEng/SaaSi/deployer/deployer-lib/deployer/app"
	"github.com/RHEcosystemAppEng/SaaSi/deployer/deployer-lib/utils"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

const (
	APP_DEPLOYER_PATH = "/deploy/application"
	HOST              = "0.0.0.0"
	POST              = "POST"
	GET               = "GET"
	CONTENT_TYPE      = "application/json"
	APPLICATION_NAME  = "app-deployer"
	STATUS            = "up"
)

var (
	err          error
	BuildVersion = "development"
	router       = mux.NewRouter()
)

type applicationInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Status  string `json:"status"`
}

func HandleRequests(args *config.Args, logger *logrus.Logger) {

	logger.Infof("Running %s with version %s", os.Args[0], BuildVersion)

	// define routes
	router.Path(APP_DEPLOYER_PATH).HandlerFunc(deploy(args, logger)).Methods(POST)
	router.Path(APP_DEPLOYER_PATH).HandlerFunc(info).Methods(GET)

	// init hosting URL
	url := fmt.Sprintf("%s:%d", HOST, args.Port)
	logger.Infof("Starting listener as %s", url)
	if err = http.ListenAndServe(url, router); err != nil {
		logger.Fatal(err)
	}
}

func handleResponse(rw http.ResponseWriter, logger *logrus.Logger, output *app.ApplicationOutput) {

	if output.Status == utils.Failed.String() {

		// set http parameters for error
		logger.Error(output.ErrorMessage)
		rw.WriteHeader(http.StatusBadRequest)
	} else {

		// set http parameters for ok
		logger.Info("Process completed successfully")
		rw.WriteHeader(http.StatusOK)
	}

	// marshal output to json format
	jsonOutput, _ := json.Marshal(output)

	// set http parameters and produce response
	rw.Header().Set("Content-Type", CONTENT_TYPE)
	rw.Write([]byte(jsonOutput))
}

func handleInfo(rw http.ResponseWriter) {

	// set output parameters and marshal to json format
	output := applicationInfo{
		Name:    APPLICATION_NAME,
		Version: BuildVersion,
		Status:  STATUS,
	}
	jsonOutput, _ := json.Marshal(output)

	// set http parameters and produce response
	rw.WriteHeader(http.StatusOK)
	rw.Header().Set("Content-Type", CONTENT_TYPE)
	rw.Write([]byte(jsonOutput))
}
