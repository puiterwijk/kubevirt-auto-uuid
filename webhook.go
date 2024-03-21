package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

var (
	optListen  = flag.String("listen", ":8443", "Listen expression")
	optTlsCert = flag.String("cert", "/etc/tls/tls.crt", "Path to TLS certificate")
	optTlsKey  = flag.String("key", "/etc/tls/tls.key", "Path to TLS key")
)

func handler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request: ", r)

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
