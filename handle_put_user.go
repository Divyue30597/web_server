package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Divyue30597/web_server/internal/token"
)

func (cfg *apiConfig) putUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds int64  `json:"expires_in_seconds"`
	}

	decode := json.NewDecoder(r.Body)

	params := parameters{}

	err := decode.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error decoding body")
		return
	}

	tkn, err := getToken(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("%v", err))
		return
	}

	// validate token and get the user id
	someTkn, err := token.VerifyToken(tkn, cfg.Jwt)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("%v", err))
		return
	}

	fmt.Println(someTkn.Claims.GetSubject())

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

	user, err := cfg.DB.UpdateUser(userId, params.Email, params.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "error updating user")
		return
	}

	respondWithJSON(w, http.StatusOK, user)
}

func getToken(r *http.Request) (string, error) {
	signedToken := r.Header.Get("Authorization")

	if signedToken == "" {
		return "", fmt.Errorf("no token provided")
	}

	newToken := strings.Split(signedToken, " ")

	if newToken[0] != "Bearer" || len(newToken) != 2 {
		return "", fmt.Errorf("invalid token")
	}

	return newToken[1], nil
}
