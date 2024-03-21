package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	admission_v1 "k8s.io/api/admission/v1"
)

var (
	optListen  = flag.String("listen", ":8443", "Listen expression")
	optTlsCert = flag.String("cert", "/etc/tls/tls.crt", "Path to TLS certificate")
	optTlsKey  = flag.String("key", "/etc/tls/tls.key", "Path to TLS key")
)

func handler(w http.ResponseWriter, r *http.Request) {
	var review admission_v1.AdmissionReview

	err := json.NewDecoder(r.Body).Decode(&review)
	if err != nil {
		log.Println("Invalid request received:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println("Request body: ", review)

	http.Error(w, "Error", 500)
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
