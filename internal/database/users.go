package database

import (
	"errors"
	"fmt"
)

var ErrAlreadyExists = errors.New("already exists")

func (db *DB) CreateUser(email, password string) (User, error) {
	// check if the email already exists
	if _, err := db.GetUserByEmail(email); !errors.Is(err, ErrNotExist) {
		return User{}, ErrAlreadyExists
	}

	dbStruct, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	newId := len(dbStruct.Users) + 1
	newUser := User{
		Id:       newId,
		Email:    email,
		Password: password,
	}

	dbStruct.Users[newId] = newUser

	err = db.writeDB(dbStruct)
	if err != nil {
		return User{}, err
	}

	return newUser, nil
}

func (db *DB) GetUserByEmail(email string) (User, error) {
	dbStruct, err := db.loadDB()
	if err != nil {
		fmt.Println(err)
		return User{}, err
	}

	for _, user := range dbStruct.Users {
		if user.Email == email {
			return user, nil
		}
	}

	return User{}, ErrNotExist
}

func (db *DB) GetUserById(id int) (User, error) {
	dbStruct, err := db.loadDB()
	if err != nil {
		return User{}, nil
	}

	user, ok := dbStruct.Users[id]
	if !ok {
		return User{}, nil
	}

	return user, nil
}

func (db *DB) GetUsers() ([]User, error) {
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

// func (db *DB) VerifyUser(body User) (User, error) {
// 	users, err := db.GetUsers()
// 	if err != nil {
// 		return User{}, err
// 	}

// 	for _, user := range users {
// 		if user.Email == body.Email {
// 			err := auth.VerifyPassword(user.Password, body.Password)
// 			if err != nil {
// 				return User{}, err
// 			}
// 			return user, nil
// 		}
// 	}

// 	return User{}, errors.New("user not found in db")
// }
