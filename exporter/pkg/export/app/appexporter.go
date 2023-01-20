package app

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/RHEcosystemAppEng/SaaSi/exporter/pkg/config"
	"github.com/RHEcosystemAppEng/SaaSi/exporter/pkg/connect"
	"github.com/konveyor/crane/cmd/apply"
	"github.com/konveyor/crane/cmd/export"
	"github.com/konveyor/crane/cmd/transform"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

type AppExporter struct {
	appConfig  *config.ApplicationConfig
	appContext *AppContext
}

func NewAppExporterFromConfig(config *config.Config) *AppExporter {
	exporter := AppExporter{appConfig: &config.Exporter.Application}

	return &exporter
}

func (e *AppExporter) Export(connectionStatus *connect.ConnectionStatus) {
	e.appContext = NewContextFromConfig(e.appConfig, connectionStatus)

	clusterRolesInspector := NewClusterRolesInspector(connectionStatus)
	clusterRolesInspector.LoadClusterRoles()

	e.PrepareOutput()
	e.ExportWithCrane()

	parametrizer := NewParametrizerFromConfig(e.appConfig, e.appContext)
	parametrizer.ExposeParameters()

	installer := NewInstallerFromConfig(e.appConfig, e.appContext, clusterRolesInspector)
	installer.BuildKustomizeInstaller()
}

func (e *AppExporter) PrepareOutput() {
	if _, err := os.Stat(e.appContext.AppFolder); !os.IsNotExist(err) {
		log.Printf("Directory exists %v", e.appContext.AppFolder)
		os.RemoveAll(e.appContext.AppFolder)
	}

	if err := os.MkdirAll(e.appContext.AppFolder, os.ModePerm); err != nil {
		log.Fatalf("Cannot create %v folder: %v", e.appContext.AppFolder, err)
	}
	log.Printf("Created output folder %s", e.appContext.AppFolder)
	for _, ns := range e.appConfig.Namespaces {
		nsFolder := filepath.Join(e.appContext.AppFolder, ns.Name)
		if err := os.Mkdir(nsFolder, os.ModePerm); err != nil {
			log.Fatalf("Cannot create %v folder: %v", nsFolder, err)
		}
		log.Printf("Created output folder %s", nsFolder)
	}
}

func (e *AppExporter) ExportWithCrane() {
	for _, ns := range e.appConfig.Namespaces {
		nsFolder := filepath.Join(e.appContext.AppFolder, ns.Name)
		log.Printf("Exporting NS %s with crane", nsFolder)

		exportFolder := e.appContext.ExportFolderForNS(ns.Name)
		transformFolder := e.appContext.TransformFolderForNS(ns.Name)
		outputFolder := e.appContext.OutputFolderForNS(ns.Name)

		doExport(e.appContext.ConnectionStatus, ns.Name, exportFolder)
		doTransform(exportFolder, transformFolder)
		doApply(exportFolder, transformFolder, outputFolder)
	}
}

func doExport(connectionStatus *connect.ConnectionStatus, namespace string, exportFolder string) {
	exportCmd := export.NewExportCommand(genericclioptions.IOStreams{
		In:     strings.NewReader(""),
		Out:    os.Stdout,
		ErrOut: os.Stderr,
	}, nil)

	exportNamespace := exportCmd.Flags().Lookup("namespace")
	exportNamespace.Value.Set(namespace)
	exportDir := exportCmd.Flags().Lookup("export-dir")
	exportDir.Value.Set(exportFolder)
	kubeconfig := exportCmd.Flags().Lookup("kubeconfig")
	kubeconfig.Value.Set(connectionStatus.KubeConfigPath)
	exportCmd.SetArgs([]string{})

	_, err := exportCmd.ExecuteC()
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func doTransform(exportFolder string, transformFolder string) {
	transformCmd := transform.NewTransformCommand(nil)

	exportDir := transformCmd.Flags().Lookup("export-dir")
	exportDir.Value.Set(exportFolder)
	transformDir := transformCmd.Flags().Lookup("transform-dir")
	transformDir.Value.Set(transformFolder)
	transformCmd.SetArgs([]string{})

	_, err := transformCmd.ExecuteC()
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func doApply(exportFolder string, transformFolder string, outputFolder string) {
	applyCmd := apply.NewApplyCommand(nil)

	exportDir := applyCmd.Flags().Lookup("export-dir")
	exportDir.Value.Set(exportFolder)
	transformDir := applyCmd.Flags().Lookup("transform-dir")
	transformDir.Value.Set(transformFolder)
	outputDir := applyCmd.Flags().Lookup("output-dir")
	outputDir.Value.Set(outputFolder)
	applyCmd.SetArgs([]string{})

	_, err := applyCmd.ExecuteC()
	if err != nil {
		log.Fatalf(err.Error())
	}

}
