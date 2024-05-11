package main

import (
	"fmt"
	"log"
	"net/http"
	"path"
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
	if path.Base(r.URL.Path) == "metrics" {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		msg := fmt.Sprintf("Hits: %v", cfg.fileserverHits)
		w.Write([]byte(msg))
	} else if path.Base(r.URL.Path) == "reset" {
		cfg.fileserverHits = 0
		msg := fmt.Sprintf("Hits reset to %v", cfg.fileserverHits)
		log.Println(msg)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(msg))
	}

}
