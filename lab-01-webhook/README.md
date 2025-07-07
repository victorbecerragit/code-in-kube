# Excercises from Programming with Kubernetes (educative.io)

# Sample code to create a simple http server (webhook server) that authenticate a user (mock sample)


Let’s write a simple HTTP server that responds with the mock authenticated user mock, when requested for the /authenticate resource over port 443.

In a real-world scenario, we only need to replace the mock codes in the function authN with our actual implementations to query our user-management system.
This function should set a valid UserInfo if the provided token is legitimate. Moreover, Status.Authenticated should be set as well—true for legitimate users and false for invalid users. 
For the error message, it can be set to the field Status.Error.

The rest of the code is a simple HTTP handler that handles requests from the kube-apiserver and deserializes the payloads, which are expected to be a TokenReview object in JSON format, and then sends the authentication result back.


Now, let’s test the HTTP server above locally to see if the service works as expected.

Create a file called /root/tokenreview.json with the content below:


{
  "apiVersion": "authentication.k8s.io/v1beta1",
  "kind": "TokenReview",
  "spec": {
    "token": "ZGVtby11c2VyOmRlbW8tcGFzc3dvcmQ="
  }
}


echo "ZGVtby11c2VyOmRlbW8tcGFzc3dvcmQ=" | base64 -d
demo-user:demo-password


For testing this service, we’ll simulate a POST request from the kube-apiserver by making it manually from our local machine.


On a new terminal , run the program:
go run ./main.go

On a second terminal, run the POST to test the webhook:

curl -k -X POST -d @tokenreview.json http://localhost:4443/authenticate

