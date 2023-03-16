package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
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

func HandleRequests() {

	router.Path(APP_DEPLOYER_PATH).HandlerFunc(deploy).Methods(POST)
	router.Path(APP_DEPLOYER_PATH).HandlerFunc(info).Methods(GET)

	// init hosting URL
	url := fmt.Sprintf("%s:%d", HOST, PORT)
	// appDeployerService.logger.Infof("Starting listener as %s", url)
	if err = http.ListenAndServe(url, router); err != nil {
		// appDeployerService.logger.Fatal(err)
		fmt.Print() //PLACEHOLDER
	}
}
