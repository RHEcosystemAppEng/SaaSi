package app

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/RHEcosystemAppEng/SaaSi/exporter/exporter-lib/config"
	"github.com/RHEcosystemAppEng/SaaSi/exporter/exporter-lib/connect"
	"github.com/RHEcosystemAppEng/SaaSi/exporter/exporter-lib/export/utils"
	"github.com/google/uuid"
	"github.com/konveyor/crane/cmd/apply"
	"github.com/konveyor/crane/cmd/export"
	"github.com/konveyor/crane/cmd/transform"
	"github.com/sirupsen/logrus"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

type AppExporter struct {
	appContext *AppContext
	logger     *logrus.Logger
}

type AppExporterOutput struct {
	Uid             string `json:"uid"`
	ApplicationName string `json:"applicationName"`
	Status          string `json:"status"`
	ErrorMessage    string `json:"errorMessage"`
	Location        string `json:"location"`
}

func NewAppExporterFromConfig(config *config.Config, exporterConfig *config.ExporterConfig, connectionStatus *connect.ConnectionStatus, logger *logrus.Logger) *AppExporter {
	exporter := AppExporter{appContext: NewAppContextFromConfig(config, exporterConfig, connectionStatus, logger)}
	exporter.logger = logger

	return &exporter
}

func (e *AppExporter) Export() AppExporterOutput {
	output := AppExporterOutput{ApplicationName: e.appContext.AppConfig.Name}
	output.Uid = uuid.New().String()
	e.appContext.InitAppFolderForUid(output.Uid)

	err := e.PrepareOutput()
	if err != nil {
		output.ErrorMessage = err.Error()
		output.Status = utils.Failed.String()
	}
	err = e.ExportWithCrane()
	if err != nil {
		output.ErrorMessage = err.Error()
		output.Status = utils.Failed.String()
	}

	parametrizer := NewParametrizerFromConfig(e.appContext)
	err = parametrizer.ExposeParameters()
	if err != nil {
		output.ErrorMessage = err.Error()
		output.Status = utils.Failed.String()
	}

	installer := NewInstallerFromConfig(e.appContext)
	err = installer.BuildKustomizeInstaller()
	if err != nil {
		output.ErrorMessage = err.Error()
		output.Status = utils.Failed.String()
	} else {
		output.ErrorMessage = ""
		output.Status = utils.Ok.String()
		output.Location = e.appContext.AppFolder
	}

	return output
}

func (e *AppExporter) PrepareOutput() error {
	if _, err := os.Stat(e.appContext.AppFolder); !os.IsNotExist(err) {
		e.logger.Infof("Directory exists %v", e.appContext.AppFolder)
		os.RemoveAll(e.appContext.AppFolder)
	}

	if err := os.MkdirAll(e.appContext.AppFolder, os.ModePerm); err != nil {
		return errors.New(fmt.Sprintf("Cannot create %v folder: %v", e.appContext.AppFolder, err))
	}
	e.logger.Infof("Created output folder %s", e.appContext.AppFolder)
	for _, ns := range e.appContext.AppConfig.Namespaces {
		nsFolder := filepath.Join(e.appContext.AppFolder, ns.Name)
		if err := os.Mkdir(nsFolder, os.ModePerm); err != nil {
			return errors.New(fmt.Sprintf("Cannot create %v folder: %v", nsFolder, err))
		}
		e.logger.Infof("Created output folder %s", nsFolder)
	}
	return nil
}

func (e *AppExporter) ExportWithCrane() error {
	for _, ns := range e.appContext.AppConfig.Namespaces {
		nsFolder := filepath.Join(e.appContext.AppFolder, ns.Name)
		e.logger.Infof("Exporting NS %s with crane", nsFolder)

		exportFolder := e.appContext.ExportFolderForNS(ns.Name)
		transformFolder := e.appContext.TransformFolderForNS(ns.Name)
		outputFolder := e.appContext.OutputFolderForNS(ns.Name)

		err := doExport(e.appContext.KubeConfigPath(), ns.Name, exportFolder, e.appContext.Debug)
		if err != nil {
			return err
		}
		err = doTransform(exportFolder, transformFolder, e.appContext.Debug)
		if err != nil {
			return err
		}
		err = doApply(exportFolder, transformFolder, outputFolder, e.appContext.Debug)
		if err != nil {
			return err
		}
	}
	return nil
}

func doExport(kubeConfigPath string, namespace string, exportFolder string, debug bool) error {
	exportCmd := export.NewExportCommand(genericclioptions.IOStreams{
		In:     strings.NewReader(""),
		Out:    os.Stdout,
		ErrOut: os.Stderr,
	}, nil)
	if debug {
		exportCmd.PersistentFlags().Bool("debug", false, "Debug the command by printing more information")
	}

	clusterScopedRbac := exportCmd.Flags().Lookup("cluster-scoped-rbac")
	clusterScopedRbac.Value.Set("true")
	exportNamespace := exportCmd.Flags().Lookup("namespace")
	exportNamespace.Value.Set(namespace)
	exportDir := exportCmd.Flags().Lookup("export-dir")
	exportDir.Value.Set(exportFolder)
	kubeconfig := exportCmd.Flags().Lookup("kubeconfig")
	kubeconfig.Value.Set(kubeConfigPath)
	if debug {
		debugFlag := exportCmd.PersistentFlags().Lookup("debug")
		debugFlag.Value.Set("true")
	}
	exportCmd.SetArgs([]string{})

	_, err := exportCmd.ExecuteC()
	return err
}

func doTransform(exportFolder string, transformFolder string, debug bool) error {
	transformCmd := transform.NewTransformCommand(nil)
	if debug {
		transformCmd.PersistentFlags().Bool("debug", false, "Debug the command by printing more information")
	}

	exportDir := transformCmd.Flags().Lookup("export-dir")
	exportDir.Value.Set(exportFolder)
	transformDir := transformCmd.Flags().Lookup("transform-dir")
	transformDir.Value.Set(transformFolder)
	if debug {
		debugFlag := transformCmd.PersistentFlags().Lookup("debug")
		debugFlag.Value.Set("true")
	}
	transformCmd.SetArgs([]string{})

	_, err := transformCmd.ExecuteC()
	return err
}

func doApply(exportFolder string, transformFolder string, outputFolder string, debug bool) error {
	applyCmd := apply.NewApplyCommand(nil)
	if debug {
		applyCmd.PersistentFlags().Bool("debug", false, "Debug the command by printing more information")
	}

	exportDir := applyCmd.Flags().Lookup("export-dir")
	exportDir.Value.Set(exportFolder)
	transformDir := applyCmd.Flags().Lookup("transform-dir")
	transformDir.Value.Set(transformFolder)
	outputDir := applyCmd.Flags().Lookup("output-dir")
	outputDir.Value.Set(outputFolder)
	if debug {
		debugFlag := applyCmd.PersistentFlags().Lookup("debug")
		debugFlag.Value.Set("true")
	}
	applyCmd.SetArgs([]string{})

	_, err := applyCmd.ExecuteC()
	return err
}
