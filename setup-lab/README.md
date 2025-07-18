# Setup kubernetes lab


Create empty kuberentes env, with custom volumes, usefull to have a "clean" environment.

./kind - > bash -x ./bootstrap.sh

To delete the cluster :

$kind delete cluster -A


# Single control-plane cluster to allow perform labs of webhooks and admissions.

./kube-docker > deploy a single VM or kubernetes in docker all-in-one image. 
 - https://kwok.sigs.k8s.io/docs/user/all-in-one-image/


# Lab
VM - Ubuntu 20.04

kubeadm,kubectl - 1.23.8

kubeadm init --config=/configfile.yml
mkdir -p $HOME/.kube
cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
chown $(id -u):$(id -g) $HOME/.kube/config
kubectl apply -f /kube-flannel.yml




