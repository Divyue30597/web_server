package main

import (
	"encoding/json"
	"net/http"
)

func (cfg *apiConfig) postUser(w http.ResponseWriter, r *http.Request) {
	type postUser struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := postUser{}

	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding JSON body")
		return
	}

	data, err := cfg.DB.CreateUser(params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating user")
		return
	}

	respondWithJSON(w, http.StatusCreated, data)
}
