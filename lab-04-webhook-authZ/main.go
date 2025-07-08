package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	authorizationapi "k8s.io/api/authorization/v1beta1"
)

func authZ(sar *authorizationapi.SubjectAccessReview) {
	// now we do some mock for demo
	// Please replace this with your logic
	if sar.Spec.User == "demo-user" {
		sar.Status.Allowed = true
	} else {
		sar.Status.Reason = fmt.Sprintf("User %q is not allowed to access %q",
			sar.Spec.User, sar.Spec.ResourceAttributes.Resource)
	}
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
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
	log.Printf("Receiving SubjectAccessReview: %s\n", string(payload))

	// Unmarshal JSON from POST request to SubjectAccessReview object
	sar := &authorizationapi.SubjectAccessReview{}
	err = json.Unmarshal(payload, sar)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	authZ(sar)

	// Marshal the SubjectAccessReview to JSON and send it back
	result, err := json.Marshal(*sar)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(result)
	w.Header().Set("Content-Type", "application/json")
}

func main() {
	// Set up a /authorize resource handler
	http.HandleFunc("/authorize", helloHandler)

	// Listen to port 443 and wait
	log.Println("Listening on port 443 for requests...")
	log.Fatal(http.ListenAndServe(":443", nil))
}