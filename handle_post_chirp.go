package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Divyue30597/web_server/internal/auth"
)

func (cfg *apiConfig) postChirp(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetToken(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("%v", err))
		return
	}

	tkn, err := auth.VerifyToken(token, cfg.Jwt)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("%v", err))
		return
	}

	id, err := tkn.Claims.GetSubject()
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Error getting user id")
		return
	}

	authorId, err := strconv.Atoi(id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error converting string to int")
		return
	}

	type postchirp struct {
		Body     string `json:"body"`
		AuthorID int    `json:"author_id"`
	}

	const validChirpLen = 140

	decoder := json.NewDecoder(r.Body)

	params := postchirp{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error decoding JSON body")
		return
	}

	if len(params.Body) > validChirpLen {
		respondWithError(w, http.StatusBadRequest, "chirp is too long")
		return
	}

	data, err := cfg.DB.CreateChirp(cleanChirp(params.Body), authorId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error creating chirp")
		return
	}

	respondWithJSON(w, http.StatusCreated, data)
}

func cleanChirp(str string) string {
	strs := strings.Split(str, " ")

	for _, bad_word := range profanity {
		for i := 0; i < len(strs); i++ {
			if strings.ToLower(strs[i]) == bad_word {
				strs[i] = "****"
			}
		}
	}

	return strings.Join(strs, " ")
}
