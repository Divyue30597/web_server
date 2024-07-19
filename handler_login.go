package main

import (
	"encoding/json"
	"net/http"

	"github.com/Divyue30597/web_server/internal/database"
)

func (cfg *apiConfig) login(w http.ResponseWriter, r *http.Request) {
	decode := json.NewDecoder(r.Body)

	params := database.User{}

	err := decode.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error decoding body")
	}

	validUser, err := cfg.DB.VerifyUser(params)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	respondWithJSON(w, http.StatusOK, struct {
		Id    int    `json:"id"`
		Email string `json:"email"`
	}{
		Id:    validUser.Id,
		Email: validUser.Email,
	})
}
