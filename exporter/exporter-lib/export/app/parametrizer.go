package app

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"

	api "github.com/openshift/api"
	"github.com/sirupsen/logrus"

	"github.com/RHEcosystemAppEng/SaaSi/exporter/exporter-lib/export/utils"
	"golang.org/x/exp/slices"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
)

type Parametrizer struct {
	appContext *AppContext
	logger     *logrus.Logger
}

func NewParametrizerFromConfig(appContext *AppContext) *Parametrizer {
	parametrizer := Parametrizer{appContext: appContext}
	parametrizer.logger = appContext.Logger()
	return &parametrizer
}

func (p *Parametrizer) ExposeParameters() error {
	for _, ns := range p.appContext.AppConfig.Namespaces {
		p.logger.Infof("Exposing parameters for NS %s", ns.Name)
		outputFolder := p.appContext.OutputFolderForNS(ns.Name)
		var yamlFiles []string
		filepath.WalkDir(outputFolder, func(s string, d fs.DirEntry, e error) error {
			if e != nil {
				return e
			}
			if filepath.Ext(d.Name()) == ".yaml" {
				yamlFiles = append(yamlFiles, s)
			}
			return nil
		})

		// Install openshift schemes
		api.InstallKube(scheme.Scheme)
		api.Install(scheme.Scheme)
		decode := scheme.Codecs.UniversalDeserializer().Decode

		for _, yamlFile := range yamlFiles {
			yfile, err := ioutil.ReadFile(yamlFile)
			if err != nil {
				return err
			}

			obj, gKV, err := decode(yfile, nil, nil)
			if err != nil {
				p.logger.Warnf("Cannot decode %s to Kubernetes resource: %s", yfile, err.Error())
				continue
			} else {
				if gKV.Kind == "ConfigMap" {
					configMap := obj.(*v1.ConfigMap)
					err := p.handleConfigMap(yamlFile, configMap)
					if err != nil {
						return err
					}
				} else if gKV.Kind == "Secret" {
					secret := obj.(*v1.Secret)
					p.handleSecret(yamlFile, secret)
				}
			}
		}
	}
	return nil
}

func (p *Parametrizer) handleConfigMap(configMapFile string, configMap *v1.ConfigMap) error {
	p.logger.Infof("Handling ConfigMap %s", configMap.Name)
	tmpParamsFolder := p.appContext.TmpParamsFolderForNS(configMap.Namespace)
	// RunCommand("oc", "extract", "-f", configMapFile, "--to", tmpParamsFolder)

	mandatoryParams := p.appContext.AppConfig.MandatoryParamsByNSAndConfigMap(configMap.Namespace, configMap.Name)

	templateFile := filepath.Join(tmpParamsFolder, fmt.Sprintf("%s.env", configMap.Name))
	os.Create(templateFile)
	for key := range configMap.Data {
		if slices.Contains(mandatoryParams, key) {
			utils.AppendToFile(templateFile, fmt.Sprintf("%s=%s\n", key, MandatoryValue))
		} else {
			utils.AppendToFile(templateFile, fmt.Sprintf("#%s=%s\n", key, NoValue))
		}
	}

	for _, mandatoryParam := range mandatoryParams {
		if _, ok := configMap.Data[mandatoryParam]; ok {
			p.logger.Infof("Removing mandatory param %s from %s", mandatoryParam, configMap.Name)
			delete(configMap.Data, mandatoryParam)
		} else {
			return fmt.Errorf("the mandatory parameter %s for ConfigMap %s does not exist", mandatoryParam, configMap.Name)
		}
	}
	// os.Rename(configMapFile, BackupFile(configMapFile))
	return nil
}

func (p *Parametrizer) handleSecret(secretFile string, secret *v1.Secret) error {
	p.logger.Infof("Handling Secret %s", secret.Name)
	if secret.Type == "Opaque" {
		tmpSecretsFolder := p.appContext.TmpSecretsFolderForNS(secret.Namespace)
		secretsFile := filepath.Join(tmpSecretsFolder, fmt.Sprintf("%s.env", secret.Name))
		os.Create(secretsFile)
		p.logger.Infof("Creating secret configuration template %s", secretsFile)

		for key, _ := range secret.Data {
			utils.AppendToFile(secretsFile, fmt.Sprintf("%s=%s\n", key, MandatoryValue))
		}
		err := os.Rename(secretFile, utils.BackupFile(secretFile))
		if err != nil {
			return err
		}
	}
	return nil
}
