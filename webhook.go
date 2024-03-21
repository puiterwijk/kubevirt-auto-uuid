package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

var (
	optListen  = flag.String("listen", ":8080", "Listen expression")
	optTlsCert = flag.String("cert", "/etc/tls.crt", "Path to TLS certificate")
	optTlsKey  = flag.String("key", "/etc/tls.key", "Path to TLS key")
)

func handler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request: ", r)

	http.Error(w, "Error", 500)
}

func main() {
	flag.Parse()

	http.HandleFunc("/", handler)

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
