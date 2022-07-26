#!/usr/bin/env bash
export KIND_EXPERIMENTAL_PROVIDER=podman

NAME="${1:-test-1}"
CLUSTER_NAME="kind-${NAME}"

ind create cluster --image kindest/node:v1.24.0 --wait 5m --name test-1 --config ./cluster.yaml

kubectl --context "${CLUSTER_NAME}" run -i --tty busybox --image=busybox --restart=Never -- sh
