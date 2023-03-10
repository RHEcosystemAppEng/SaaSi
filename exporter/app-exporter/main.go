package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/RHEcosystemAppEng/SaaSi/exporter/exporter-lib/config"
	"github.com/RHEcosystemAppEng/SaaSi/exporter/exporter-lib/export"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type AppExporter struct {
	config   *config.Config
	logger   *logrus.Logger
	exporter *export.Exporter
}

var router = mux.NewRouter()

func main() {
	appExporter := AppExporter{}
	appExporter.config = config.ReadConfig()
	appExporter.logger = appExporter.config.Logger
	appExporter.exporter = export.NewExporterFromConfig(appExporter.config)

	router.Path("/export/application").HandlerFunc(appExporter.export).Methods("POST")

	host := "0.0.0.0"
	url := fmt.Sprintf("%s:%d", host, 8080)
	appExporter.config.Logger.Infof("Starting listener as %s", url)
	if err := http.ListenAndServe(url, router); err != nil {
		appExporter.config.Logger.Fatal(err)
	}
}

func (e *AppExporter) export(rw http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/export/application" {
		if req.Method == "POST" {
			reqBody, err := io.ReadAll(req.Body)
			if err != nil {
				message := fmt.Sprintf("Cannot execute export service, IO error whiole reading the data file: %s", err.Error())
				e.logger.Errorf(message)
				http.Error(rw, message, http.StatusUnprocessableEntity)
				return
			}

			exporterConfig := config.ExporterConfig{}
			err = yaml.Unmarshal(reqBody, &exporterConfig)
			if err != nil {
				message := fmt.Sprintf("Cannot unmarshal request body to expected model: %s", err.Error())
				e.logger.Errorf(message)
				e.logger.Errorf("Request body is: %# v", string(reqBody))
				http.Error(rw, message, http.StatusUnprocessableEntity)
				return
			}

			e.logger.Infof("Running export request: %# v", string(reqBody))
			output := e.exporter.Export(&exporterConfig)

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
			http.Error(rw, fmt.Sprintf("Expect method POST at /, got %v", req.Method), http.StatusMethodNotAllowed)
		}
	} else {
		http.Error(rw, fmt.Sprintf("Unmanaged path %s", req.URL.Path), http.StatusNotFound)
		http.NotFound(rw, req)
		return
	}
}