package database

import (
	"fmt"
)

type Chirp struct {
	Id       int    `json:"id"`
	Body     string `json:"body"`
	AuthorID int    `json:"author_id"`
}

func (db *DB) CreateChirp(body string, authorID int) (Chirp, error) {
	dbStruct, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	newId := len(dbStruct.Chirps) + 1
	newChirp := Chirp{
		Id:       newId,
		Body:     body,
		AuthorID: authorID,
	}

	dbStruct.Chirps[newId] = newChirp

	err = db.writeDB(dbStruct)
	if err != nil {
		return Chirp{}, err
	}

	return newChirp, nil
}

func (db *DB) GetChirps() ([]Chirp, error) {
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

func (db *DB) GetChirpsByAuthorID(authorID int) ([]Chirp, error) {
	dbStruct, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	chirps := make([]Chirp, 0, len(dbStruct.Chirps))

	for _, val := range dbStruct.Chirps {
		if val.AuthorID == authorID {
			chirps = append(chirps, val)
		}
	}

	return chirps, nil
}

func (db *DB) DeleteChirp(id int) error {
	dbStruct, err := db.loadDB()
	if err != nil {
		return err
	}

	_, ok := dbStruct.Chirps[id]
	if !ok {
		return fmt.Errorf("chirp with id %d not found", id)
	}

	delete(dbStruct.Chirps, id)

	err = db.writeDB(dbStruct)
	if err != nil {
		return err
	}

	return nil
}
