package database

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func (db *DB) CreateUser(body User) (User, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	dbStruct, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	newId := len(dbStruct.Users) + 1

	newUser := User{
		Id:       newId,
		Email:    body.Email,
		Password: body.Password,
	}

	dbStruct.Users[newId] = newUser

	err = db.writeDB(dbStruct)
	if err != nil {
		fmt.Println(err, "writeDB error")
		return User{}, err
	}

	return newUser, nil
}

func (db *DB) GetUsers() ([]User, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	dbStruct, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	users := make([]User, 0, len(dbStruct.Users))
	for _, user := range dbStruct.Users {
		users = append(users, user)
	}

	return users, nil
}

func (db *DB) VerifyUser(body User) (User, error) {
	users, err := db.GetUsers()
	if err != nil {
		return User{}, err
	}

	for _, user := range users {
		if user.Email == body.Email {
			err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
			if err != nil {
				return User{}, err
			}
			return user, nil
		}
	}

	return User{}, errors.New("user not found in db")
}
