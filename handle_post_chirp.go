package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func (cfg *apiConfig) postChirp(w http.ResponseWriter, r *http.Request) {
	type postchirp struct {
		Body string `json:"body"`
	}

	const validChirpLen = 140

	decoder := json.NewDecoder(r.Body)

	params := postchirp{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding JSON body")
		return
	}

	if len(params.Body) > validChirpLen {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	data, err := cfg.DB.CreateChirp(cleanChirp(params.Body))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating chirp")
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
