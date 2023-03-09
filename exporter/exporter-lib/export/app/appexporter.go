package app

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/RHEcosystemAppEng/SaaSi/exporter/exporter-lib/config"
	"github.com/RHEcosystemAppEng/SaaSi/exporter/exporter-lib/connect"
	"github.com/konveyor/crane/cmd/apply"
	"github.com/konveyor/crane/cmd/export"
	"github.com/konveyor/crane/cmd/transform"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

type AppExporter struct {
	appContext *AppContext
}

func NewAppExporterFromConfig(config *config.Config, connectionStatus *connect.ConnectionStatus) *AppExporter {
	exporter := AppExporter{appContext: NewAppContextFromConfig(config, connectionStatus)}

	return &exporter
}

func (e *AppExporter) Export() {
	e.PrepareOutput()
	e.ExportWithCrane()

	parametrizer := NewParametrizerFromConfig(e.appContext)
	parametrizer.ExposeParameters()

	installer := NewInstallerFromConfig(e.appContext)
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
	for _, ns := range e.appContext.AppConfig.Namespaces {
		nsFolder := filepath.Join(e.appContext.AppFolder, ns.Name)
		if err := os.Mkdir(nsFolder, os.ModePerm); err != nil {
			log.Fatalf("Cannot create %v folder: %v", nsFolder, err)
		}
		log.Printf("Created output folder %s", nsFolder)
	}
}

func (e *AppExporter) ExportWithCrane() {
	for _, ns := range e.appContext.AppConfig.Namespaces {
		nsFolder := filepath.Join(e.appContext.AppFolder, ns.Name)
		log.Printf("Exporting NS %s with crane", nsFolder)

		exportFolder := e.appContext.ExportFolderForNS(ns.Name)
		transformFolder := e.appContext.TransformFolderForNS(ns.Name)
		outputFolder := e.appContext.OutputFolderForNS(ns.Name)

		doExport(e.appContext.KubeConfigPath(), ns.Name, exportFolder)
		doTransform(exportFolder, transformFolder)
		doApply(exportFolder, transformFolder, outputFolder)
	}
}

func doExport(kubeConfigPath string, namespace string, exportFolder string) {
	exportCmd := export.NewExportCommand(genericclioptions.IOStreams{
		In:     strings.NewReader(""),
		Out:    os.Stdout,
		ErrOut: os.Stderr,
	}, nil)

	clusterScopedRbac := exportCmd.Flags().Lookup("cluster-scoped-rbac")
	clusterScopedRbac.Value.Set("true")
	exportNamespace := exportCmd.Flags().Lookup("namespace")
	exportNamespace.Value.Set(namespace)
	exportDir := exportCmd.Flags().Lookup("export-dir")
	exportDir.Value.Set(exportFolder)
	kubeconfig := exportCmd.Flags().Lookup("kubeconfig")
	kubeconfig.Value.Set(kubeConfigPath)
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
