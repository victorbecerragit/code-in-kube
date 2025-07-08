# Excercises from Programming with Kubernetes (educative.io)

# Sample code to create a simple http server (webhook server) for authorization services

You have deployed the webhook in lab-05, setup the certificate on lab-04 ready for kubernetes services and now you are ready to create the container for the webhook.

docker build -t <github_account>/authz-webhook:v0.1 ./

# docker build with github actions

https://github.com/victorbecerragit/action-docker/tree/main/authz-webhook


# Login to the control-plane node and follow these steps.

```

To Test the webhook in kubernetes ,first we need to add the file webhook-config.yml to the control-plane node path /etc/kubernetes/pki/

And apply the deploy.


kubectl apply -f ./authz-webhook-deploy.yml

In order to let our static kube-apiserver Pod access the webhook service with the DNS name my-authz.kube-system.svc, 
we need to export the DNS record inside the control-plane node (/etc/hosts)

echo "$(kubectl get svc -n kube-system my-authz | grep -v NAME | awk '{print $3}') my-authz.kube-system.svc" >> /etc/hosts

Now, we need to tell the kube-apiserver how to use our authorization webhook. This can be configured with the flag --authorization-webhook-config-file=<my-config-file-path> of the kube-apiserver. 
Then, append the webhook mode to the flag --authorization-mode. In this lesson, our kube-apiserver is running as a static Pod, whose lifecycle is managed by the kubelet. We only need to update the kube-apiserver manifest file, and then the kubelet will be notified of the change and restart it later.

Enable webhook authorization in the kube-apiserver:

# Enable webhook mode in kube-apiserver
sed -i "s|Node,RBAC|Node,RBAC,Webhook|" /etc/kubernetes/manifests/kube-apiserver.yaml

# Include authorization webhook config file to kube-apiserver
sed -i "18i\    - --authorization-webhook-config-file=/etc/kubernetes/pki/webhook-config.yml" /etc/kubernetes/manifests/kube-apiserver.yaml

# Remove the current kube-apiserver docker to recreate it.
docker ps -a | grep apiserver | grep -v "pause" | awk '{print $1}' | xargs docker rm -f

Let’s do some tests to see if the kube-apiserver could successfully use the token webhook authorization service at https://my-authz.kube-system.svc/authorize. We use kubectl --as to impersonate users for operations.

 Try to get ns using the user "abc"

kubectl get ns --as abc

The output should be similar to this:
Error from server (Forbidden): namespaces is forbidden: User "abc" cannot list resource "namespaces" in API group "" at the cluster scope: User "abc" is not allowed to access "namespaces"

Now, let's try to access the namespaces for the demo-user:

kubectl get ns --as demo-user

The output should be similar to this:
NAME              STATUS   AGE
default           Active   10m
kube-node-lease   Active   10m
kube-public       Active   10m
kube-system       Active   10m

When we impersonate a random user, such as abc, to list namespaces, the request is forbidden. This is exactly what we expected. On the other hand, when we use our mock user demo-user, listing all the namespaces is allowed.

```