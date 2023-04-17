package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/RHEcosystemAppEng/SaaSi/deployer/deployer-lib/config"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

const (
	INFRA_DEPLOYER_PATH = "/deploy/infra"
	HOST                = "0.0.0.0"
	POST                = "POST"
	GET                 = "GET"
	CONTENT_TYPE        = "application/json"
	APPLICATION_NAME    = "infra-deployer"
	STATUS              = "up"
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
	router.Path(INFRA_DEPLOYER_PATH).HandlerFunc(deploy(args, logger)).Methods(POST)
	router.Path(INFRA_DEPLOYER_PATH).HandlerFunc(info).Methods(GET)

	// init hosting URL
	url := fmt.Sprintf("%s:%d", HOST, args.Port)
	logger.Infof("Starting listener as %s", url)
	if err = http.ListenAndServe(url, router); err != nil {
		logger.Fatal(err)
	}
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
