#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

kubectl delete pod crtest
kubectl apply -f crtest.yaml

if [ $? -ne 0 ]; then
  echo "create pod crtest failed"


kubectl delete podcheckpoints podcheckpoint-sample
kubectl migrate -f /data/go/src/k8s.io/kubernetes/vendor/k8s.io/podcheckpoint/crds/podcheckpoint.yaml --node=worker01

if [ $? -ne 0 ]; then
  echo "migrate pod crtest failed"