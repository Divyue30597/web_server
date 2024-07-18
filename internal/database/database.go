package database

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

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

func (db *DB) GetSingleChirp(id int) (Chirp, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	dbStruct, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	chirp, ok := dbStruct.Chirps[id]
	if !ok {
		return Chirp{}, fmt.Errorf("chirp with id %d not found", id)
	}

	return chirp, nil

}
