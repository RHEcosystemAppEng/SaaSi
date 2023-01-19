package export

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/RHEcosystemAppEng/SaaSi/exporter/pkg/config"
)

var ParamsFolder = "params"
var SecretsFolder = "secrets"
var NoValue = "__DEFAULT__"
var MandatoryValue = "__MANDATORY__"
var KustomizationFile = "kustomization.yaml"

type Context struct {
	OutputFolder string
	AppFolder    string
}

func NewContextFromConfig(config *config.Config) *Context {
	context := Context{}
	pwd, _ := os.Getwd()
	log.Printf("Running from %v", pwd)
	context.OutputFolder = pwd + "/output"
	// TODO
	context.AppFolder = filepath.Join(context.OutputFolder, config.Exporter.Application.Name)

	return &context
}

func (i *Context) ExportFolderForNS(namespace string) string {
	return filepath.Join(i.AppFolder, namespace, "export")
}

func (i *Context) TransformFolderForNS(namespace string) string {
	return filepath.Join(i.AppFolder, namespace, "transform")
}

func (i *Context) OutputFolderForNS(namespace string) string {
	return filepath.Join(i.AppFolder, namespace, "output")
}

func (i *Context) lookupOrCreateFolder(path ...string) string {
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

func (i *Context) TmpParamsFolderForNS(namespace string) string {
	return i.lookupOrCreateFolder(namespace, ParamsFolder)
}

func (i *Context) TmpSecretsFolderForNS(namespace string) string {
	return i.lookupOrCreateFolder(namespace, SecretsFolder)
}

func (i *Context) InstallerFolder() string {
	return i.lookupOrCreateFolder("installer")
}

func (i *Context) kustomizeFolderForNS(namespace string) string {
	return i.lookupOrCreateFolder(i.InstallerFolder(), "kustomize", namespace)
}

func (i *Context) KustomizationFileFrom(folder string) string {
	return filepath.Join(folder, KustomizationFile)
}

func (i *Context) BaseKustomizeFolderForNS(namespace string) string {
	return i.lookupOrCreateFolder(i.kustomizeFolderForNS(namespace), "base")
}

func (i *Context) KustomizeTemplateFolderForNS(namespace string) string {
	return i.lookupOrCreateFolder(i.kustomizeFolderForNS(namespace), "template")
}

func (i *Context) KustomizeParamsFolderForNS(namespace string) string {
	return i.lookupOrCreateFolder(i.BaseKustomizeFolderForNS(namespace), ParamsFolder)
}

func (i *Context) KustomizeTemplateParamsFolderForNS(namespace string) string {
	return i.lookupOrCreateFolder(i.KustomizeTemplateFolderForNS(namespace), ParamsFolder)
}

func (i *Context) KustomizeSecretsFolderForNS(namespace string) string {
	return i.lookupOrCreateFolder(i.KustomizeTemplateFolderForNS(namespace), SecretsFolder)
}
