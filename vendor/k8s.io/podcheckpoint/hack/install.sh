mv /usr/bin/kubelet /usr/bin/kubelet.back
mv /usr/bin/kubectl /usr/bin/kubectl.back
mv _output/dockerized/bin/linux/amd64/kubelet /usr/bin/kubelet
mv _output/dockerized/bin/linux/amd64/kubectl /usr/bin/kubectl
systemctl restart kubelet