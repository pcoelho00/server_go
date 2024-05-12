package main

import (
	"fmt"
	"html/template"
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
	}

}

type MetricsData struct {
	Visits int
}

func (cfg *apiConfig) Metrics(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/metrics/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := MetricsData{
		Visits: cfg.fileserverHits,
	}
	tmpl.Execute(w, data)
}

func (cfg *apiConfig) ResetStats(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits = 0
	msg := fmt.Sprintf("Hits reset to %v", cfg.fileserverHits)
	log.Println(msg)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(msg))
}
