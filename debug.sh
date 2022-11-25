systemctl stop kubelet
dlv --check-go-version=false --listen=:2345 --headless=true --api-version=2 --accept-multiclient exec /usr/bin/kubelet -- --bootstrap-kubeconfig=/etc/kubernetes/bootstrap-kubelet.conf \
--kubeconfig=/etc/kubernetes/kubelet.conf --cgroup-driver=systemd \
--config=/var/lib/kubelet/config.yaml --cgroup-driver=systemd --network-plugin=cni \
--pod-infra-container-image=registry.aliyuncs.com/google_containers/pause:3.2
