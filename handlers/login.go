package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pcoelho00/server_go/auth"
	"github.com/pcoelho00/server_go/jsondecoders"
)

type LoginUser struct {
	Id       int    `json:"id"`
	Email    string `json:"email"`
	JwtToken string `json:"token"`
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

	log.Println(User)

	if err != nil {
		jsondecoders.RespondWithError(w, http.StatusUnauthorized, "User not Found/Wrong Password")
		return
	}
	token, err := auth.CreateToken(cfg.JwtSecret, postUser.ExpireSecs, User.Id)
	if err != nil {
		jsondecoders.RespondWithError(w, http.StatusInternalServerError, "Error creating token")
		return
	}

	jsondecoders.RespondWithJson(w, http.StatusOK, LoginUser{
		Id:       User.Id,
		Email:    User.Email,
		JwtToken: token,
	})

}

func (cfg *ApiConfig) PutLoginUserHandler(w http.ResponseWriter, r *http.Request) {

	token_string := r.Header.Get("Authorization")
	token_string = strings.Replace(token_string, "Bearer ", "", 1)

	log.Println(token_string)

	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(token_string, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.JwtSecret), nil
	})

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

	UpdatedUser, err := cfg.DB.UpdateUser(id, postUser.Email, postUser.Password)

	if err != nil {
		jsondecoders.RespondWithError(w, http.StatusInternalServerError, "Failed to Update User")
		return
	}

	jsondecoders.RespondWithJson(w, http.StatusOK, ResponseUser{
		Id:    UpdatedUser.Id,
		Email: UpdatedUser.Email,
	})

}
