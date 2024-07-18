package database

func (db *DB) CreateUser(body string) (User, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	dbStruct, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	newId := len(dbStruct.Users) + 1
	newUser := User{
		Id:    newId,
		Email: body,
	}

	dbStruct.Users[newId] = newUser

	err = db.writeDB(dbStruct)
	if err != nil {
		return User{}, err
	}

	return newUser, nil
}
