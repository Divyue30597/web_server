package main

import (
	"encoding/json"
	"net/http"

	"github.com/Divyue30597/web_server/internal/database"
	"golang.org/x/crypto/bcrypt"
)

func (cfg *apiConfig) postUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := database.User{}

	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding JSON body")
		return
	}

	params.Password = hashPassword(params.Password)
	data, err := cfg.DB.CreateUser(params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error creating user")
		return
	}

	respondWithJSON(w, http.StatusCreated, struct {
		Id    int    `json:"id"`
		Email string `json:"email"`
	}{
		Id:    data.Id,
		Email: data.Email,
	})
}

func hashPassword(password string) string {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err.Error()
	}

	return string(hashedPass)
}
