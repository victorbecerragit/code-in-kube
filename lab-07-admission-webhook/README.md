# Excercises from Programming with Kubernetes (educative.io)

# Sample code to create an validating-admission-webhook.

```

A validating admission webhook service is a web server, because the kube-apiserver invokes it through HTTPS POST requests. 
Now, let’s implement such a service step by step.

```go

import:
    admissionapi "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"


```

Let’s write a simple HTTP server at the path /validate on port 443. It checks the Pod name and rejects all Pods having mock-app in their names. 

In a real-world scenario, we only need to replace the mock code in the function validate() (line 21 in main.go) with our actual business logic. In the demo above, we’re handling Pods, and other resources can be handled too, only if we’ve set matching rules in ValidatingAdmissionConfiguration

func validate(admissionReq *admissionapi.AdmissionRequest) (admissionResp *admissionapi.AdmissionResponse)


One of the notable qualities of good admission controllers is their capability for easy diagnostics. So, when rejecting a request, we should return an error with precise or detailed information about the denial reasons or suggested remediations, like above. We’re rejecting all the Pods whose names contain mock-app and returning a hint message.

The rest of the code (lines 68–117 in main.go) is a simple HTTP handler that handles requests from the kube-apiserver and deserializes the payloads, which are expected to be an AdmissionReview object in JSON format. This HTTP handler then sends the validating result back.

Another notable thing is that, after we finish handling AdmissionReview, we should not forget to set the Response.UID, which matches the Request.UID.

# Generate Certificates.

In production environments, the kube-apiserver uses the CA certificate to visit our admission webhooks, which are serving securely with HTTPS. Now, let’s use cfssl to create some certificates for safe serving.

```

```shell

echo '{"CA":{"expiry": "87600h","pathlen":0},"CN":"CA","key":{"algo":"rsa","size":2048}}' | cfssl gencert -initca - | cfssljson -bare ca -

echo '{"signing":{"default":{"expiry":"87600h","usages":["signing","key encipherment","server auth","client auth"]}}}' > ca-config.json

export ADDRESS=localhost,validating-admission-demo.kube-system.svc

export NAME=server

echo '{"CN":"'$NAME'","hosts":[""],"key":{"algo":"rsa","size":2048}}' | cfssl gencert -config=ca-config.json -ca=ca.pem -ca-key=ca-key.pem -hostname="$ADDRESS" - | cfssljson -bare $NAME

``` 

```
Here, we generate a self-signed certificate with the CN field set to validating-admission-demo.kube-system.svc. We’re going to run this service in Kubernetes, which is the endpoint that we want to expose. The DNS name validating-admission-demo.kube-system.svc indicates that there’s a Service named validating-admission-demo running in the namespace kube-system.


Build and deploy 

docker build -t dixudx/pwk:validating-admission-webhook ./

kubectl apply -f ./webhook-deploy.yaml

Then, we declare ValidatingWebhookConfiguration in the kube-apiserver.

CABUNDLE=`base64 -w 0 ca.pem` envsubst < /usercode/webhook-config.yml | kubectl apply -f -

The env var CABUNDLE is defined in webhook-config.yml, envsubst will overwrite the variable with the certificate decoded.

In order to let our static kube-apiserver Pod access the webhook service with the DNS name validating-admission-demo.kube-system.svc, we need to export the DNS record as follows. Run it in the terminal above.

echo "$(kubectl get svc -n kube-system validating-admission-demo | grep -v NAME | awk '{print $3}') validating-admission-demo.kube-system.svc" >> /etc/hosts

In this lesson, our kube-apiserver is running as a static Pod, whose lifecycle is managed by the kubelet. We only need to update the kube-apiserver manifest file, and then the kubelet will be notified of the change and restart it later.

docker ps -a | grep apiserver | grep -v "pause" | awk '{print $1}' | xargs docker rm -f


Test:

kubectl run a-mock-app --image=nginx:1.23

The command above will create a Pod with the name a-mock-app-pod, which contains mock-app. Let’s run it in the terminal above to see what happens

Error from server: admission webhook "validating-admission-demo.kube-system.svc" denied the request: Keep calm and this is a webhook demo in the cluster!

We can see that the Pod a-mock-app is rejected correctly by our webhook, which is working as expected.

```