package main

import (
	"fmt"
	"net/http"

	"github.com/Divyue30597/web_server/internal/auth"
)

func (cfg *apiConfig) revoke(w http.ResponseWriter, r *http.Request) {
	refresh_token, err := auth.GetToken(r)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("%v", err))
		return
	}

	err = cfg.DB.UpdateUserFromRefreshToken(refresh_token)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error revoking token")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
