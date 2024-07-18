package database

import "fmt"

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
