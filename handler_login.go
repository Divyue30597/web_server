package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Divyue30597/web_server/internal/auth"
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
		return
	}

	if params.ExpiresInSeconds == 0 || params.ExpiresInSeconds > 24*60*60 {
		params.ExpiresInSeconds = time.Now().Add(time.Hour).Unix()
	}

	user, err := cfg.DB.GetUserByEmail(params.Email)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = auth.VerifyPassword(params.Password, user.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid password")
		return
	}

	token, err := auth.CreateToken(user.Id, time.Hour, cfg.Jwt)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "couldn't create token")
		return
	}

	refreshToken, err := auth.CreateRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "couldn't create refresh token")
		return
	}

	// save refresh token to DB with expiration time
	err = cfg.DB.SaveRefreshToken(user.Id, refreshToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("%v", err))
		return
	}

	w.Header().Set("Authorization", "Bearer "+token)

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			Id:           user.Id,
			Email:        user.Email,
			Token:        token,
			RefreshToken: refreshToken,
		},
	})
}
