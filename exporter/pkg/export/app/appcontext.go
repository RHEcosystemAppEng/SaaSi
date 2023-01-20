package app

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/RHEcosystemAppEng/SaaSi/exporter/pkg/config"
	"github.com/RHEcosystemAppEng/SaaSi/exporter/pkg/connect"
)

var ParamsFolder = "params"
var SecretsFolder = "secrets"
var NoValue = "__DEFAULT__"
var MandatoryValue = "__MANDATORY__"
var KustomizationFile = "kustomization.yaml"

type AppContext struct {
	OutputFolder     string
	AppFolder        string
	ConnectionStatus *connect.ConnectionStatus
}

func NewContextFromConfig(appConfig *config.ApplicationConfig, connectionStatus *connect.ConnectionStatus) *AppContext {
	context := AppContext{ConnectionStatus: connectionStatus}

	pwd, _ := os.Getwd()
	log.Printf("Running from %v", pwd)
	context.OutputFolder = pwd + "/output"
	// TODO
	context.AppFolder = filepath.Join(context.OutputFolder, appConfig.Name)

	return &context
}

func (i *AppContext) ExportFolderForNS(namespace string) string {
	return filepath.Join(i.AppFolder, namespace, "export")
}

func (i *AppContext) TransformFolderForNS(namespace string) string {
	return filepath.Join(i.AppFolder, namespace, "transform")
}

func (i *AppContext) OutputFolderForNS(namespace string) string {
	return filepath.Join(i.AppFolder, namespace, "output")
}

func (i *AppContext) lookupOrCreateFolder(path ...string) string {
	fullPath := ""
	if strings.HasPrefix(path[0], i.AppFolder) {
		fullPath = filepath.Join(path...)
	} else {
		fullPath = filepath.Join(append([]string{i.AppFolder}, path...)...)
	}
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		if err := os.MkdirAll(fullPath, os.ModePerm); err != nil {
			log.Fatalf("Cannot create %v folder: %v", fullPath, err)
		}
	}
	return fullPath
}

func (i *AppContext) TmpParamsFolderForNS(namespace string) string {
	return i.lookupOrCreateFolder(namespace, ParamsFolder)
}

func (i *AppContext) TmpSecretsFolderForNS(namespace string) string {
	return i.lookupOrCreateFolder(namespace, SecretsFolder)
}

func (i *AppContext) InstallerFolder() string {
	return i.lookupOrCreateFolder("installer")
}

func (i *AppContext) kustomizeFolderForNS(namespace string) string {
	return i.lookupOrCreateFolder(i.InstallerFolder(), "kustomize", namespace)
}

func (i *AppContext) KustomizationFileFrom(folder string) string {
	return filepath.Join(folder, KustomizationFile)
}

func (i *AppContext) BaseKustomizeFolderForNS(namespace string) string {
	return i.lookupOrCreateFolder(i.kustomizeFolderForNS(namespace), "base")
}

func (i *AppContext) KustomizeTemplateFolderForNS(namespace string) string {
	return i.lookupOrCreateFolder(i.kustomizeFolderForNS(namespace), "template")
}

func (i *AppContext) KustomizeParamsFolderForNS(namespace string) string {
	return i.lookupOrCreateFolder(i.BaseKustomizeFolderForNS(namespace), ParamsFolder)
}

func (i *AppContext) KustomizeTemplateParamsFolderForNS(namespace string) string {
	return i.lookupOrCreateFolder(i.KustomizeTemplateFolderForNS(namespace), ParamsFolder)
}

func (i *AppContext) KustomizeSecretsFolderForNS(namespace string) string {
	return i.lookupOrCreateFolder(i.KustomizeTemplateFolderForNS(namespace), SecretsFolder)
}
