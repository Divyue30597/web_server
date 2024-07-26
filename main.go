package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Divyue30597/web_server/internal/database"
	"github.com/joho/godotenv"
)

const port = "8080"

var profanity = []string{"kerfuffle", "sharbert", "fornax"}

type apiConfig struct {
	fileserverhits int
	DB             *database.DB
	Jwt            string
}

type User struct {
	Id           int    `json:"id"`
	Email        string `json:"email"`
	Password     string `json:"-"`
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

func main() {
	godotenv.Load()
	jwtSecret := os.Getenv("JWT_SECRET")
	mux := http.NewServeMux()

	db, err := database.NewDB("database.json")
	if err != nil {
		log.Fatal(err)
	}

	cfg := &apiConfig{
		fileserverhits: 0,
		DB:             db,
		Jwt:            jwtSecret,
	}

	mux.Handle("/", http.FileServer(http.Dir(".")))

	mux.Handle("/app/*", http.StripPrefix("/app",
		cfg.middlewareMetricsInc(http.FileServer(http.Dir("."))),
	))

	mux.HandleFunc("GET /api/healthz", healthHandler)
	mux.HandleFunc("GET /admin/metrics", cfg.handlerMetrics)
	mux.HandleFunc("GET /api/reset", cfg.handlerReset)
	mux.HandleFunc("POST /api/chirps", cfg.postChirp)
	mux.HandleFunc("GET /api/chirps", cfg.getChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", cfg.getSingleChirp)
	mux.HandleFunc("POST /api/users", cfg.postUser)
	mux.HandleFunc("POST /api/login", cfg.login)
	mux.HandleFunc("PUT /api/users", cfg.putUser)
	mux.HandleFunc("POST /api/refresh", cfg.refreshToken)
	mux.HandleFunc("POST /api/revoke", cfg.revoke)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", cfg.deleteChirps)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	fmt.Println("Serving on port " + port)
	log.Fatal(server.ListenAndServe())
}

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
