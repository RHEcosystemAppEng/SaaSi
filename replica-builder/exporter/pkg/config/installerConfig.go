package config

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

var ParamsFolder = "params"
var SecretsFolder = "secrets"
var NoValue = "__DEFAULT__"
var MandatoryValue = "__MANDATORY__"
var KustomizationFile = "kustomization.yaml"

type InstallerConfig struct {
	OutputFolder string
	AppFolder    string
}

func NewInstallerConfigFromApplicationConfig(appConfig *ApplicationConfig) *InstallerConfig {
	config := InstallerConfig{}
	pwd, _ := os.Getwd()
	log.Printf("Running from %v", pwd)
	config.OutputFolder = pwd + "/output"
	config.AppFolder = filepath.Join(config.OutputFolder, appConfig.Application.Name)

	return &config
}

func (i *InstallerConfig) ExportFolderForNS(namespace string) string {
	return filepath.Join(i.AppFolder, namespace, "export")
}

func (i *InstallerConfig) TransformFolderForNS(namespace string) string {
	return filepath.Join(i.AppFolder, namespace, "transform")
}

func (i *InstallerConfig) OutputFolderForNS(namespace string) string {
	return filepath.Join(i.AppFolder, namespace, "output")
}

func (i *InstallerConfig) lookupOrCreateFolder(path ...string) string {
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

func (i *InstallerConfig) TmpParamsFolderForNS(namespace string) string {
	return i.lookupOrCreateFolder(namespace, ParamsFolder)
}

func (i *InstallerConfig) TmpSecretsFolderForNS(namespace string) string {
	return i.lookupOrCreateFolder(namespace, SecretsFolder)
}

func (i *InstallerConfig) InstallerFolder() string {
	return i.lookupOrCreateFolder("installer")
}

func (i *InstallerConfig) kustomizeFolderForNS(namespace string) string {
	return i.lookupOrCreateFolder(i.InstallerFolder(), "kustomize", namespace)
}

func (i *InstallerConfig) KustomizationFileFrom(folder string) string {
	return filepath.Join(folder, KustomizationFile)
}

func (i *InstallerConfig) BaseKustomizeFolderForNS(namespace string) string {
	return i.lookupOrCreateFolder(i.kustomizeFolderForNS(namespace), "base")
}

func (i *InstallerConfig) KustomizeTemplateFolderForNS(namespace string) string {
	return i.lookupOrCreateFolder(i.kustomizeFolderForNS(namespace), "template")
}

func (i *InstallerConfig) KustomizeParamsFolderForNS(namespace string) string {
	return i.lookupOrCreateFolder(i.BaseKustomizeFolderForNS(namespace), ParamsFolder)
}

func (i *InstallerConfig) KustomizeTemplateParamsFolderForNS(namespace string) string {
	return i.lookupOrCreateFolder(i.KustomizeTemplateFolderForNS(namespace), ParamsFolder)
}

func (i *InstallerConfig) KustomizeSecretsFolderForNS(namespace string) string {
	return i.lookupOrCreateFolder(i.KustomizeTemplateFolderForNS(namespace), SecretsFolder)
}
