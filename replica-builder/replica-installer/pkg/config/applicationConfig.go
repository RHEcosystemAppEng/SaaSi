package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type ApplicationConfig struct {
	Application Application `yaml:"application"`
}

type Application struct {
	Name       string            `yaml:"name"`
	Namespaces []SourceNamespace `yaml:"namespaces"`
}

type SourceNamespace struct {
	Name            string           `yaml:"name"`
	Target			string			 `yaml:"target"`
	// MandatoryParams []MandatoryParam `yaml:"mandatory-params"`
}

// type MandatoryParam struct {
// 	ConfigMap string   `yaml:"configMap"`
// 	Params    []string `yaml:"params"`
// }

func ReadApplicationConfig(configFile string) *ApplicationConfig {
	yfile, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatal(err)
	}

	config := ApplicationConfig{}
	err = yaml.Unmarshal(yfile, &config)
	if err != nil {
		log.Fatal(err)
	}
	return &config
}

// func (a *ApplicationConfig) MandatoryParamsByNSAndConfigMap(namespace string, configMap string) []string {
// 	for _, ns := range a.Application.Namespaces {
// 		if ns.Name == namespace {
// 			for _, params := range ns.MandatoryParams {
// 				if params.ConfigMap == configMap {
// 					return params.Params
// 				}
// 			}
// 		}
// 	}
// 	return []string{}
// }
