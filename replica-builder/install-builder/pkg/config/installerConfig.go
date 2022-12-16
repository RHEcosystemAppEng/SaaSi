package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var ParamsFolder = "params"
var SecretsFolder = "secrets"
var NoSecretValue = "__EMPTY__"

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

func (i *InstallerConfig) TmpParamsFolderForNS(namespace string) string {
	tmpParamsFolder := filepath.Join(i.AppFolder, namespace, ParamsFolder)
	if _, err := os.Stat(tmpParamsFolder); os.IsNotExist(err) {
		if err := os.Mkdir(tmpParamsFolder, os.ModePerm); err != nil {
			log.Fatalf("Cannot create %v folder: %v", tmpParamsFolder, err)
		}
	}
	return tmpParamsFolder
}

func (i *InstallerConfig) TmpSecretsFolderForNS(namespace string) string {
	tmpSecretsFolder := filepath.Join(i.AppFolder, namespace, SecretsFolder)
	if _, err := os.Stat(tmpSecretsFolder); os.IsNotExist(err) {
		if err := os.Mkdir(tmpSecretsFolder, os.ModePerm); err != nil {
			log.Fatalf("Cannot create %v folder: %v", tmpSecretsFolder, err)
		}
	}
	return tmpSecretsFolder
}

func (i *InstallerConfig) TmpParamsFolderForConfigMap(namespace string, configmap string) string {
	tmpParamsFolder := filepath.Join(i.TmpParamsFolderForNS(namespace), configmap)
	if _, err := os.Stat(tmpParamsFolder); os.IsNotExist(err) {
		if err := os.Mkdir(tmpParamsFolder, os.ModePerm); err != nil {
			log.Fatalf("Cannot create %v folder: %v", tmpParamsFolder, err)
		}
	}
	return tmpParamsFolder
}

func (i *InstallerConfig) InstallerFolder() string {
	installerFolder := filepath.Join(i.AppFolder, "installer")
	if _, err := os.Stat(installerFolder); os.IsNotExist(err) {
		if err := os.Mkdir(installerFolder, os.ModePerm); err != nil {
			log.Fatalf("Cannot create %v folder: %v", installerFolder, err)
		}
	}
	return installerFolder
}

func (i *InstallerConfig) kustomizeFolderForNS(namespace string) string {
	return filepath.Join(i.InstallerFolder(), namespace, "source", "resources", fmt.Sprintf("%s-versionchanged-parameterized", namespace),
		"kustomize")
}

func (i *InstallerConfig) BaseKustomizeFolderForNS(namespace string) string {
	return filepath.Join(i.kustomizeFolderForNS(namespace), "base")
}

func (i *InstallerConfig) KustomizeOverlaysFolderForNS(namespace string) string {
	return filepath.Join(i.kustomizeFolderForNS(namespace), "overlays")
}

func (i *InstallerConfig) KustomizationFileFrom(root string) string {
	return filepath.Join(root, "kustomization.yaml")
}

func (i *InstallerConfig) KustomizeParamsFolderForNS(namespace string) string {
	return filepath.Join(i.BaseKustomizeFolderForNS(namespace), ParamsFolder)
}

func (i *InstallerConfig) KustomizeSecretsFolderForNS(namespace string) string {
	return filepath.Join(i.KustomizeOverlaysFolderForNS(namespace), SecretsFolder)
}
