# Excercises from Programming with Kubernetes (educative.io)

# Sample code to create a simple http server (webhook server) for authorization services

```

A webhook authorization service is a web server, because the kube-apiserver invokes it through HTTPS POST requests.

- authorizationapi "k8s.io/api/authorization/v1beta1"

  func authZ(sar *authorizationapi.SubjectAccessReview)

Let’s write a simple HTTP server that responds with the mock authenticated user mock when requested for the /authorize resource over port 443.

In a real-world scenario, we only need to replace the mock codes in the function authZ() with our actual implementations to determine user privileges. This function should set Status.Allowed to true if the request performed by the user is allowed and false for invalid requests. For the denied reason, it can be set to the field Status.Reason.

The rest of the code is a simple HTTP handler that handles requests from the kube-apiserver and deserializes the payloads, which are expected to be a SubjectAccessReview object in JSON format, and then sends the authorization result back.

"SubjectAccessReview" is a Kubernetes API object used to check whether a user or entity has permission to perform a specific action.

In the lesson context, it is the JSON payload received by the HTTP handler, which is deserialized to determine user privileges via the authZ() function.

- subjectaccessreview.json


For testing this service, we’ll simulate such a POST request from the kube-apiserver by making it manually from our local machine.


curl -k -X POST -d @subjectaccessreview.json http://localhost:443/authorize

```

