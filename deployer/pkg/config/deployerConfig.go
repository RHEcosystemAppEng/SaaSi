package config

import (
	"io/ioutil"
	"log"

	"github.com/jessevdk/go-flags"
	"gopkg.in/yaml.v2"
)

type FlagArgs struct {
	ConfigFile    string `short:"f" description:"Application configuration file for deployemnt" required:"true"`
	RootOutputDir string `short:"o" long:"output-dir" default:"output" description:"Root output folder"`
	RootSourceDir string `short:"s" long:"source-dir" description:"Root source folder" required:"true"`
}

// ----------------------
// ----Deployer Config----
// ----------------------

type DeployerConfig struct {
	Deployer ComponentConfig `yaml:"deployer"`
}

type ComponentConfig struct {
	ClusterConfig     ClusterConfig     `yaml:"cluster"`
	ApplicationConfig ApplicationConfig `yaml:"application"`
	FlagArgs          FlagArgs
}

// ----------------------
// ----Cluster Config----
// ----------------------

type ClusterConfig struct {
	Server        string        `yaml:"server"`
	User          string        `yaml:"user"`
	Token         string        `yaml:"token"`
	FromClusterId string        `yaml:"fromClusterId"`
	ClusterId     string        `yaml:"clusterId"`
	Aws           AwsSettings   `yaml:"aws"`
	Params        ClusterParams `yaml:"params"`
}

type AwsSettings struct {
	AwsPublicDomain    string `yaml:"aws_public_domain"`
	AwsAccountName     string `yaml:"aws_account_name"`
	AwsAccessKeyId     string `yaml:"aws_access_key_id"`
	AwsSecretAccessKey string `yaml:"aws_secret_access_key"`
}

type ClusterParams struct {
	ClusterName       string `yaml:"CLUSTER_NAME"`
	ClusterBaseDomain string `yaml:"CLUSTER_BASE_DOMAIN"`
	WorkerCount       string `yaml:"WORKER_COUNT"`
}

// ----------------------
// ------App Config------
// ----------------------

type ApplicationConfig struct {
	Name                   string       `yaml:"name"`
	NamespaceMappingFormat string       `yaml:"namespaceMappingFormat"`
	Namespaces             []Namespaces `yaml:"namespaces"`
}

type Namespaces struct {
	Name       string       `yaml:"name"`
	Target     string       `yaml:"target"`
	ConfigMaps []ConfigMaps `yaml:"params"`
	Secrets    []Secrets    `yaml:"secrets"`
}

type ConfigMaps struct {
	ConfigMap string              `yaml:"configMap"`
	Params    []ApplicationParams `yaml:"params"`
}

type Secrets struct {
	Secret string              `yaml:"secret"`
	Params []ApplicationParams `yaml:"params"`
}

type ApplicationParams struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

func InitDeployerConfig() *ComponentConfig {

	// init flags for input arguments
	flagArgs := FlagArgs{}

	// parse input arguments from os into flags
	_, err := flags.Parse(&flagArgs)
	if err != nil {
		log.Fatal("Failed to parse os input arguments")
	}

	// read deployer config file
	yfile, err := ioutil.ReadFile(flagArgs.ConfigFile)
	if err != nil {
		log.Fatal(err)
	}

	// Unmarshal deployer config
	config := DeployerConfig{}
	err = yaml.Unmarshal(yfile, &config)
	if err != nil {
		log.Fatal(err)
	}

	// set flag input arguments for components
	config.Deployer.FlagArgs = flagArgs

	return &config.Deployer
}
