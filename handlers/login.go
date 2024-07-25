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
)

type LoggedUser struct {
	Id           int    `json:"id"`
	Email        string `json:"email"`
	JwtToken     string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	IsChirpyRed  bool   `json:"is_chirpy_red"`
}

func (cfg *ApiConfig) PostLoginHandler(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	var postUser PostUser

	err := decoder.Decode(&postUser)

	if err != nil {
		jsondecoders.RespondWithError(w, http.StatusBadRequest, "Bad Arguments")
		return
	}

	User, err := cfg.DB.GetUserFromLogin(postUser.Email, postUser.Password)
	if err != nil {
		jsondecoders.RespondWithError(w, http.StatusUnauthorized, "User not Found/Wrong Password")
		return
	}

	refresh_token, err := auth.CreateRefreshToken(32)
	if err != nil {
		jsondecoders.RespondWithError(w, http.StatusInternalServerError, "Error creating token")
		return
	}

	err = cfg.DB.SaveRefreshToken(refresh_token, User.Id)
	if err != nil {
		jsondecoders.RespondWithError(w, http.StatusInternalServerError, "Error saving token")
		return
	}

	token, err := auth.CreateToken(cfg.JwtSecret, postUser.ExpireSecs, User.Id)
	if err != nil {
		jsondecoders.RespondWithError(w, http.StatusInternalServerError, "Error creating token")
		return
	}

	jsondecoders.RespondWithJson(w, http.StatusOK, LoggedUser{
		Id:           User.Id,
		Email:        User.Email,
		JwtToken:     token,
		RefreshToken: refresh_token,
		IsChirpyRed:  User.IsChirpyRed,
	})

}

func (cfg *ApiConfig) PutLoginUserHandler(w http.ResponseWriter, r *http.Request) {

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

	id, err := strconv.Atoi(subject)
	if err != nil {
		log.Println(err.Error())
		jsondecoders.RespondWithError(w, http.StatusInternalServerError, "Couldn't Retrieve user")
		return
	}

	decoder := json.NewDecoder(r.Body)
	var postUser PostUser

	err = decoder.Decode(&postUser)

	if err != nil {
		jsondecoders.RespondWithError(w, http.StatusBadRequest, "Bad Arguments")
		return
	}

	log.Println(postUser)

	UpdatedUser, err := cfg.DB.UpdateUser(id, postUser.Email, postUser.Password)

	if err != nil {
		jsondecoders.RespondWithError(w, http.StatusInternalServerError, "Failed to Update User")
		return
	}

	jsondecoders.RespondWithJson(w, http.StatusOK, database.PublicUser{
		Id:          UpdatedUser.Id,
		Email:       UpdatedUser.Email,
		IsChirpyRed: UpdatedUser.IsChirpyRed,
	})

}

func (cfg *ApiConfig) RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {

	token_string := r.Header.Get("Authorization")
	token_string = strings.Replace(token_string, "Bearer ", "", 1)

	log.Println(token_string)

	id, err := cfg.DB.FindRefreshToken(token_string)
	if err != nil {
		log.Println(err.Error())
		jsondecoders.RespondWithError(w, http.StatusInternalServerError, "Error accessing the Database")
		return
	}

	if id == 0 {
		jsondecoders.RespondWithError(w, http.StatusUnauthorized, "Unauthorized User")
		return
	} else {
		type TokenResponse struct {
			Token string `json:"token"`
		}
		token, err := auth.CreateToken(cfg.JwtSecret, 60*60, id)
		if err != nil {
			log.Println(err.Error())
			jsondecoders.RespondWithError(w, http.StatusInternalServerError, "Error Creating Token")
			return
		}

		jsondecoders.RespondWithJson(w, http.StatusOK, TokenResponse{
			Token: token,
		})
		return
	}

}

func (cfg *ApiConfig) RevokeRefreshTokenHandler(w http.ResponseWriter, r *http.Request) {

	token_string := r.Header.Get("Authorization")
	token_string = strings.Replace(token_string, "Bearer ", "", 1)

	log.Println(token_string)

	err := cfg.DB.RevokeRefreshToken(token_string)
	if err != nil {
		log.Println(err.Error())
		jsondecoders.RespondWithError(w, http.StatusInternalServerError, "Error accessing the Database")
		return
	}

	jsondecoders.RespondWithNoBody(w, http.StatusNoContent)

}
