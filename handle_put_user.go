package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Divyue30597/web_server/internal/auth"
)

func (cfg *apiConfig) putUser(w http.ResponseWriter, r *http.Request) {
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

	tkn, err := auth.GetToken(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("%v", err))
		return
	}

	// validate token and get the user id
	someTkn, err := auth.VerifyToken(tkn, cfg.Jwt)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("%v", err))
		return
	}

	id, err := someTkn.Claims.GetSubject()
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "error getting user id")
		return
	}

	// convert to int
	userId, err := strconv.Atoi(id)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "error converting user id")
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "error hashing password")
		return
	}

	user, err := cfg.DB.UpdateUser(userId, params.Email, hashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "error updating user")
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			Id:          user.Id,
			Email:       user.Email,
			IsChirpyRed: user.IsChirpyRed,
		},
	})
}
