# deployer
A Golang CLI tool to deploy configurations extracted from a live OpenShit/Kubernetes environment from the [exporter](../exporter/README.md) tool.

## Dependencies
* Use [kustomize](konveyor.io/tools/crane/) to export the original configuration and remove cluster specific settings
  (e.g. IP addresses, status, ...)

## Installer configuration
The `deployer` runs using a configuration to specify how to install the packaged configuration:
* For each namespace we can define a target namespace name, reuse the original one or apply a common strategy to transform
the initial name into the new one
* For all the mandatory and optional parameters defined in `ConfigMaps`, we can define the desired value
* For all the parameters defined as 1`Secret` keys, we must provide the desired value

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
    target: NEW-NS1
    params:
     - configMap: MAP-1
        params:
        - name: PARAM-1
          value: VALUE-1
        - name: "..."
          value: "..."
     - configMap: MAP-2
         params:
        - name: PARAM-1
          value: VALUE-1
       - name: "..."
         value: "..."
    secrets:
      - secret: SECRET-1
          params:
          - name: PARAM-1
            value: VALUE-1
          - name: "..."
            value: "..."
      - secret: SECRET-2
          params:
          - name: PARAM-1
            value: VALUE-1
          - name: "..."
            value: "..."
```

## Running the installer
Prerequisites:
* Install `oc` and `kustomize` CLI tools
* Install `go` CLI tool (at least version 1.19)
* `oc` login to the source OpenShift cluster

Run this command to install the package extracted in `./exporter` from the given configuration `myapp.yaml`:
```bash
go run main.go -f myapp.yaml -e ./exporter
```

Available options:
```bash
> go run main.go --help
  -f string
        Application configuration file
  -e string
        Root folder of exporter application (shorthand)
  -exporter-dir string
        Root folder of exporter application
  -o string
        Root output folder (shorthand) (default "<CURRENT_DIR>/output")
  -output-dir string
        Root output folder (default "<CURRENT_DIR>/output")
```

```bash
├── output
│   └── clusters
│       └── CLUSTER_ID1
│           └── TBD: metadata to specify source cluster ID, target cluster ID and target customizations
...
│   └── applications
│       └── APPLICATION1
│           └── UUID1: one such folder for each different installation of APPLICATION1
│               ├── kustomize
│               │   ├── NAMESPACE1
│               │   │   ├── base
│               │   │   │   └──   ...
│               │   │   └── template
│               │   │       └──   ...
│               │   ├── NAMESPACE2
│               │   │   └── ...
│               │   └── ...
│               └── deploy
│                   ├── NAMESPACE1.yaml
│                   ├── NAMESPACE2.yaml
│                   └── ...

```

**TODO**:
* Manage mandatory params
  * If __DEFAULT__ stop installation

## Open points
* Which user permissions are needed to install

