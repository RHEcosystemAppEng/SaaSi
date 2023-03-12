# deployer
A Golang CLI tool to deploy configurations extracted from a live OpenShit/Kubernetes environment from the [exporter](../exporter/README.md) tool.

## Dependencies
* Use [kustomize](konveyor.io/tools/crane/) to export the original configuration and remove cluster specific settings
  (e.g. IP addresses, status, ...)
* Need [Ansible](https://docs.ansible.com/ansible/latest/installation_guide/intro_installation.html) installed in order to run playbook to Provision an OCP Cluster on AWS.
* Need [Jinja](https://pypi.org/project/Jinja2/) template engine Binary in order to render template that will become yaml parameters file of the ansible playbook.
* Need [htpasswd](https://command-not-found.com/htpasswd) util in order to define identity provider in OCP cluster
* need [Openshift Cli Tool](https://console.redhat.com/openshift/downloads)
* Need to download a [pull secret](https://console.redhat.com/openshift/install/pull-secret) for the Openshift installer, and copy paste it into [manifest.j2](./infra/templates/manifest.j2) pullSecret key
```yaml
    {%- if PROV_CLOUD_PROVIDER != 'None' %}
    platform:
            {{ PROV_CLOUD_PROVIDER | lower }}:
              region: "{{ PROV_CLOUD_REGION }}"
        {%- endif %}

    sshKey: ''
    pullSecret: 'Paste Here between quotes'
```
* The deployer dependent on Ansible playbook to provision the cluster, so before running the deployer with the option to provision a new cluster enabled, run the following command from root directory of repo:

```shell
git submodule add git@github.com:RHEcosystemAppEng/ocp_labs_provisioner.git deployer/infra/playbook
```
**Note: If .gitmodules file already exists, need to run the following command:**
```shell
git submodule update --checkout  -- deployer/infra/playbook 
```
Note: If you're not sure if you're at the top level of the git repo or not, just run the following two commands:
```shell
pwd
git rev-parse --show-toplevel
```
the two commands will bring the same directory only if you're in root directory of repo.

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
## Deployer Configuration with Cluster Provisioning
```yaml
deployer:
  cluster:
    # Must be unique
    clusterId: SAASI_DEMO_CLUSTER_ID
    # These will be eventually put in a secured store (after MVP1)
    # For MVP1 just store the AWS settings in a fixed vault/secret
    aws:
      aws_public_domain: fsi-partner.rhecoeng.com
      aws_account_name: AWS_ACCOUNT_NAME
      aws_access_key_id: CHANGEME
      aws_secret_access_key: CHANGEME
    params:
      # All are optional - but needed to be overridden if values weren't exported 
      # from source cluster or if extracted values from 
      # the source cluster are not relevant for your AWS Account. 
      # Anyway supplying CLUSTER_NAME is always recommended and sometimes mandatory.
      CLUSTER_NAME: saasi-cluster
      CLUSTER_BASE_DOMAIN: fsi-partner.rhecoeng.com
      WORKER_COUNT: 3
      MASTER_COUNT: 1
      CLUSTER_VERSION: 4.10.0
      PROV_CLOUD_REGION: eu-west-3
      REGISTRY_ROUTE_HOSTNAME: openshift-registry.fsi-partner.rhecoeng.com
  application:
    name: infinity
    namespaceMappingFormat: "%s-prod"
    namespaces:
    - name: campaign
      # If missing, the namespaceMappingFormat or the original name are applied
      target: campaign1
      params:
      - configMap: openshift-service-ca.crt
        params:
        - name: service-ca.crt
          value: VALUE-1
      - configMap: campaignms-api-config
        params:
        - name: DB_CONNECTION_URL
          value: VALUE-1
      secrets:
      - secret: campaignmsdbsecret
        params:
        - name: DB_PASS_ENCRYPTION_KEY
          value: VALUE-1
        - name: PASSWORD
          value: VALUE-1
        - name: USER_NAME
          value: VALUE-1
    - name: arrangement
      # If missing, the namespaceMappingFormat or the original name are applied
      target: arrangement1
      params:
      - configMap: openshift-service-ca.crt
        params:
        - name: service-ca.crt
          value: VALUE-1
      secrets:
      - secret: arrangementdbsecret
        params:
        - name: DB_PASS_ENCRYPTION_KEY
          value: VALUE-1
        - name: PASSWORD
          value: VALUE-1
        - name: USER_NAME
          value: VALUE-1
    - name: adapterms
      # If missing, the namespaceMappingFormat or the original name are applied
      secrets:
      - secret: adaptermsdbsecret
        params:
        - name: DB_PASS_ENCRYPTION_KEY
          value: VALUE-1
        - name: PASSWORD
          value: VALUE-1
        - name: USER_NAME
          value: VALUE-1
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
        Application configuration file for deployment
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
* Need to propagate logs of openshift installer to console output of deployer application , to bring an enhanced user experience.
* Need to add an optional lifecycle hook, that will be placed at the point of time that is between provisioning the cluster and deploying the application , that its implementation will be provided by user, this will give the user an option to deploy dependencies in a form of script location, that will be passed as command line argument - This script will take care of installing all the application' dependencies that weren't migrated from source cluster( for example, kafka cluster, databases, etc).
* Need to add cluster config parameter to the input yaml file, of the pull secret value that will be injected into [manifest.yaml](./infra/templates/manifest.j2)  


