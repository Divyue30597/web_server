package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

const port = "8080"

var profanity = []string{"kerfuffle", "sharbert", "fornax"}

type apiConfig struct {
	fileserverhits int
}

type DB struct {
	path string
	mux  *sync.Mutex
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}

type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
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
	mux.HandleFunc("POST /api/chirps", postChirp)
	mux.HandleFunc("GET /api/chirps", getChirps)

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

func NewDB(path string) (*DB, error) {
	db := &DB{
		path: path,
		mux:  &sync.Mutex{},
	}

	if err := db.ensureDB(); err != nil {
		return nil, err
	}

	return db, nil
}

func (db *DB) ensureDB() error {
	// check if the file exists
	_, err := os.Stat(db.path)
	if os.IsNotExist(err) {
		file, err := os.Create(db.path)
		if err != nil {
			return err
		}

		file.Write([]byte(`
			{
				"chirps": {}
			}
		`))

		file.Close()
	}

	return nil
}

func (db *DB) loadDB() (DBStructure, error) {
	var dbStruct DBStructure
	data, err := os.ReadFile(db.path)
	if err != nil {
		fmt.Println(err, "readfile")
		return DBStructure{}, err
	}

	err = json.Unmarshal(data, &dbStruct)
	if err != nil {
		fmt.Println(err, "json unmarshal")
		return DBStructure{}, err
	}

	return dbStruct, nil
}

func (db *DB) writeDB(dbStructure DBStructure) error {
	data, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}

	err = os.WriteFile(db.path, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) CreateChirp(body string) (Chirp, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	// we read the database
	dbStruct, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	newId := len(dbStruct.Chirps) + 1
	newChirp := Chirp{
		Id:   newId,
		Body: body,
	}

	dbStruct.Chirps[newId] = newChirp

	err = db.writeDB(dbStruct)
	if err != nil {
		return Chirp{}, err
	}

	return newChirp, nil
}

func (db *DB) GetChirps() ([]Chirp, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	dbStruct, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	chirps := make([]Chirp, 0, len(dbStruct.Chirps))
	for _, v := range dbStruct.Chirps {
		chirps = append(chirps, v)
	}

	return chirps, nil
}

func postChirp(w http.ResponseWriter, r *http.Request) {
	db, err := NewDB("chirps.json")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating DB")
	}
	type postchirp struct {
		Body string `json:"body"`
	}

	const validChirpLen = 140

	decoder := json.NewDecoder(r.Body)

	params := postchirp{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding JSON body")
		return
	}

	if len(params.Body) > validChirpLen {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	data, _ := db.CreateChirp(cleanChirp(params.Body))

	respondWithJSON(w, http.StatusCreated, data)
}

func getChirps(w http.ResponseWriter, r *http.Request) {
	db, err := NewDB("chirps.json")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating DB")
	}

	chirps, err := db.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error getting chirps")
	}

	respondWithJSON(w, http.StatusOK, chirps)
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

func respondWithError(w http.ResponseWriter, code int, message string) {
	if code > 499 {
		log.Fatalf("Responding with 5XX error: %s", message)
	}
	type errorResponse struct {
		Err string `json:"error"`
	}

	respondWithJSON(w, code, errorResponse{Err: message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(payload)
	if err != nil {
		log.Fatalf("error marshalling json: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(code)
	w.Write(data)
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
