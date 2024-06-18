package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/pcoelho00/server_go/database"
	"github.com/pcoelho00/server_go/handlers"
)

func main() {
	const port = "8080"
	const root = "."
	const templates = root + "/templates"

	_, err := os.Stat("database.json")
	if !os.IsNotExist(err) {
		err := os.Remove("database.json")
		if err != nil {
			log.Fatal("Couldn't delete the database")
		}
	}

	db, err := database.NewDB("database.json")
	if err != nil {
		log.Fatal("Can't connect with the Database")
	}

	mux := http.NewServeMux()
	apiCfg := handlers.ApiConfig{
		FileserverHits: 0,
		DB:             db,
	}

	handler := http.StripPrefix("/app", http.FileServer(http.Dir(templates)))
	mux.Handle("/app/*", apiCfg.MiddlewareMetricsInc(handler))
	mux.HandleFunc("GET /api/healthz", handlers.HealthsResponseHandler)

	mux.HandleFunc("GET /api/reset", apiCfg.ResetStatsHandler)
	mux.HandleFunc("GET /admin/metrics", apiCfg.MetricsHandler)
	mux.HandleFunc("POST /api/chirps", apiCfg.PostJsonHandler)
	mux.HandleFunc("GET /api/chirps", apiCfg.GetChirpsMsgHandler)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.GetChirpHandler)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	fmt.Printf("Server started at port: %s\n", port)
	log.Fatal(server.ListenAndServe())

}
