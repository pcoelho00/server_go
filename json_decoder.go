package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/pcoelho00/server_go/database"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type ValidMsg struct {
	Valid string `json:"cleaned_body"`
}

func (cfg *apiConfig) PostJsonHandler(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	msg := database.ChirpsMsg{}

	err := decoder.Decode(&msg)

	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(msg.Body) <= 140 {
		new_msg := profaneWords(msg.Body)
		chirp_msg, err := cfg.db.CreateChirp(new_msg)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Couldn't write to database")
		}
		respondWithJson(w, http.StatusCreated, chirp_msg)
		dbStructure, err := cfg.db.UpdateDB(chirp_msg)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Couldn't write to database")
		}
		cfg.db.WriteDB(dbStructure)

	} else {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
	}

}

func (cfg *apiConfig) GetJsonHandler(w http.ResponseWriter, r *http.Request) {

	ChirpsMsgs, err := cfg.db.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't return msg from the database")
	}

	respondWithJson(w, http.StatusOK, ChirpsMsgs)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type ErrorResponse struct {
		Error string `json:"error"`
	}

	respondWithJson(w, code, ErrorResponse{msg})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	u, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error encoding parameters: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(u)
}

func profaneWords(msg string) string {
	title := cases.Title(language.English)
	upper := cases.Upper(language.English)

	for _, word := range [3]string{"kerfuffle", "sharbert", "fornax"} {
		r := strings.NewReplacer(word, "****", title.String(word), "****", upper.String(word), "****")
		msg = r.Replace(msg)
	}
	return msg
}
