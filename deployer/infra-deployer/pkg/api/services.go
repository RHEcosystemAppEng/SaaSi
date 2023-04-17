package api

import (
	"fmt"
	"net/http"

	"github.com/RHEcosystemAppEng/SaaSi/deployer/deployer-lib/config"
	"github.com/sirupsen/logrus"
)

func deploy(args *config.Args, logger *logrus.Logger) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
	}
}

func info(rw http.ResponseWriter, req *http.Request) {

	// validate requested path
	if req.URL.Path != INFRA_DEPLOYER_PATH {
		http.Error(rw, fmt.Sprintf("Unmanaged path %s", req.URL.Path), http.StatusNotFound)
		http.NotFound(rw, req)
		return
	}

	// validate requested method
	if req.Method != GET {
		http.Error(rw, fmt.Sprintf("Expect method %s at %s, got %v", GET, INFRA_DEPLOYER_PATH, req.Method), http.StatusMethodNotAllowed)
		http.NotFound(rw, req)
		return
	}

	handleInfo(rw)
}
