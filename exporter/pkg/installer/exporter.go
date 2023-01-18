package installer

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/RHEcosystemAppEng/SaaSi/exporter/pkg/config"
	"github.com/konveyor/crane/cmd/apply"
	"github.com/konveyor/crane/cmd/export"
	"github.com/konveyor/crane/cmd/transform"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

type Exporter struct {
	appConfig       *config.ApplicationConfig
	installerConfig *config.InstallerConfig
}

func NewExporterFromConfig(appConfig *config.ApplicationConfig, installerConfig *config.InstallerConfig) *Exporter {
	exporter := Exporter{appConfig: appConfig, installerConfig: installerConfig}

	return &exporter
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

		doExport(ns.Name, exportFolder)
		doTransform(exportFolder, transformFolder)
		doApply(exportFolder, transformFolder, outputFolder)
	}
}

func doExport(namespace string, exportFolder string) {
	exportCmd := export.NewExportCommand(genericclioptions.IOStreams{
		In:     strings.NewReader(""),
		Out:    os.Stdout,
		ErrOut: os.Stderr,
	}, nil)

	exportNamespace := exportCmd.Flags().Lookup("namespace")
	exportNamespace.Value.Set(namespace)
	exportDir := exportCmd.Flags().Lookup("export-dir")
	exportDir.Value.Set(exportFolder)
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
