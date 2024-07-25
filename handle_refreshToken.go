package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Divyue30597/web_server/internal/auth"
)

func (cfg *apiConfig) refreshToken(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	authHeader, err := auth.GetToken(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("%v", err))
		return
	}

	user, err := cfg.DB.GetUserFromRefreshToken(authHeader)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "error getting user")
		return
	}

	// create new token
	token, err := auth.CreateToken(user.Id, time.Hour, cfg.Jwt)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error creating token")
		return
	}

	if r.ContentLength > 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		Token: token,
	})
}
