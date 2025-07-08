
# Sample code to create a simple http server (webhook server) that authenticate a user (mock sample)

You have deployed the webhook in lab-01, setup the certificate on lab-02 ready for kubernetes services and now you are ready to create the container for the webhook.

docker build -t <github_account>/authn-webhook:v0.1 ./

# docker build with github actions

https://github.com/victorbecerragit/action-docker/tree/main/authn-webhook


# Login to the control-plane node and follow these steps.

```

To Test the webhook in kubernetes ,first we need to add the file webhook-config.yml to the control-plane node path /etc/kubernetes/pki/

And apply the deploy.


kubectl apply -f ./authn-webhook-deploy.yml

In order to let our static kube-apiserver Pod access the webhook service with the DNS name my-authn.kube-system.svc, 
we need to export the DNS record inside the control-plane node (/etc/hosts)

echo "$(kubectl get svc -n kube-system my-authn | grep -v NAME | awk '{print $3}') my-authn.kube-system.svc" >> /etc/hosts

Now, we need to tell the kube-apiserver how to use our authentication webhook. This can be configured using the flag --authentication-token-webhook-config-file of the kube-apiserver.

Our kube-apiserver is running as a static Pod, whose lifecycle is managed by the kubelet. 
We only need to update the kube-apiserver manifest file. Then, the kubelet will be notified of the change and restart it later. 

Enable webhook authentication in the kube-apiserver:

#sed -i "18i\    - --authentication-token-webhook-config-file=/etc/kubernetes/pki/webhook-config.yml" /etc/kubernetes/manifests/kube-apiserver.yaml

sed -i "18i\    - --authorization-webhook-config-file=/etc/kubernetes/pki/webhook-config.yml" /etc/kubernetes/manifests/kube-apiserver.yaml
docker ps -a | grep apiserver | grep -v "pause" | awk '{print $1}' | xargs docker rm -f

Let’s do some tests to see if the kube-apiserver can successfully use the token webhook authentication service at https://my-authn.kube-system.svc/authorize. 
We can use the token created in lab-01 "ZGVtby11c2VyOmRlbW8tcGFzc3dvcmQ=" (tokenreview.json) to construct a kubeconfig file:


cp ~/.kube/config ./
kubectl config --kubeconfig=./config set-credentials webhook-test --token="ZGVtby11c2VyOmRlbW8tcGFzc3dvcmQ="
kubectl config --kubeconfig=./config set-context kubernetes-admin@kubernetes --user=webhook-test

With this kubeconfig file, we can run the command below in the terminal above to see what happened:

kubectl get ns --kubeconfig=./config

The output should be similar to the following:
Error from server (Forbidden): namespaces is forbidden: User "mock" cannot list resource "namespaces" in API group "" at the cluster scope

We can see that the token is identified correctly by our webhook. Additionally, this output just shows that the webhook is working as expected, because the webhook treats the token as the user mock, who apparently has no right (we didn’t assign any RBAC roles for this mock user) to list namespaces.

```