package main

import (
	"net/http"
	"sort"
	"strconv"

	"github.com/Divyue30597/web_server/internal/database"
)

func (cfg *apiConfig) getChirps(w http.ResponseWriter, r *http.Request) {
	s := r.URL.Query().Get("author_id")
	sort := r.URL.Query().Get("sort")

	var chirps []database.Chirp
	var err error

	if s != "" {
		authorID, err := strconv.Atoi(s)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "error converting string to int")
			return
		}
		chirps, err = cfg.DB.GetChirpsByAuthorID(authorID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "error getting chirps")
			return
		}
	} else {
		chirps, err = cfg.DB.GetChirps()
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "error getting chirps")
			return
		}
	}

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error getting chirps")
		return
	}

	sortedChirps := sortChirps(chirps, sort)
	respondWithJSON(w, http.StatusOK, sortedChirps)
}

func sortChirps(chirps []database.Chirp, sortStr string) []database.Chirp {
	sort.Slice(chirps, func(i, j int) bool {
		if sortStr == "desc" {
			return chirps[i].Id > chirps[j].Id
		}
		return chirps[i].Id < chirps[j].Id
	})

	return chirps
}
