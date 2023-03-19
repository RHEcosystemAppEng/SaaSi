package api

import (
	"fmt"
	"net/http"

	"github.com/RHEcosystemAppEng/SaaSi/deployer/deployer-lib/config"
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
