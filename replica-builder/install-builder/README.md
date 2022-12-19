# install-builder
A Golang CLI tool to extract configurations from a live OpenShit/Kubernetes environment and generate a reusable, configurable
installer for the [replica-installer](../replica-installer/README.md) tool.

## Dependencies
* Use [Konveyor crane](https://konveyor.io/tools/crane/) to export the original configuration and remove cluster specific settings 
(e.g. IP addresses, status, ...)
* Use [Konveyor move2kube](https://move2kube.konveyor.io/) to generate the [kustomize](https://kustomize.io/) installer 
for the given environment 

## Features
* The generated installer can replicate the resources of the original namespaces
* The generated installer allows to override all the application parameters defined in the `ConfigMap` and `Secrets` (encrypted)
  * Original values fom all `ConfigMap`s are used as defaults but can be overridden individually
    * We can identify mandatory properties whose values will not be copied to be reused but has to be overridden at 
    installation time 
  * Original values fom all `Secret`s are hidden and must be overridden at installation time
* The generated installer is agnostic from the original namespaces

## Feature design
* Based on a Golang application that runs the `Konveyor crane` CLI to export and normalize the original configuration
* Post-execution manipulations are performed
  * Clear the reference to the original namespace
  * Create a sample `template` for the overlayed configuration where the developer can apply customizations
    * All the keys in any given `ConfigMap` can be overridden in properties files called `CONFIGMAP_NAME.env`
      * The skeleton of these files are automatically generated with all the key names commented out 
    * All the keys in the managed `Secret`s are customizable in the same way
      * The `template` overlay re-creates the `Secret`s from the template files called `SECRET_NAME.env`
      * The base `kustomize` configuration does not re-create the `Secret`s, so its deployment would actually fail
      * Errors must be raised while trying to install the default template for the secrets

**TODO**:
* [Handle `Secret` securely ](https://github.com/zvigrinberg/handle-secrets-with-kustomize/blob/main/README.md)
* Manage mandatory params
  * Remove from copied ConfigMaps
  * Put in custom.env as__DEFAULT__
* Handle properties that are not managed as ConfigMap/Secret keys
* Export of cluster-wide resources
* Filter out automatically created resources (e.g., some RoleBindings)
* Management of OpenShift resources (e.g. Route)
* Consider cross-namespace references (e.g. a Service URL like "<svc name>.<ns-name>")
* Skip `kubernetes.io/service-account-token` Secrets
* Manage image registries

## Builder configuration
The `install-builder` runs using a configuration that specifies the packaging behavior: 
```yaml
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
* `oc`, `crane` and `move2kube` CLI are installed
* Login `oc` to the source OpenShift cluster
* `go` at least version 1.19 

Run this command to create the installer from the given configuration `myapp.yaml`:
```bash
go run main.go myapp.yaml
```

The command creates an `output/<APP NAME>/installer` folder in the current directory with the whole installer package.

### Output specification
The `installer/kustomize` folder contains the base configuration and an overlay template for each of the configured 
namespaces:
```bash
├── output
│   └── "APPLICATION"
│       └── installer
│           └── kustomize
│               └── "NAMESPACE1"
│                   ├── base
│                   │   ├── "RESOURCE1.yaml"
│                   │   ├── "RESOURCE2.yaml"
│                   │   └── kustomization.yaml
│                   └── template
│                       ├── kustomization.yaml
│                       ├── params
│                       │   ├── "CONFIG_MAP1.env"
│                       │   └── "CONFIG_MAP2.env"
│                       └── secrets
│                           └── "SECRET1.env"
```

The `base` kustomization contains the resources extracted from the source cluster, stripped of the status information,
cluster specific settings and namespace configuration.

The `template` overlay is an example of a possible `kustomize` overlay, with skeletons of environment files to override the
`ConfigMap`s parameters (remove the `#` comment and set the desired value) and all the `Secret` values:
```bash
> cat output/Infinity/installer/kustomize/holdings/template/params/CONFIG_MAP1.env
#KEY1=__EMPTY
#KEY2=__EMPTY

> cat output/Infinity/installer/kustomize/holdings/template/secrets/SECRET1.env
KEY1=__EMPTY
```

## Customize and install the template
Simple procedure that will be automated using the [replica-installer](../replica-installer/README.md) tool.

```bash
cd output/APPLICATION/installer/kustomize/NAMESPACE
cp -r template MYCONFIG
cd MYCONFIG
# Apply namespaxce update (must be created first)
kustomize edit set namespace MYNAMESPACE
# Edit changes to params/*.env and secrets/*env
kustomize build . | oc apply -f-
```
## Open points
* Which permissions are needed to export
* Parametrize overlay names:
    * [code](https://github.com/konveyor/move2kube/blob/3d57835d897596bed2bd42d937b6c5f2ac173f73/transformer/kubernetes/parameterizer/parameterizer.go#L57)
    * [Parameterizing custom fields in Helm Chart, Kustomize, OC Templates](https://move2kube.konveyor.io/tutorials/customizing-the-output/custom-parameterization-of-helm-charts-kustomize-octemplates)
* `move2kube`: add externalizer script to automate the extraction
    * What for Helm and OpenShift template?
* `crane`: Create transformer to automatically remove namespaces
    * [customplugins](https://konveyor.github.io/crane/tools/customplugins/)
* What if there are Jobs needed to run before installing the app? (e.g., dbinit)
* Convert/adapt cluster versions (e.g. adapt to different K8s API versions)