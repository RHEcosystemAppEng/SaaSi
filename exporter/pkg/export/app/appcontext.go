package app

import (
	"path/filepath"

	"github.com/RHEcosystemAppEng/SaaSi/exporter/pkg/config"
	"github.com/RHEcosystemAppEng/SaaSi/exporter/pkg/connect"
	"github.com/RHEcosystemAppEng/SaaSi/exporter/pkg/context"
)

const (
	ApplicationsFolder = "applications"
	CraneFolder        = "crane"
	SecretsFolder      = "secrets"
	KustomizeFolder    = "kustomize"
	BaseFolder         = "base"
	TemplateFolder     = "template"
	ParamsFolder       = "params"
	NoValue            = "__DEFAULT__"
	MandatoryValue     = "__MANDATORY__"
	KustomizationFile  = "kustomization.yaml"
)

type AppContext struct {
	context.ExporterContext

	AppConfig *config.ApplicationConfig
	AppFolder string
}

func NewAppContextFromConfig(config *config.Config, connectionStatus *connect.ConnectionStatus) *AppContext {
	appContext := AppContext{AppConfig: &config.Exporter.Application}

	appContext.InitFromConfig(config, connectionStatus)
	appContext.AppFolder = filepath.Join(appContext.OutputFolder, ApplicationsFolder, appContext.AppConfig.Name)

	return &appContext
}
func (c *AppContext) RootFolder() string {
	return c.AppFolder
}

func (c *AppContext) ExportFolderForNS(namespace string) string {
	return filepath.Join(c.AppFolder, CraneFolder, namespace, "export")
}

func (c *AppContext) TransformFolderForNS(namespace string) string {
	return filepath.Join(c.AppFolder, CraneFolder, namespace, "transform")
}

func (c *AppContext) OutputFolderForNS(namespace string) string {
	return filepath.Join(c.AppFolder, CraneFolder, namespace, "output")
}

func (c *AppContext) TmpParamsFolderForNS(namespace string) string {
	return context.LookupOrCreateFolder(c, namespace, ParamsFolder)
}

func (c *AppContext) TmpSecretsFolderForNS(namespace string) string {
	return context.LookupOrCreateFolder(c, namespace, SecretsFolder)
}

func (c *AppContext) kustomizeFolderForNS(namespace string) string {
	return context.LookupOrCreateFolder(c, KustomizeFolder, namespace)
}

func (c *AppContext) KustomizationFileFrom(folder string) string {
	return filepath.Join(folder, KustomizationFile)
}

func (c *AppContext) BaseKustomizeFolderForNS(namespace string) string {
	return context.LookupOrCreateFolder(c, c.kustomizeFolderForNS(namespace), BaseFolder)
}

func (c *AppContext) KustomizeTemplateFolderForNS(namespace string) string {
	return context.LookupOrCreateFolder(c, c.kustomizeFolderForNS(namespace), TemplateFolder)
}

func (c *AppContext) KustomizeParamsFolderForNS(namespace string) string {
	return context.LookupOrCreateFolder(c, c.BaseKustomizeFolderForNS(namespace), ParamsFolder)
}

func (c *AppContext) KustomizeTemplateParamsFolderForNS(namespace string) string {
	return context.LookupOrCreateFolder(c, c.KustomizeTemplateFolderForNS(namespace), ParamsFolder)
}

func (c *AppContext) KustomizeSecretsFolderForNS(namespace string) string {
	return context.LookupOrCreateFolder(c, c.KustomizeTemplateFolderForNS(namespace), SecretsFolder)
}
