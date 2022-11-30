#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

PROJECT_ROOT="${GOPATH}/src/k8s.io/kubernetes/"
echo ${PROJECT_ROOT}

systemctl stop kubelet
mv /usr/bin/kubelet /usr/bin/kubelet.back
mv /usr/bin/kubectl /usr/bin/kubectl.back
cp ${PROJECT_ROOT}/_output/dockerized/bin/linux/amd64/kubelet /usr/bin/kubelet
cp ${PROJECT_ROOT}/_output/dockerized/bin/linux/amd64/kubectl /usr/bin/kubectl
systemctl restart kubelet