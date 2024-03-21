package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	admission_v1 "k8s.io/api/admission/v1"
	kubevirt_v1 "kubevirt.io/api/core/v1"
)

var (
	optListen  = flag.String("listen", ":8443", "Listen expression")
	optTlsCert = flag.String("cert", "/etc/tls/tls.crt", "Path to TLS certificate")
	optTlsKey  = flag.String("key", "/etc/tls/tls.key", "Path to TLS key")
)

func handler(w http.ResponseWriter, r *http.Request) {
	var review admission_v1.AdmissionReview

	if err := json.NewDecoder(r.Body).Decode(&review); err != nil {
		log.Println("Invalid request received:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var vm kubevirt_v1.VirtualMachine
	if err := json.Unmarshal(review.Request.Object.Raw, &vm); err != nil {
		log.Println("Invalid VM request received:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println("Request VM: ", vm)
	fmt.Println("Request UID: ", review.Request.UID)

	var response admission_v1.AdmissionReview
	response.Response = new(admission_v1.AdmissionResponse)
	response.Response.UID = review.Request.UID
	response.Response.Allowed = false

	if vm.Spec.Template.Spec.Domain.Firmware == nil || vm.Spec.Template.Spec.Domain.Firmware.UUID == "" {
		fmt.Println("VM created without UUID, patching in...")
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println("Error writing response:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	flag.Parse()

	http.HandleFunc("/", handler)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})

	fmt.Printf("Listening on %s, tls cert %s, tls key %s\n", *optListen, *optTlsCert, *optTlsKey)
	log.Fatal(
		http.ListenAndServeTLS(
			*optListen,
			*optTlsCert,
			*optTlsKey,
			nil,
		),
	)
}
