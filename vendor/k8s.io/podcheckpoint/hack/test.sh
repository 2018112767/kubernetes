kubectl delete podcheckpoints podcheckpoint-sample
kubectl migrate -f /data/go/src/k8s.io/kubernetes/vendor/k8s.io/podcheckpoint/crds/podcheckpoint.yaml --node=migrate=dst