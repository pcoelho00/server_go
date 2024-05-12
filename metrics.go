package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
)

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
