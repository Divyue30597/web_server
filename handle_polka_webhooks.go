package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/Divyue30597/web_server/internal/database"
)

func (cfg *apiConfig) polkaWebhooks(w http.ResponseWriter, r *http.Request) {
	token, err := GetApiKeyAuthVal(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("%v", err))
		return
	}

	if token != cfg.ApiKey {
		respondWithError(w, http.StatusUnauthorized, "invalid token")
		return
	}

	type params struct {
		Event string `json:"event"`
		Data  struct {
			UserID int `json:"user_id"`
		}
	}

	decoder := json.NewDecoder(r.Body)
	parameters := params{}

	err = decoder.Decode(&parameters)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error decoding params")
		return
	}

	if parameters.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// user, err := cfg.DB.GetUserById(parameters.Data.UserID)
	// if err != nil {
	// 	respondWithError(w, http.StatusNotFound, "user not found")
	// 	return
	// }

	// _, err = cfg.DB.UpdateUser(user.Id, user.Email, user.Password, true)
	// if err != nil {
	// 	respondWithError(w, http.StatusNotFound, "error updating user info")
	// 	return
	// }

	_, err = cfg.DB.UpgradeChirpyRed(parameters.Data.UserID)
	if err != nil {
		if errors.Is(err, database.ErrNotExist) {
			respondWithError(w, http.StatusNotFound, "couldn't find user")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "couldn't update user")
		return
	}

	// when the user is updated in the database successfully
	w.WriteHeader(http.StatusNoContent)
}

func GetApiKeyAuthVal(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	tokens := strings.Split(authHeader, " ")

	if len(tokens) != 2 || tokens[0] != "ApiKey" {
		return "", fmt.Errorf("token not provided or invalid token")
	}

	if authHeader == "" {
		return "", fmt.Errorf("token not found")
	}

	return tokens[1], nil
}
