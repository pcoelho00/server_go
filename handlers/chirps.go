package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/pcoelho00/server_go/database"
	"github.com/pcoelho00/server_go/jsondecoders"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type ApiConfig struct {
	FileserverHits int
	DB             *database.DB
}

type ValidMsg struct {
	Valid string `json:"cleaned_body"`
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

func (cfg *ApiConfig) PostJsonHandler(w http.ResponseWriter, r *http.Request) {

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
		chirp_msg, err := cfg.DB.CreateChirp(new_msg)
		if err != nil {
			jsondecoders.RespondWithError(w, http.StatusBadRequest, "Couldn't write to database")
		}
		jsondecoders.RespondWithJson(w, http.StatusCreated, chirp_msg)
		dbStructure, err := cfg.DB.WriteChirpsToDB(chirp_msg)
		if err != nil {
			jsondecoders.RespondWithError(w, http.StatusBadRequest, "Couldn't write to database")
		}
		cfg.DB.WriteDB(dbStructure)

	} else {
		jsondecoders.RespondWithError(w, http.StatusBadRequest, "Chirp is too long")
	}

}

func (cfg *ApiConfig) GetChirpsMsgHandler(w http.ResponseWriter, r *http.Request) {

	ChirpsMsgs, err := cfg.DB.GetChirps()
	if err != nil {
		jsondecoders.RespondWithError(w, http.StatusBadRequest, "Couldn't return msg from the database")
	}

	jsondecoders.RespondWithJson(w, http.StatusOK, ChirpsMsgs)
}

func (cfg *ApiConfig) GetChirpHandler(w http.ResponseWriter, r *http.Request) {

	ChirpId, err := strconv.Atoi(r.PathValue("chirpID"))
	if err != nil {
		jsondecoders.RespondWithError(w, http.StatusBadRequest, "Error getting ID")
	}
	ChirpsMsg, err := cfg.DB.GetChirp(ChirpId)
	if err != nil {
		jsondecoders.RespondWithError(w, http.StatusNotFound, "Chirp ID doesn't exists")
	} else {
		jsondecoders.RespondWithJson(w, http.StatusOK, ChirpsMsg)
	}

}
