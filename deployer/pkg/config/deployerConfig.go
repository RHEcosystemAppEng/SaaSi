package config

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

const (
	DEFAULT_OUTPUT_DIR = "output"
)

type DeployerConfig struct {
	Deployer ComponentConfig `yaml:"deployer"`
}

type ComponentConfig struct {
	Cluster       Cluster     `yaml:"cluster"`
	Application   Application `yaml:"application"`
	RootOutputDir string
}

// ----------------------
// ----Cluster Config----
// ----------------------

type Cluster struct {
	Server      string        `yaml:"server"`
	User        string        `yaml:"user"`
	Token       string        `yaml:"token"`
	FromCluster string        `yaml:"fromCluster"`
	UserName    string        `yaml:"userName"`
	Aws         AwsSettings   `yaml:"aws"`
	Params      ClusterParams `yaml:"params"`
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

type Application struct {
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

func ReadDeployerConfig() *ComponentConfig {

	// define config file flag
	configFile := flag.String("f", "", "Application configuration file for deployemnt")

	// define default output directory
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get root directory: %s", err)
	}
	defaultRootOutputDir := filepath.Join(pwd, DEFAULT_OUTPUT_DIR)

	// define output directory flag
	var rootOutputDir string
	flag.StringVar(&rootOutputDir, "output-dir", defaultRootOutputDir, "Root output folder")
	flag.StringVar(&rootOutputDir, "o", defaultRootOutputDir, "Root output folder (shorthand)")

	// get os arguments
	flag.Parse()

	// read deployer config file
	yfile, err := ioutil.ReadFile(*configFile)
	if err != nil {
		log.Fatal(err)
	}

	// Unmarshal deployer config
	config := DeployerConfig{}
	err = yaml.Unmarshal(yfile, &config)
	if err != nil {
		log.Fatal(err)
	}

	// set output directory for components
	config.Deployer.RootOutputDir = rootOutputDir

	return &config.Deployer
}
