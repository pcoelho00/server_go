package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type MetricsData struct {
	Visits int
}

func (cfg *ApiConfig) MetricsHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/metrics/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := MetricsData{
		Visits: cfg.FileserverHits,
	}
	tmpl.Execute(w, data)
}

func (cfg *ApiConfig) ResetStatsHandler(w http.ResponseWriter, r *http.Request) {
	cfg.FileserverHits = 0
	msg := fmt.Sprintf("Hits reset to %v", cfg.FileserverHits)
	log.Println(msg)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(msg))
}

func (cfg *ApiConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.FileserverHits++
		log.Printf("middleware %v", cfg.FileserverHits)
		next.ServeHTTP(w, r)
	})
}
