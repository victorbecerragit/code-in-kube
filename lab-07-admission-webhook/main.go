package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	admissionapi "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	podResource = metav1.GroupVersionResource{Version: "v1", Resource: "pods"}
)

func validate(admissionReq *admissionapi.AdmissionRequest) (admissionResp *admissionapi.AdmissionResponse) {
	admissionResp = &admissionapi.AdmissionResponse{}
	// copy uid from tr.Request
	admissionResp.UID = admissionReq.UID

	var err error
	defer func() {
		// If the handler returned an error, incorporate the error message into the response
		if err != nil {
			admissionResp.Allowed = false
			admissionResp.Result = &metav1.Status{
				Message: err.Error(),
			}
		}
	}()

	if admissionReq.Resource != podResource {
		log.Printf("expect resource to be %s, but got %s", podResource, admissionReq.Resource)
		err = fmt.Errorf("expect resource to be %s", podResource)
		return
	}

	// this demo webhook cares about pods
	// Parse the Pod object.
	// TODO: Please replace this with your logic
	pod := corev1.Pod{}
	if err = json.Unmarshal(admissionReq.Object.Raw, &pod); err != nil {
		return
	}

	// now we do some mock for demo
	// TODO: Please replace this with your logic
	switch admissionReq.Operation {
	case admissionapi.Create, admissionapi.Update:
		if strings.Contains(pod.Name, "mock-app") {
			admissionResp.Allowed = false
			admissionResp.Result = &metav1.Status{
				Message: "Keep calm and this is a webhook demo in the cluster!",
			}
			return
		}
	}

	admissionResp.Allowed = true
	return
}

func sampleHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Receiving %s", r.Method)

	if r.Method != "POST" {
		http.Error(w, "Only Accept POST requests", http.StatusMethodNotAllowed)
		return
	}

	// Read body of POST request
	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Unmarshal JSON from POST request to AdmissionReview object
	admissionReview := admissionapi.AdmissionReview{}
	err = json.Unmarshal(payload, &admissionReview)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	admissionReview.Response = validate(admissionReview.Request)

	// Marshal the AdmissionReview to JSON and send it back
	result, err := json.Marshal(admissionReview)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(result)
	w.Header().Set("Content-Type", "application/json")
}

func main() {
	var (
		certFile string
		keyFile  string
	)
	flag.StringVar(&certFile, "tls-cert-file", "", "File containing the default x509 Certificate for HTTPS.")
	flag.StringVar(&keyFile, "tls-private-key-file", "", "File containing the default x509 private key matching --tls-cert-file.")
	flag.Parse()

	// Set up a /validate resource handler
	http.HandleFunc("/validate", sampleHandler)

	// Listen to port 443 and wait
	log.Println("Listening on port 443 for requests...")
	log.Fatal(http.ListenAndServeTLS(":443", certFile, keyFile, nil))
}