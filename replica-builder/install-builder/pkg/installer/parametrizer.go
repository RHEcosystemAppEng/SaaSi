package installer

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"

	"k8s.io/cli-runtime/pkg/printers"
	"k8s.io/client-go/kubernetes/scheme"

	"github.com/RHEcosystemAppEng/SaaSi/replica-builder/install-builder/pkg/config"
	v1 "k8s.io/api/core/v1"
)

type Parametrizer struct {
	appConfig       *config.ApplicationConfig
	installerConfig *config.InstallerConfig
}

func NewParametrizerFromConfig(appConfig *config.ApplicationConfig, installerConfig *config.InstallerConfig) *Parametrizer {
	parametrizer := Parametrizer{appConfig: appConfig, installerConfig: installerConfig}
	return &parametrizer
}

func (p *Parametrizer) ExposeParameters() {
	for _, ns := range p.appConfig.Application.Namespaces {
		log.Printf("Exposing parameters for NS %s", ns.Name)
		outputFolder := p.installerConfig.OutputFolderForNS(ns.Name)
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

		for _, yamlFile := range yamlFiles {
			yfile, err := ioutil.ReadFile(yamlFile)
			if err != nil {
				log.Fatal(err)
			}

			decode := scheme.Codecs.UniversalDeserializer().Decode
			obj, gKV, err := decode(yfile, nil, nil)
			if err != nil {
				workaround := false
				for _, k := range []string{"RoleBinding", "Route"} {
					if FileContains(yamlFile, fmt.Sprintf("kind: %s", k)) {
						// TODO
						// 1- Skip RoleBindings that are managed by controllers
						// 2- Manage OpenShift resources like Route
						log.Printf("Skipping %s resource %s", k, yamlFile)
						ReplaceInFile(yamlFile, fmt.Sprintf("namespace: %s", ns.Name), "# Removed NS setting")
						workaround = true
						break
					}
				}
				if !workaround {
					log.Fatal(err)
				}
			} else {
				if gKV.Kind == "ConfigMap" {
					configMap := obj.(*v1.ConfigMap)
					p.extractConfigMap(yamlFile, configMap)
				} else if gKV.Kind == "Secret" {
					secret := obj.(*v1.Secret)
					p.handleSecret(yamlFile, secret)
				} else {
					value := reflect.Indirect(reflect.ValueOf(obj))
					ns := value.FieldByName("Namespace")
					log.Printf("yamlFile is %s", yamlFile)
					log.Printf("ns is %+v", ns)
					if !ns.IsZero() {
						namespace := reflect.Indirect(ns).String()
						name := value.FieldByName("Name").String()

						log.Printf("Resetting namespace %s at %s/%s", namespace, gKV.Kind, name)
						ns.SetString("")

						os.Rename(yamlFile, BackupFile(yamlFile))
						newFile, err := os.Create(yamlFile)
						if err != nil {
							log.Fatal(err)
						}
						y := printers.YAMLPrinter{}
						defer newFile.Close()
						y.PrintObj(obj, newFile)
					} else {
						name := value.FieldByName("Name").String()
						log.Printf("Found not namespaced resource %s/%s", gKV.Kind, name)
					}
				}
			}
		}
	}
}

func (p *Parametrizer) extractConfigMap(configMapFile string, configMap *v1.ConfigMap) {
	log.Printf("Extracting ConfigMap %s", configMap.Name)
	tmpParamsFolder := p.installerConfig.TmpParamsFolderForConfigMap(configMap.Namespace, configMap.Name)
	RunCommand("oc", "extract", "-f", configMapFile, "--to", tmpParamsFolder)

	templateFile := filepath.Join(tmpParamsFolder, fmt.Sprintf("%s.env", configMap.Name))
	os.Create(templateFile)
	for key := range configMap.Data {
		AppendToFile(templateFile, fmt.Sprintf("#%s=%s\n", key, config.NoValue))
	}

	os.Rename(configMapFile, BackupFile(configMapFile))
}

func (p *Parametrizer) handleSecret(secretFile string, secret *v1.Secret) {
	log.Printf("Handling Secret %s", secret.Name)
	if secret.Type != "Opaque" {
		log.Printf("Removing non-Opaque Secret %s", secret.Name)
	} else {
		tmpSecretsFolder := p.installerConfig.TmpSecretsFolderForNS(secret.Namespace)
		secretsFile := filepath.Join(tmpSecretsFolder, fmt.Sprintf("%s.env", secret.Name))
		os.Create(secretsFile)
		log.Printf("Creating secret configuration template %s", secretsFile)

		for key, _ := range secret.Data {
			AppendToFile(secretsFile, fmt.Sprintf("%s=%s\n", key, config.NoValue))
		}
	}
	os.Rename(secretFile, BackupFile(secretFile))
}
