package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type ChirpsMsg struct {
	Body string
}

type ValidMsg struct {
	Valid string `json:"cleaned_body"`
}

func (cfg *apiConfig) JsonHandler(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	msg := ChirpsMsg{}

	err := decoder.Decode(&msg)

	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(msg.Body) <= 140 {
		new_msg := profaneWords(msg.Body)
		respondWithJson(w, http.StatusOK, ValidMsg{new_msg})
	} else {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
	}

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
