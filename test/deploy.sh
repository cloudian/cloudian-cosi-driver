#! /bin/bash

set -euxo pipefail
cd "$(dirname "$0")"

# Teardown old deployments if they exist, just ignore errors rather than being clever
k3d cluster delete cosi-driver || true

k3d cluster create --config ../utils/k3d/k3d-config.yaml

# Copy an up-to-date image to the cluster
pushd ..
    make image
    k3d image import -c cosi-driver cloudian-cosi-driver:v0.0.0
popd

# 3rd Party Resources - Lock to a SHA? Given releases are thin on the ground
kubectl apply -k github.com/kubernetes-sigs/container-object-storage-interface-api
kubectl apply -k github.com/kubernetes-sigs/container-object-storage-interface-controller
