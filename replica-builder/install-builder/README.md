# install-builder
A Golang CLI tool to extract configurations from a live OpenShit/Kubernetes environment and generate a reusable, configurable
installer for the [replica-installer](../replica-installer/README.md) tool.

## Dependencies
* Use [Konveyor crane](konveyor.io/tools/crane/) to export the original configuration and remove cluster specific settings 
(e.g. IP addresses, status, ...)
* Use [Konveyor move2kube](https://move2kube.konveyor.io/) to generate the [kustomize](https://kustomize.io/) installer 
for the given environment 

## Features
* The generated installer can replicate the resources of the original namespaces
* The generated installer allows to override all the application parameters defined in the ConfigMap and Secrets (encrypted)
* The generated installer is agnosti from the original namespaces

## Feature design
* Based on a Golang application that runs the Konveyor CLI tools and performs post-execution manipulation
  * Externalize all the `ConfigMap` and `Secret` keys to allow the customizations of each single property or just use the 
  default values
    * `oc extract` command is used for the purpose
    * The base `kustomize` configuration re-creates the ConfigMaps and Secrets using files that are extracted from the
    extracted configurations (1 file per key)
    * The `kustomize` overlays instead use a merged approach and can override only the needed keys using a properties
    file `custom.env`
  * Clear the reference to the original namespace

**TODO**:
* Export and encrypt the `Secret` values
* Manage mandatory params
  * Remove from defaults
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
        configMaps:
          # Paramaters are given as a list of ConfigMap name and key values  
          - name: MAP-1
            param: PARAM-1
          - name: "..."
            param: "..."
          - name: MAP-M
            param: PARAM-N
        secrets:
          # Paramaters are given as a list of Secret name and key values  
          - name: SECRET-1
            param: PARAM-1
          - name: "..."
            param: "..."
          - name: SECRET-M
            param: PARAM-N
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