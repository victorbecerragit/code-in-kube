# Excercises from Programming with Kubernetes (educative.io)

# Sample code to create a simple http server (webhook server) that authenticate a user (mock sample)

You have deployed the webhook in lab-01.

In production environments, using HTTPS is more secure. Now, let’s create some certificates for safe serving:

openssl req -x509 -newkey rsa:2048 -nodes -subj "/CN=my-authn.kube-system.svc" -keyout key.pem -out cert.pem

Here, we generate a self-signed certificate with the CN field set to my-authn.kube-system.svc.
We’re going to run this service in Kubernetes, which is the endpoint that we want to expose. 
This DNS name my-authn.kube-system.svc indicates that there’s a Service named my-authn running in the namespace kube-system.

In the main.go, just was added flag module, certificate and key files, which are called by the http listiner server.

var (
		certFile string
		keyFile  string
	)
	flag.StringVar(&certFile, "tls-cert-file", "", "File containing the default x509 Certificate for HTTPS.")
	flag.StringVar(&keyFile, "tls-private-key-file", "", "File containing the default x509 private key matching --tls-cert-file.")
	flag.Parse()

 	// Listen to port 443 and wait
	log.Println("Listening on port 443 for requests...")
	log.Fatal(http.ListenAndServeTLS(":443", certFile, keyFile, nil))




