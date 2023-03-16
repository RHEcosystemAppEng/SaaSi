package config

import (
	"errors"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"gopkg.in/yaml.v2"
)

type Config struct {
	RootOutputFolder string
	Debug            bool
	exportConfigFile string
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

func ReadConfigFromFlags() *Config {
	defaultRoot, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	defaultOutput := filepath.Join(defaultRoot, "output")

	config := Config{}
	var outputFolder string
	flag.StringVar(&outputFolder, "output-dir", defaultOutput, "Root output folder")
	flag.StringVar(&outputFolder, "o", defaultOutput, "Root output folder (shorthand)")
	flag.StringVar(&config.exportConfigFile, "f", "", "Application configuration file")
	flag.BoolVar(&config.Debug, "debug", false, "Debug the command by printing more information")
	flag.Parse()

	config.RootOutputFolder = outputFolder
	return &config
}

func ReadConfigFromEnvVars() *Config {
	config := Config{}
	v, ok := os.LookupEnv("OUTPUT_DIR")
	if !ok {
		log.Fatal("missing mandatory variable OUTPUT_DIR")
	}
	config.RootOutputFolder = v

	v, ok = os.LookupEnv("DEBUG")
	config.Debug = ok && strings.ToLower(v) == "true"

	return &config
}

func (c *Config) ReadExporterConfig() *ExporterConfig {
	yfile, err := ioutil.ReadFile(c.exportConfigFile)
	if err != nil {
		log.Fatal(err)
	}

	exporterConfig := ExporterConfig{}
	err = yaml.Unmarshal(yfile, &exporterConfig)
	if err != nil {
		log.Fatal(err)
	}
	return &exporterConfig
}

func (c *ClusterConfig) ValidateClusterConfig() error {
	if reflect.ValueOf(c).IsZero() {
		return errors.New("missing cluster configuration")
	} else {
		if c.ClusterId == "" {
			return errors.New("missing clusterId configuration")
		}
		if c.Server == "" {
			return errors.New("missing server configuration")
		}
		if c.Token == "" {
			return errors.New("missing token configuration")
		}
	}
	return nil
}

func (e *ExporterConfig) Validate() error {
	if reflect.ValueOf(e.Cluster).IsZero() {
		return errors.New("missing cluster configuration")
	} else {
		err := e.Cluster.ValidateClusterConfig()
		if err != nil {
			return err
		}
	}
	if reflect.ValueOf(e.Application).IsZero() {
		return errors.New("missing application configuration")
	} else {
		if e.Application.Name == "" {
			return errors.New("missing application name configuration")
		}
	}
	return nil
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
