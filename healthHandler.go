package main

import (
	"net/http"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	switch r.Method {
	case http.MethodGet:
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(http.StatusText(http.StatusOK)))
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
