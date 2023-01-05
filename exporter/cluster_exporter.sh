#!/bin/bash
# This script requires to be logged in a K8s/Openshift cluster and the CLI configured

# Check cluster API connection
[[ $? -ne 0 ]] && { echo "Not connected to a cluster. Exiting..."; exit; }


## Global vars
################################################################################
PREFIX="cloned-"
RESULTS_DIR="./results"
# Extracting API address
ORI_CLUSTER_API=$(kubectl cluster-info | \
	# removing color characters
  sed -e 's/\x1b\[[0-9;]*m//g' | \
	# API URL
  sed -r 's/.*is running at (.*)/\1/' | \
  head -1
)
ORI_CLUSTER_ID="$(oc get clusterversion -o jsonpath='{.items[].spec.clusterID}{"\n"}')"


## Init
echo "Initializating..."
[[ -d $RESULTS_DIR ]] || { mkdir -p $RESULTS_DIR; }


## Main
################################################################################
echo "Copying cluster with ID: $ORI_CLUSTER_ID at $ORI_CLUSTER_API"

echo "Getting Cluster Info..."
export CLUSTER_NAME=$(echo "$PREFIX$ORI_CLUSTER_API" | sed -r 's/https:\/\/api.([0-9a-zA-Z\-]*)\.(.*):6443/\1/')
export CLUSTER_BASE_DOMAIN=$(echo "$ORI_CLUSTER_API" | sed -r 's/https:\/\/api.([0-9a-zA-Z\-]*)\.(.*):6443/\2/')
export CLUSTER_VERSION=$(oc get clusterversion -o go-template='{{range .items}}{{.spec.desiredUpdate.version}}{{"\n"}}{{end}}')


echo "Getting Infrastructure Info..."
export WORKER_COUNT=$(oc get nodes --selector=node-role.kubernetes.io/worker --no-headers | wc -l)
export CLUSTER_NETWORK=$(oc get network.config/cluster -o go-template='{{range .spec.clusterNetwork}}{{.cidr}}{{"\n"}}{{end}}')
export HOST_PREFIX=$(oc get network.config/cluster -o go-template='{{range .spec.clusterNetwork}}{{.hostPrefix}}{{"\n"}}{{end}}')
export SERVICE_NETWORK=$(oc get network.config/cluster -o go-template='{{range .spec.serviceNetwork}}{{.}}{{"\n"}}{{end}}')
export NETWORK_TYPE=$(oc get network.config/cluster -o go-template='{{.spec.networkType}}')


echo "Getting Registry Info..."
export REGISTRY_ROUTE_HOSTNAME=$(oc get routes image-registry -n openshift-image-registry -o go-template='{{.spec.host}}' | sed -r "s/(.*).apps..*/\1.$CLUSTER_NAME.$CLUSTER_BASE_DOMAIN/g")
if [[ "$(echo $REGISTRY_ROUTE_INFO | wc -l)" -ne 0 ]]; then
  export REGISTRY_IS_EXPOSED="true"
else
  export REGISTRY_IS_EXPOSED="false"
fi


echo "Getting Cloud/Bare-Metal Provider Info..."
export PROV_CLOUD_PROVIDER=$(oc get Infrastructure cluster -o go-template='{{.status.platform}}')
export PROV_CLOUD_REGION=$(oc get Infrastructure cluster -o go-template="{{.status.platformStatus.$(echo $PROV_CLOUD_PROVIDER | tr '[:upper:]' '[:lower:]').region}}")


## Generating Template
################################################################################
echo "Generating Manifest..."
j2 templates/manifest.j2 > $RESULTS_DIR/cluster_manifest_${CLUSTER_NAME}_${CLUSTER_BASE_DOMAIN}.yaml


echo "Done"
