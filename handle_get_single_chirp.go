package main

import (
	"net/http"
	"strconv"
)

func (cfg *apiConfig) getSingleChirp(w http.ResponseWriter, r *http.Request) {
	chirpId := r.PathValue("chirpID")

	chirpID, err := strconv.Atoi(chirpId)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Chirp ID")
	}

	data, err := cfg.DB.GetSingleChirp(chirpID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	respondWithJSON(w, http.StatusOK, data)
}
