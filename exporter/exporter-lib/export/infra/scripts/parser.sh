#!/bin/sh
# This script reads the cluster config ENV vars and the customizations ENV vars to generate the manifest.yaml file to create the cluster infrastructure


## Global vars
################################################################################
PREFIX="cloned_"
RESULTS_DIR="./results"

# Print script usage flags
function print_usage() {
  echo "
  $0
    -s: Specifies the source cluster env vars file.
        (Example values in: ./exporter/customize.env)
    -c: Specifies the file with the customization env vars.
        (Example values in: ./exporter/customize.env)
    -h: Prints script's usage
  "
}

## Main
################################################################################
# Flags
while getopts 's:c:h' flag; do
  case "${flag}" in
    s) SOURCE_VARS_FILE="$OPTARG" ;;
    c) USER_DEF_VARS_FILE="$OPTARG" ;;
    h) print_usage; exit 0 ;;
    *) echo "Error: invalid flag. Exiting"; print_usage; exit 1 ;;
  esac
done

# Check source file var
if [ -z $SOURCE_VARS_FILE ]; then
  echo "Error: missing flags. Exiting"
  print_usage
  exit 1
fi


## Reading Source cluster Vars
################################################################################
if [ -f $SOURCE_VARS_FILE ]; then
  echo "Source file: $SOURCE_VARS_FILE"
  source $SOURCE_VARS_FILE
fi
TEMPLATE_FILE="$RESULTS_DIR/customized_${CLUSTER_NAME}_${CLUSTER_BASE_DOMAIN}.yaml"

## Reading User defined Vars
################################################################################
if [[ ! -z $USER_DEF_VARS_FILE && -f $USER_DEF_VARS_FILE ]]; then
  echo "Custom vars file: $USER_DEF_VARS_FILE"
  source $USER_DEF_VARS_FILE
fi


## Generating Template
################################################################################
echo "Generating Manifest at: $TEMPLATE_FILE"
j2 templates/manifest.j2 > $TEMPLATE_FILE
