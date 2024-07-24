package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Divyue30597/web_server/internal/auth"
	"github.com/Divyue30597/web_server/internal/token"
)

func (cfg *apiConfig) login(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds int64  `json:"expires_in_seconds"`
	}

	type response struct {
		User
	}

	decode := json.NewDecoder(r.Body)

	params := parameters{}

	err := decode.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error decoding body")
	}

	if params.ExpiresInSeconds == 0 || params.ExpiresInSeconds > 24*60*60 {
		params.ExpiresInSeconds = time.Now().Add(time.Duration(24*60*60) * time.Second).Unix()
	}

	user, err := cfg.DB.GetUserByEmail(params.Email)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = auth.VerifyPassword(params.Password, user.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid password")
		return
	}

	token, err := token.CreateToken(user.Id, params.ExpiresInSeconds, cfg.Jwt)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't create token")
		return
	}

	w.Header().Set("Authorization", "Bearer "+token)

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			Id:    user.Id,
			Email: user.Email,
			Token: token,
		},
	})
}
