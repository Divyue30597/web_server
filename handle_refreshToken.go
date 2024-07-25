package main

import (
	"fmt"
	"net/http"

	"github.com/Divyue30597/web_server/internal/auth"
)

func (cfg *apiConfig) refreshToken(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	authHeader, err := auth.GetToken(r)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("%v", err))
		return
	}

	user, err := cfg.DB.GetUserFromRefreshToken(authHeader)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error getting user")
		return
	}

	// create new token
	token, err := auth.CreateToken(user.Id, cfg.Jwt)
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
