package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	const port = "8080"
	const root = "."

	mux := http.NewServeMux()
	apiCfg := apiConfig{}

	handler := http.StripPrefix("/app", http.FileServer(http.Dir(root)))
	mux.Handle("/app/*", apiCfg.middlewareMetricsInc(handler))
	mux.Handle("/healthz", healthzReponse{})
	mux.Handle("/metrics", &apiCfg)
	mux.Handle("/reset", &apiCfg)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	fmt.Printf("Server started at port: %s\n", port)
	log.Fatal(server.ListenAndServe())

}
