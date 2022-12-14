# replica-installer
A Golang CLI tool to extract configurations from a live OpenShit/Kubernetes environment and generate a reusable, configurable
installer for the [replica-installer](../replica-installer/README.md) tool.

## Dependencies
* Use [kustomize](konveyor.io/tools/crane/) to export the original configuration and remove cluster specific settings
  (e.g. IP addresses, status, ...)

## Installer configuration
The `replica-installer` runs using a configuration that specifies the installation behavior:
```yaml
application:
  # This creates an installer package named APP
  name: APP
  # Optional, defines the automatic mapping for the target namespaces
  # %s will be replaced with the original namespace name
  namespaceMappingFormat: "%s-prod"
  namespaces:
    - name: NS1
      # If missing, the namespaceMappingFormat or the original name are applied 
      target: NEWNS1
      params:
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

## Running the installer
Prerequisites:
* `oc`, `kustomize` CLI are installed
* Login `oc` to the target OpenShift cluster
* `go` at least version 1.19

Run this command to install the package extracted in `./installer/MYAPP` from the given configuration `myapp.yaml`:
```bash
go run main.go myapp.yaml ./installer/MYAPP
```

**TODO**:
* Manage mandatory params
  * If __DEFAULT__ stop installation

## Open points
* Which user permissions are needed to install

