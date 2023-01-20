package config

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Exporter   ExporterConfig `yaml:"exporter"`
	RootFolder string
}

type ExporterConfig struct {
	Cluster     ClusterConfig     `yaml:"cluster"`
	Application ApplicationConfig `yaml:"application"`
}

type ClusterConfig struct {
	ClusterId string `yaml:"clusterId"`
	Server    string `yaml:"server"`
	Token     string `yaml:"token"`
}

type ApplicationConfig struct {
	Name       string            `yaml:"name"`
	Namespaces []SourceNamespace `yaml:"namespaces"`
}

type SourceNamespace struct {
	Name            string           `yaml:"name"`
	MandatoryParams []MandatoryParam `yaml:"mandatory-params"`
}

type MandatoryParam struct {
	ConfigMap string   `yaml:"configMap"`
	Params    []string `yaml:"params"`
}

func ReadConfig() *Config {
	defaultRoot, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	config := Config{}
	configFile := flag.String("f", "", "Application configuration file")
	rootFolder := flag.String("root", defaultRoot, "Root installation folder")
	flag.Parse()

	yfile, err := ioutil.ReadFile(*configFile)
	if err != nil {
		log.Fatal(err)
	}

	err = yaml.Unmarshal(yfile, &config)
	if err != nil {
		log.Fatal(err)
	}

	config.RootFolder = *rootFolder
	return &config
}

func (c *ApplicationConfig) MandatoryParamsByNSAndConfigMap(namespace string, configMap string) []string {
	for _, ns := range c.Namespaces {
		if ns.Name == namespace {
			for _, params := range ns.MandatoryParams {
				if params.ConfigMap == configMap {
					return params.Params
				}
			}
		}
	}
	return []string{}
}
