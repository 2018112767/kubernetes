# 1. install go
wget -c https://dl.google.com/go/go1.19.2.linux-amd64.tar.gz -O - |  tar -xz -C /usr/local

echo "export GOPATH=$HOME/go" >> /etc/profile
echo "export GOBIN=$GOPATH/bin" >> /etc/profile
echo "export PATH=$PATH:/usr/local/go/bin:$GOBIN" >> /etc/profile

source /etc/profile
go version

# 2. install docker

apt-get update
apt-get install \
    apt-transport-https \
    ca-certificates \
    curl \
    gnupg-agent \
    software-properties-common -y
curl -fsSL https://mirrors.aliyun.com/docker-ce/linux/ubuntu/gpg | sudo apt-key add -
sudo add-apt-repository \
     "deb [arch=amd64] https://mirrors.aliyun.com/docker-ce/linux/ubuntu \
     $(lsb_release -cs) \
     stable"
apt-get update
apt-get install docker-ce=5:19.03.5~3-0~ubuntu-bionic docker-ce-cli=5:19.03.5~3-0~ubuntu-bionic containerd.io -y

# 3. install gcc make git
apt install git make gcc -y

# 4. install kubeadm
vim /etc/ssh/sshd_config
PermitRootLogin yes
service ssh restart
ufw disable
vim /etc/fstab #注释swap
cat >> /etc/hosts << EOF
192.168.153.130 master
EOF

cat >> /etc/sysctl.d/k8s.conf << EOF
net.bridge.bridge-nf-call-ip6tables = 1
net.bridge.bridge-nf-call-iptables = 1
EOF
sysctl --system




apt-get update && apt-get install -y apt-transport-https -y

curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key add -

cat <<EOF >/etc/apt/sources.list.d/kubernetes.list
deb http://apt.kubernetes.io/ kubernetes-xenial main
EOF

apt-get update
apt-get install kubeadm=1.18.0-00 kubelet=1.18.0-00 kubectl=1.18.0-00
kubeadm init --pod-network-cidr=10.244.0.0/16 --ignore-preflight-errors=NumCPU --apiserver-advertise-address=192.168.153.130 --image-repository registry.aliyuncs.com/google_containers

mkdir -p $HOME/.kube
cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
chown $(id -u):$(id -g) $HOME/.kube/config
export KUBECONFIG=/etc/kubernetes/admin.conf



# 5 install docker code
mkdir -p $GOPATH/src/github.com/docker
cd $GOPATH/src/github.com/docker
git clone https://github.com/2018112767/moby.git
cd moby
git checkout -b v19.03.9 origin/tag-19.03.9
make install

cd ..
git clone https://github.com/2018112767/cli.git
cd cli
git checkout -b t19 origin/t19


mv /usr/bin/kubelet /usr/bin/kubelet.back
mv /usr/bin/kubectl /usr/bin/kubectl.back
mv _output/dockerized/bin/linux/amd64/kubelet /usr/bin/kubelet
mv _output/dockerized/bin/linux/amd64/kubectl /usr/bin/kubectl
systemctl restart kubelet
journalctl -xeu kubelet

kubectl get podcheckpoints
kubectl delete podcheckpoints podcheckpoint-sample
kubectl migrate -f vendor/k8s.io/podcheckpoint/crds/podcheckpoint.yaml --node=migrate=true
