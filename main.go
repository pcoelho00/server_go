package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	const port = "8080"
	const root = "."
	const templates = root + "/templates"

	mux := http.NewServeMux()
	apiCfg := apiConfig{}

	handler := http.StripPrefix("/app", http.FileServer(http.Dir(templates)))
	mux.Handle("/app/*", apiCfg.middlewareMetricsInc(handler))
	mux.Handle("GET /api/healthz", healthzReponse{})
	mux.Handle("GET /api/reset", &apiCfg)
	mux.Handle("GET /admin/metrics", &apiCfg)
	mux.Handle("POST /api/validate_chirp", &apiCfg)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	fmt.Printf("Server started at port: %s\n", port)
	log.Fatal(server.ListenAndServe())

}
