package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Divyue30597/web_server/internal/auth"
)

func (cfg *apiConfig) deleteChirps(w http.ResponseWriter, r *http.Request) {
	authHeader, err := auth.GetToken(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("%v", err))
		return
	}

	token, err := auth.VerifyToken(authHeader, cfg.Jwt)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("%v", err))
		return
	}

	id, err := token.Claims.GetSubject()
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "error getting user id")
		return
	}

	authorId, err := strconv.Atoi(id)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "error converting string to int")
		return
	}

	chirpID := r.PathValue("chirpID")
	chirpId, err := strconv.Atoi(chirpID)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "error converting string to int")
		return
	}

	chirp, err := cfg.DB.GetSingleChirp(chirpId)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("%v", err))
		return
	}

	if chirp.AuthorID == authorId {
		err = cfg.DB.DeleteChirp(chirpId)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("%v", err))
			return
		}
		w.WriteHeader(http.StatusNoContent)
	} else {
		w.WriteHeader(http.StatusForbidden)
		return
	}
}
