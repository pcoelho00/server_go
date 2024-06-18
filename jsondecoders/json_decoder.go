package jsondecoders

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func RespondWithError(w http.ResponseWriter, code int, msg string) {
	type ErrorResponse struct {
		Error string `json:"error"`
	}

	RespondWithJson(w, code, ErrorResponse{msg})
}

func RespondWithJson(w http.ResponseWriter, code int, payload interface{}) {
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

func ProfaneWords(msg string) string {
	title := cases.Title(language.English)
	upper := cases.Upper(language.English)

	for _, word := range [3]string{"kerfuffle", "sharbert", "fornax"} {
		r := strings.NewReplacer(word, "****", title.String(word), "****", upper.String(word), "****")
		msg = r.Replace(msg)
	}
	return msg
}
