#!/bin/bash
# This script requires to be logged in a K8s/Openshift cluster and the CLI configured

# Check cluster API connection
[[ $? -ne 0 ]] && { echo "Not connected to a cluster. Exiting..."; exit; }

# removing color characters
 # Extracting API address
cluster_control_pane=$(kubectl cluster-info | \
  sed -e 's/\x1b\[[0-9;]*m//g' | \
  sed -r 's/.*is running at (.*)/\1/' | \
  head -1
)

# Getting Node details
## instanceType could be null if there is no cloud-provider
nodes="$(
  kubectl get nodes -o json | \
    jq '
      .items[] |
        {
          nodeName: .metadata.name,
          role: .metadata.labels | to_entries | .[] | select(.key | test("node-role.kubernetes.io/.*")) | .key,
          resources: .status.capacity,
          instanceType: .metadata.labels."node.kubernetes.io/instance-type"
        }
    ' | jq -n '[inputs]')"

# Getting Storage Classes
storage_classes="$(
  kubectl get storageclasses -o json | \
    jq '
      .items[] |
        {
          name: .metadata.name,
          provisioner: .provisioner
        }
    ' | jq -n '[inputs]')"

# Printing JSON report
echo "{
  \"cluster\":
  {
    \"apiAddress\": \"$cluster_control_pane\"
  },
  \"nodes\": $nodes,
  \"storage_classes\": $storage_classes
}" | jq
