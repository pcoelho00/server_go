package main

import (
	"log"
	"net/http"
	"strings"
)

type apiConfig struct {
	fileserverHits int
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++
		log.Printf("middleware %v", cfg.fileserverHits)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.Contains(r.URL.Path, "/admin/metrics") {
		cfg.Metrics(w, r)
	} else if strings.Contains(r.URL.Path, "/api/reset") {
		cfg.ResetStats(w, r)
	} else if strings.Contains(r.URL.Path, "/api/validate_chirp") {
		cfg.JsonHandler(w, r)
	}

}
