#!/usr/bin/bash -xe
SERVICE_CA_KEY=$(oc -n openshift-service-ca get secret/signing-key -o go-template='{{index .data "tls.crt"}}')
oc apply -f manifests/deploy.yaml
cat manifests/webhook.yaml | sed -e 's/SERVICECACRT/$SERVICE_CA_KEY/' | oc apply -f -
