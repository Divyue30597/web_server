package main

import "net/http"

func (cfg *apiConfig) getChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.DB.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error getting chirps")
	}

	respondWithJSON(w, http.StatusOK, chirps)
}
