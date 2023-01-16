package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type DeployerConfig struct {
	Deployer ComponentConfig `yaml:"deployer"`
}

type ComponentConfig struct {
	Cluster     Cluster     `yaml:"cluster"`
	Application Application `yaml:"application"`
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
	Name       string       `yaml:"name"`
	Namespaces []Namespaces `yaml:"namespaces"`
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

func ReadDeployerConfig(configFile string) *ComponentConfig {
	yfile, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatal(err)
	}

	config := DeployerConfig{}
	err = yaml.Unmarshal(yfile, &config)
	if err != nil {
		log.Fatal(err)
	}
	return &config.Deployer
}
