package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	admission_v1 "k8s.io/api/admission/v1"
	kubevirt_v1 "kubevirt.io/api/core/v1"
)

var (
	optDebug   = flag.Bool("debug", false, "Enable debugging")
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

	var response admission_v1.AdmissionReview
	response.APIVersion = review.APIVersion
	response.Kind = review.Kind
	response.Response = new(admission_v1.AdmissionResponse)
	response.Response.UID = review.Request.UID
	response.Response.Allowed = true

	if vm.Spec.Template.Spec.Domain.Firmware == nil || vm.Spec.Template.Spec.Domain.Firmware.UUID == "" {
		newUuid := uuid.New()

		fmt.Printf(
			"VM created without UUID, identifier: %s.%s, new UUID: %s\n",
			vm.Namespace,
			vm.Name,
			newUuid,
		)

		patch := []byte(fmt.Sprintf(
			`[{"op": "add", "path": "/spec/template/spec/domain/firmware/uuid", "value": "%s"}]`,
			newUuid,
		))

		response.Response.Warnings = []string{
			fmt.Sprintf(
				"No UUID in request, assigned %s",
				newUuid,
			),
		}
		var patchtype = admission_v1.PatchTypeJSONPatch
		response.Response.PatchType = &patchtype
		response.Response.Patch = patch
	}

	if *optDebug {
		if data, err := json.Marshal(response); err == nil {
			log.Println("Marshalled data: ", string(data))
		}
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
