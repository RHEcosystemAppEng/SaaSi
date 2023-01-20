package app

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/RHEcosystemAppEng/SaaSi/exporter/pkg/config"
	"github.com/RHEcosystemAppEng/SaaSi/exporter/pkg/connect"
	"k8s.io/client-go/rest"
)

var ParamsFolder = "params"
var SecretsFolder = "secrets"
var NoValue = "__DEFAULT__"
var MandatoryValue = "__MANDATORY__"
var KustomizationFile = "kustomization.yaml"

type AppContext struct {
	AppConfig        *config.ApplicationConfig
	connectionStatus *connect.ConnectionStatus

	OutputFolder string
	AppFolder    string
}

func NewAppContextFromConfig(config *config.Config, connectionStatus *connect.ConnectionStatus) *AppContext {
	context := AppContext{connectionStatus: connectionStatus, AppConfig: &config.Exporter.Application}

	pwd, _ := os.Getwd()
	log.Printf("Running from %v", pwd)
	context.OutputFolder = pwd + "/output"
	// TODO
	context.AppFolder = filepath.Join(context.OutputFolder, context.AppConfig.Name)

	return &context
}

func (c *AppContext) KubeConfig() *rest.Config {
	return c.connectionStatus.KubeConfig
}

func (c *AppContext) KubeConfigPath() string {
	return c.connectionStatus.KubeConfigPath
}

func (c *AppContext) ExportFolderForNS(namespace string) string {
	return filepath.Join(c.AppFolder, namespace, "export")
}

func (c *AppContext) TransformFolderForNS(namespace string) string {
	return filepath.Join(c.AppFolder, namespace, "transform")
}

func (c *AppContext) OutputFolderForNS(namespace string) string {
	return filepath.Join(c.AppFolder, namespace, "output")
}

func (c *AppContext) lookupOrCreateFolder(path ...string) string {
	fullPath := ""
	if strings.HasPrefix(path[0], c.AppFolder) {
		fullPath = filepath.Join(path...)
	} else {
		fullPath = filepath.Join(append([]string{c.AppFolder}, path...)...)
	}
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		if err := os.MkdirAll(fullPath, os.ModePerm); err != nil {
			log.Fatalf("Cannot create %v folder: %v", fullPath, err)
		}
	}
	return fullPath
}

func (c *AppContext) TmpParamsFolderForNS(namespace string) string {
	return c.lookupOrCreateFolder(namespace, ParamsFolder)
}

func (c *AppContext) TmpSecretsFolderForNS(namespace string) string {
	return c.lookupOrCreateFolder(namespace, SecretsFolder)
}

func (c *AppContext) InstallerFolder() string {
	return c.lookupOrCreateFolder("installer")
}

func (c *AppContext) kustomizeFolderForNS(namespace string) string {
	return c.lookupOrCreateFolder(c.InstallerFolder(), "kustomize", namespace)
}

func (c *AppContext) KustomizationFileFrom(folder string) string {
	return filepath.Join(folder, KustomizationFile)
}

func (c *AppContext) BaseKustomizeFolderForNS(namespace string) string {
	return c.lookupOrCreateFolder(c.kustomizeFolderForNS(namespace), "base")
}

func (c *AppContext) KustomizeTemplateFolderForNS(namespace string) string {
	return c.lookupOrCreateFolder(c.kustomizeFolderForNS(namespace), "template")
}

func (c *AppContext) KustomizeParamsFolderForNS(namespace string) string {
	return c.lookupOrCreateFolder(c.BaseKustomizeFolderForNS(namespace), ParamsFolder)
}

func (c *AppContext) KustomizeTemplateParamsFolderForNS(namespace string) string {
	return c.lookupOrCreateFolder(c.KustomizeTemplateFolderForNS(namespace), ParamsFolder)
}

func (c *AppContext) KustomizeSecretsFolderForNS(namespace string) string {
	return c.lookupOrCreateFolder(c.KustomizeTemplateFolderForNS(namespace), SecretsFolder)
}
