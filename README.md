# SaaSi
One Stop shop for hybrid cloud native application packing and deployment tools 

## Exporter
The exporter is the component which collects the needed infrastructure details
to generate a YAML manifest that could be used to setup the lab by the
[`ocp-lab-provisioner`](https://github.com/RHEcosystemAppEng/ocp_labs_provisioner)
tool.

The exporter requires the `oc` CLI installed and configured to access the
Openshift cluster that you want to clone. To run it, please use the following
commands:
```sh
# Access the exporter's folder
cd exporter

# Check your oc cli connection
oc cluster-info

# Run the exporter
./cluster_exporter.sh

# Check results
ls ./results
```

**Note**: Each exported manifests could be found under `results` folder named
as: `cluster_manifest_cloned-<{TARGET_CLUSTER_DOMAIN>.yaml
