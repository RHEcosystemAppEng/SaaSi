package installer

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/RHEcosystemAppEng/SaaSi/replica-builder/install-builder/pkg/config"
)

type Exporter struct {
	appConfig       *config.ApplicationConfig
	installerConfig *config.InstallerConfig
}

func NewExporterFromConfig(appConfig *config.ApplicationConfig, installerConfig *config.InstallerConfig) *Exporter {
	validateRequirements()
	exporter := Exporter{appConfig: appConfig, installerConfig: installerConfig}

	return &exporter
}

func validateRequirements() {
	_, err := exec.LookPath("crane")
	if err != nil {
		log.Fatal("crane command not found")
	}

	// _, err = exec.LookPath("move2kube")
	// if err != nil {
	// 	log.Fatal("move2kube command not found")
	// }

	_, err = exec.LookPath("oc")
	if err != nil {
		log.Fatal("oc command not found")
	}

	_, err = exec.LookPath("kustomize")
	if err != nil {
		log.Fatal("kustomize command not found")
	}
	log.Printf("Requirements validated")
}

func (e *Exporter) PrepareOutput() {
	if _, err := os.Stat(e.installerConfig.AppFolder); !os.IsNotExist(err) {
		log.Printf("Directory exists %v", e.installerConfig.AppFolder)
		os.RemoveAll(e.installerConfig.AppFolder)
	}

	if err := os.MkdirAll(e.installerConfig.AppFolder, os.ModePerm); err != nil {
		log.Fatalf("Cannot create %v folder: %v", e.installerConfig.AppFolder, err)
	}
	log.Printf("Created output folder %s", e.installerConfig.AppFolder)
	for _, ns := range e.appConfig.Application.Namespaces {
		nsFolder := filepath.Join(e.installerConfig.AppFolder, ns.Name)
		if err := os.Mkdir(nsFolder, os.ModePerm); err != nil {
			log.Fatalf("Cannot create %v folder: %v", nsFolder, err)
		}
		log.Printf("Created output folder %s", nsFolder)
	}
}

func (e *Exporter) ExportWithCrane() {
	for _, ns := range e.appConfig.Application.Namespaces {
		nsFolder := filepath.Join(e.installerConfig.AppFolder, ns.Name)
		log.Printf("Exporting NS %s with crane", nsFolder)

		exportFolder := e.installerConfig.ExportFolderForNS(ns.Name)
		transformFolder := e.installerConfig.TransformFolderForNS(ns.Name)
		outputFolder := e.installerConfig.OutputFolderForNS(ns.Name)

		RunCommand("crane", "export", "--namespace", ns.Name, "--export-dir", exportFolder)
		RunCommand("crane", "transform", "--export-dir", exportFolder,
			"--transform-dir", transformFolder)
		RunCommand("crane", "apply", "--export-dir", exportFolder,
			"--transform-dir", transformFolder,
			"--output-dir", outputFolder)
	}
}
