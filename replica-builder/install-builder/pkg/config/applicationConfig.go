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
	Name            string          `yaml:"name"`
	MandatoryParams MandatoryParams `yaml:"mandatory-params"`
}

type MandatoryParams struct {
	ConfigMaps []MandatoryParam `yaml:"configMaps"`
	Secrets    []MandatoryParam `yaml:"secrets"`
}

type MandatoryParam struct {
	ContainerName string `yaml:"name"`
	ParamName     string `yaml:"param"`
}

func ReadConfig(configFile string) *ApplicationConfig {
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
