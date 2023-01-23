# exporter
A Golang CLI tool to extract configurations from a live OpenShit/Kubernetes environment and generate a reusable, configurable
installer for the [deployer](../deployer/README.md) tool.

## Dependencies
* [Konveyor Crane](https://konveyor.io/tools/crane/) [Golang](https://go.dev/) packages to export the original configuration and remove cluster specific settings 
(e.g. IP addresses, status, ...)

## Features
* The generated installer can replicate the resources of the original namespaces
* The generated installer allows to override all the application parameters defined in the `ConfigMap` and `Secrets` (encrypted)
  * Original values from all `ConfigMap`s are used as defaults but can be overridden individually
    * Mandatory properties will be identified and their values will not be copied to be reused but will have to be overridden at 
    installation time 
  * Original values from all `Secret`s are hidden and must be overridden at installation time
* The generated installer is agnostic from the original namespaces
* The installer also exports cluster-wide resources like the `ClusterRoleBindings` to 

## Feature design
* Apply [Golang](https://go.dev/) packages imported from the `Crane` migration tool (under the `Konveyor` community) to export and normalize the original configuration
* Post-execution manipulations are performed
  * Clear the reference to the original namespace
  * Create a sample `template` for the overlayed configuration where the developer can apply customizations
    * All the keys in any given `ConfigMap` can be overridden in properties files called `CONFIGMAP_NAME.env`
      * The skeleton of these files are automatically generated
      * All the key names for non-mandatory parameters are commented out and set to `__DEFAULT__` value
      * For mandatory parameters, the keys are set to `__MANDATORY__` value 
    * All the keys in the managed `Secret`s are customizable in the same way
      * The `template` overlay re-creates the `Secret`s from the template files called `SECRET_NAME.env`
      * The base `kustomize` configuration does not re-create the `Secret`s, so its deployment would actually fail
      * Errors must be raised while trying to install the default template for the secrets

## Builder configuration
The `exporter` runs using a configuration that specifies the desired behavior: 
```yaml
exporter:
  cluster:
    # Must be unique
    clusterId: UNIQUE-ID
    server: API-SERVER-URL
    # Must be valid at the moment we export the configuration
    token: TOKEN
  application:
    # This creates an installer package named APP
    name: APP
    namespaces:
    - name: NS1
      # No default values are generated for each of the following mandatory params
      mandatory-params:
      # Provide the name of one of the exported ConfigMaps
      - configMap: MAP-1
        params:
        # Paramaters are given as a list of key names  
        - PARAM-1
        - "..."
        - PARAM-N
```

## Running the builder
Prerequisites:
* Install `oc` CLI tool
* Install `go` CLI tool (at least version 1.19)
* `oc` login to the source OpenShift cluster

Run this command to create the installer from the given configuration `myapp.yaml`:
```bash
go run main.go -f myapp.yaml
```
The command creates an `output/<APP NAME>/installer` folder in the current directory with the whole installer package.

Available options:
```bash
> go run main.go --help
  -f string
        Application configuration file
  -i string
        Root installation folder (shorthand) (default "<CURRENT_DIR>")
  -install-dir string
        Root installation folder (default "<CURRENT_DIR>")
  -o string
        Root output folder (shorthand) (default "<CURRENT_DIR>/output")
  -output-dir string
        Root output folder (default "<CURRENT_DIR>/output")
```

### Output specification
Under the configured output folder the exporter creates:
* one `clusters` folder:
  * one folder for each exported cluster configuration, using the `clusterId` specified in the `cluster` configuration
* one `applications` folder:
  * one folder for each exported application, using the `name` specified in the `application` configuration
  * one subfolder for each managed namespace configured in `namespaces`, with `base` and `template` configurations for `kustomize`

```bash
├── output
│   └── clusters
│       └── CLUSTER_ID1
│           └── CLUSTER_ID1.env
...
│   └── applications
│       └── APPLICATION1
│           └── kustomize
│               └── NAMESPACE1
│                   ├── base
│                   │   ├── RESOURCE1.yaml
│                   │   ├── RESOURCE2.yaml
│                   │   └── kustomization.yaml
│                   └── template
│                       ├── kustomization.yaml
│                       ├── params
│                       │   ├── CONFIG_MAP1.env
│                       │   └── CONFIG_MAP2.env
│                       └── secrets
│                           └── SECRET1.env
```

The `base` kustomization contains the resources extracted from the source cluster, stripped of the status information,
cluster specific settings and namespace configuration.

The `template` overlay is an example of a possible `kustomize` overlay, with skeletons of environment files to override the
`ConfigMap`s parameters (remove the `#` comment and set the desired value) and all the `Secret` values:
```bash
> cat output/Infinity/kustomize/holdings/template/params/CONFIG_MAP1.env
#KEY1=__DEFAULT__
#KEY2=__DEFAULT__
KEY3=__MANDATORY__

> cat output/Infinity/kustomize/holdings/template/secrets/SECRET1.env
KEY1=__DEFAULT__
```

## Customize and install the template
Simple procedure that will be automated using the [deployer](./deployer/README.md) tool.

```bash
cd output/APPLICATION/installer/kustomize/NAMESPACE
cp -r template MYCONFIG
cd MYCONFIG
# Apply namespaxce update (must be created first)
kustomize edit set namespace MYNAMESPACE
# Edit changes to params/*.env and secrets/*env
kustomize build . | oc apply -f-
```

## Issues
GitHub [exporter issues](https://github.com/RHEcosystemAppEng/SaaSi/issues?q=is%3Aopen+is%3Aissue+label%3Aexporter)