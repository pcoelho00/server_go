package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/pcoelho00/server_go/auth"
	"github.com/pcoelho00/server_go/database"
	"github.com/pcoelho00/server_go/jsondecoders"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type ApiConfig struct {
	FileserverHits int
	DB             *database.DB
	JwtSecret      string
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

func (cfg *ApiConfig) PostChirpsHandler(w http.ResponseWriter, r *http.Request) {

	token_string := r.Header.Get("Authorization")
	token_string = strings.Replace(token_string, "Bearer ", "", 1)

	log.Println(token_string)

	claims, err := auth.GetJWTClaims(token_string, cfg.JwtSecret)
	if err != nil {
		log.Println(err.Error())
		jsondecoders.RespondWithError(w, http.StatusUnauthorized, "Couldn't authenticate user")
		return
	}

	subject, err := claims.GetSubject()
	log.Printf("Subject: %s\n", subject)

	if err != nil {
		log.Println(err.Error())
		jsondecoders.RespondWithError(w, http.StatusInternalServerError, "Couldn't authenticate user")
		return
	}

	author_id, err := strconv.Atoi(subject)
	if err != nil {
		log.Println(err.Error())
		jsondecoders.RespondWithError(w, http.StatusInternalServerError, "Couldn't Retrieve user")
		return
	}

	decoder := json.NewDecoder(r.Body)
	msg := database.ChirpsMsg{}

	err = decoder.Decode(&msg)

	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(msg.Body) <= 140 {
		new_msg := profaneWords(msg.Body)
		chirp_msg, err := cfg.DB.CreateChirp(author_id, new_msg)
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

func (cfg *ApiConfig) DeleteChirpHandler(w http.ResponseWriter, r *http.Request) {
	token_string := r.Header.Get("Authorization")
	token_string = strings.Replace(token_string, "Bearer ", "", 1)

	log.Println(token_string)

	claims, err := auth.GetJWTClaims(token_string, cfg.JwtSecret)
	if err != nil {
		log.Println(err.Error())
		jsondecoders.RespondWithError(w, http.StatusUnauthorized, "Couldn't authenticate user")
		return
	}

	subject, err := claims.GetSubject()
	log.Printf("Subject: %s\n", subject)

	if err != nil {
		log.Println(err.Error())
		jsondecoders.RespondWithError(w, http.StatusInternalServerError, "Couldn't authenticate user")
		return
	}

	author_id, err := strconv.Atoi(subject)
	if err != nil {
		log.Println(err.Error())
		jsondecoders.RespondWithError(w, http.StatusInternalServerError, "Couldn't Retrieve user")
		return
	}

	ChirpId, err := strconv.Atoi(r.PathValue("chirpID"))
	if err != nil {
		jsondecoders.RespondWithError(w, http.StatusBadRequest, "Error getting ID")
	}

	err = cfg.DB.DeleteChirp(author_id, ChirpId)
	if err != nil {
		jsondecoders.RespondWithError(w, http.StatusForbidden, "Couldn't delete msg")
	} else {
		jsondecoders.RespondWithNoBody(w, http.StatusNoContent)
	}

}
