package config

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Exporter               ExporterConfig `yaml:"exporter"`
	RootInstallationFolder string
	RootOutputFolder       string
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
	defaultOutput := filepath.Join(defaultRoot, "output")

	config := Config{}
	configFile := flag.String("f", "", "Application configuration file")
	var rootFolder string
	flag.StringVar(&rootFolder, "install-dir", defaultRoot, "Root installation folder")
	flag.StringVar(&rootFolder, "i", defaultRoot, "Root installation folder (shorthand)")
	var outputFolder string
	flag.StringVar(&outputFolder, "output-dir", defaultOutput, "Root output folder")
	flag.StringVar(&outputFolder, "o", defaultOutput, "Root output folder (shorthand)")
	flag.Parse()

	yfile, err := ioutil.ReadFile(*configFile)
	if err != nil {
		log.Fatal(err)
	}

	err = yaml.Unmarshal(yfile, &config)
	if err != nil {
		log.Fatal(err)
	}

	config.RootInstallationFolder = rootFolder
	config.RootOutputFolder = outputFolder
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
