package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/pcoelho00/server_go/database"
)

type apiConfig struct {
	fileserverHits int
	db             *database.DB
}

func main() {
	const port = "8080"
	const root = "."
	const templates = root + "/templates"

	db, err := database.NewDB("database.json")
	if err != nil {
		log.Fatal("Can't connect with the Database")
	}

	mux := http.NewServeMux()
	apiCfg := apiConfig{
		fileserverHits: 0,
		db:             db,
	}

	handler := http.StripPrefix("/app", http.FileServer(http.Dir(templates)))
	mux.Handle("/app/*", apiCfg.middlewareMetricsInc(handler))
	mux.HandleFunc("GET /api/healthz", HealthsResponseHandler)

	mux.HandleFunc("GET /api/reset", apiCfg.ResetStatsHandler)
	mux.HandleFunc("GET /admin/metrics", apiCfg.MetricsHandler)
	mux.HandleFunc("POST /api/chirps", apiCfg.PostJsonHandler)
	mux.HandleFunc("GET /api/chirps", apiCfg.GetJsonHandler)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	fmt.Printf("Server started at port: %s\n", port)
	log.Fatal(server.ListenAndServe())

}
