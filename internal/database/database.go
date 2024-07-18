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
	Users  map[int]User  `json:"users"`
}

type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}

type User struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
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

		defer file.Close()

		initialData := DBStructure{
			Chirps: make(map[int]Chirp),
			Users:  make(map[int]User),
		}

		jsonData, err := json.MarshalIndent(initialData, "", "	")
		if err != nil {
			return err
		}

		if _, err := file.Write(jsonData); err != nil {
			return err
		}
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
