package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/pcoelho00/server_go/jsondecoders"
)

type ResponseData struct {
	UserId int `json:"user_id"`
}

type WebHookBody struct {
	Event string       `json:"event"`
	Data  ResponseData `json:"data"`
}

func (cfg *ApiConfig) ChirpyRedHandler(w http.ResponseWriter, r *http.Request) {

	token_string := r.Header.Get("Authorization")
	token_string = strings.Replace(token_string, "ApiKey ", "", 1)

	if token_string != cfg.PolkaKey {
		jsondecoders.RespondWithNoBody(w, http.StatusUnauthorized)
		return
	}

	decoder := json.NewDecoder(r.Body)
	msg := WebHookBody{}

	err := decoder.Decode(&msg)
	if err != nil {
		log.Println(err.Error())
		jsondecoders.RespondWithNoBody(w, http.StatusInternalServerError)
		return
	}

	if msg.Event != "user.upgraded" {
		jsondecoders.RespondWithNoBody(w, http.StatusNoContent)
		return
	}

	updated, err := cfg.DB.UpdateToChirpyRed(msg.Data.UserId)
	if err != nil {
		log.Println(err.Error())
		jsondecoders.RespondWithNoBody(w, http.StatusInternalServerError)
		return
	}

	if !updated {
		jsondecoders.RespondWithNoBody(w, http.StatusNotFound)
		return
	} else {
		jsondecoders.RespondWithNoBody(w, http.StatusNoContent)
		return
	}

}
