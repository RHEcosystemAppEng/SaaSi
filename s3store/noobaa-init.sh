#!/bin/bash

TARGET_NS=$1
export AWS_ACCESS_KEY_ID=$(kubectl get secret noobaa-admin -n $TARGET_NS -o json | jq -r '.data.AWS_ACCESS_KEY_ID|@base64d')
export AWS_SECRET_ACCESS_KEY=$(kubectl get secret noobaa-admin -n $TARGET_NS -o json | jq -r '.data.AWS_SECRET_ACCESS_KEY|@base64d')
export S3_ENDPOINT=$(oc get route s3 -n $TARGET_NS -ojsonpath='{.spec.host}')
export S3_REGION=$(oc get noobaa noobaa -n $TARGET_NS -ojsonpath='{.spec.region}')
if test -z "$S3_REGION" 
then
      export S3_REGION="us-east-1"
fi

alias s3='aws --endpoint "https://$S3_ENDPOINT" --no-verify-ssl s3'
