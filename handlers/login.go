package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/pcoelho00/server_go/jsondecoders"
)

func (cfg *ApiConfig) PostLoginHandler(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	user := PostUser{}
	err := decoder.Decode(&user)

	if err != nil {
		jsondecoders.RespondWithError(w, http.StatusBadRequest, "Bad Arguments")
	}

	User, err := cfg.DB.GetUserFromLogin(user.Email, user.Password)
	if err != nil {
		jsondecoders.RespondWithError(w, http.StatusUnauthorized, "User not Found/Wrong Password")
	} else {
		jsondecoders.RespondWithJson(w, http.StatusOK, User)
	}

}
