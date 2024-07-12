package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const port = "8080"

type apiConfig struct {
	fileserverhits int
}

type chirp struct {
	Body string `json:"body"`
}

type errorResponse struct {
	Err string `json:"error"`
}

type validChirp struct {
	Valid bool `json:"valid"`
}

func main() {
	mux := http.NewServeMux()

	cfg := &apiConfig{
		fileserverhits: 0,
	}

	mux.Handle("/", http.FileServer(http.Dir(".")))

	mux.Handle("/app/*", http.StripPrefix("/app",
		cfg.middlewareMetricsInc(http.FileServer(http.Dir("."))),
	))

	mux.HandleFunc("GET /api/healthz", healthHandler)
	mux.HandleFunc("GET /admin/metrics", cfg.handlerMetrics)
	mux.HandleFunc("GET /api/reset", cfg.handlerReset)
	mux.HandleFunc("POST /api/validate_chirp", postChirp)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	fmt.Println("Serving on port " + port)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func postChirp(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	params := chirp{}
	err := decoder.Decode(&params)
	if err != nil {

		errVal := errorResponse{
			Err: "Something went wrong",
		}
		data, err := json.Marshal(errVal)
		if err != nil {
			log.Printf("error marshalling json: %v", err)
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(data)
		return
	}

	if len(params.Body) > 140 {
		errVal := errorResponse{
			Err: "Chirp is too long",
		}

		data, err := json.Marshal(errVal)
		if err != nil {
			log.Printf("error marshalling json: %v", err)
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(data)
		return
	}

	validChirp := validChirp{
		Valid: true,
	}

	data, err := json.Marshal(validChirp)
	if err != nil {
		errVal := errorResponse{
			Err: "Something went wrong",
		}
		data, err := json.Marshal(errVal)
		if err != nil {
			log.Printf("error marshalling json: %v", err)
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(data)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// func respondWithJSON


func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`
	<html>

		<body>
				<h1>Welcome, Chirpy Admin</h1>
				<p>Chirpy has been visited %d times!</p>
		</body>

	</html>
	`, cfg.fileserverhits)))
}




func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverhits = 0
	w.WriteHeader(http.StatusOK)
	// w.Write([]byte("Hits reset to 0"))
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverhits++
		w.Header().Add("Cache-Control", "no-cache")
		next.ServeHTTP(w, r)
	})
}
