# Excercises from Programming with Kubernetes (educative.io)

# Sample code to create an mutating-admission-webhook.

```
A mutating admission webhook service is a web server, because the kube-apiserver invokes it through HTTPS POST requests. 

import
	admissionapi "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"


Let’s write a simple HTTP server at the path /mutate over port 443. It adds labels "my-labe-mock=test" for Pods.

    // Inject labels
		if podCopy.Labels == nil {
			podCopy.Labels = map[string]string{}
		}
		podCopy.Labels["my-label-mock"] = "test"

...

// Now we create a mock as an example
	// TODO: Please replace this with your logic
	switch admissionReq.Operation {
	case admissionapi.Create, admissionapi.Update:
		// Create a copy, so that we can modify it
		podCopy := pod.DeepCopy()

		// Inject labels
		if podCopy.Labels == nil {
			podCopy.Labels = map[string]string{}
		}
		podCopy.Labels["my-label-mock"] = "test"

		var patchList jsonpatch.JSONPatchList
		patchList, err = jsonpatch.CreateJSONPatch(podCopy, pod)
		if err != nil {
			return
		}
		admissionResp.Patch = patchList.Raw()
		admissionResp.PatchType = new(admissionapi.PatchType)
		*admissionResp.PatchType = admissionapi.PatchTypeJSONPatch
	}

	admissionResp.Allowed = true
	return


In a real-world scenario, we only need to replace the mock code in the function admit() (line 21 in main.go) with our actual business logic. In the demo above, we’re handling pods, and other resources can be handled as well if we’ve set matching rules in MutatingAdmissionConfiguration.

func validate(admissionReq *admissionapi.AdmissionRequest) (admissionResp *admissionapi.AdmissionResponse)

One of the notable qualities of good admission controllers is their capability for easy diagnostics. So, when rejecting a request, we should return an error with precise or detailed information about the denial reasons or suggested solutions, as shown in the code widget above. We reject the requests if we failed to patch the pods and return a hint message.

The rest of the code (lines 78–127 in main.go) is a simple HTTP handler that handles requests from the kube-apiserver and deserializes the payloads, which are expected to be an AdmissionReview object in JSON format, and then sends the mutating result back.

Another notable thing is that, after we finish handling the AdmissionReview, we should not forget to set Response.UID, which matches Request.UID.

# Generate Certificates.

In production environments, the kube-apiserver uses the CA certificate to visit our admission webhooks, which are serving securely with HTTPS. Now, let’s use cfssl to create some certificates for safe serving.

```bash 

echo '{"CA":{"expiry": "87600h","pathlen":0},"CN":"CA","key":{"algo":"rsa","size":2048}}' | cfssl gencert -initca - | cfssljson -bare ca -

echo '{"signing":{"default":{"expiry":"87600h","usages":["signing","key encipherment","server auth","client auth"]}}}' > ca-config.json

export ADDRESS=localhost,mutating-admission-demo.kube-system.svc

export NAME=server

echo '{"CN":"'$NAME'","hosts":[""],"key":{"algo":"rsa","size":2048}}' | cfssl gencert -config=ca-config.json -ca=ca.pem -ca-key=ca-key.pem -hostname="$ADDRESS" - | cfssljson -bare $NAME

```

 we generate a self-signed certificate with the CN field set to mutating-admission-demo.kube-system.svc. 
 We’re going to run this service in Kubernetes, which is the endpoint that we want to expose. 
 The DNS name mutating-admission-demo.kube-system.svc indicates that there’s a Service named mutating-admission-demo running in the namespace kube-system.


Build and deploy 

docker build -t dixudx/pwk:mutating-admission-webhook ./

kubectl apply -f ./webhook-deploy.yaml

We declare MutatingWebhookConfiguration in the kube-apiserver.

CABUNDLE=`base64 -w 0 ca.pem` envsubst < /usercode/webhook-config.yml | kubectl apply -f -

The env var CABUNDLE is defined in webhook-config.yaml, envsubst will overwrite the variable with the certificate decoded.

In order to let our static kube-apiserver Pod access the webhook service with the DNS name mutating-admission-demo.kube-system.svc, we need to export the DNS record as follows. 

echo "$(kubectl get svc -n kube-system mutating-admission-demo | grep -v NAME | awk '{print $3}') mutating-admission-demo.kube-system.svc" >> /etc/hosts

In this lesson, our kube-apiserver is running as a static Pod, whose lifecycle is managed by the kubelet. We only need to update the kube-apiserver manifest file, and then the kubelet will be notified of the change and restart it later.

docker ps -a | grep apiserver | grep -v "pause" | awk '{print $1}' | xargs docker rm -f


Test:

```bash 

kubectl run a-mock-app --image=nginx:1.23

kubectl get pod a-mock-app -o yaml | grep -C 5 labels

The output will be as follows:
apiVersion: v1
kind: Pod
metadata:
  creationTimestamp: "2022-08-17T13:32:20Z"
  labels:
    my-label-mock: test
    run: a-mock-app
  name: a-mock-app
  namespace: default
  resourceVersion: "1600"

  We can see that a new label my-label-mock=test has already been injected to the Pod a-mock-app by our webhook, which is working as expected.