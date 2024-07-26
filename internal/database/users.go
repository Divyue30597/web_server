package database

import (
	"errors"
)

var ErrAlreadyExists = errors.New("already exists")

type User struct {
	Id       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

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

func (db *DB) UpdateUser(id int, email, hashedPassword string) (User, error) {
	// No update logic yet
	dbStruct, err := db.loadDB()
	if err != nil {
		return User{}, nil
	}

	user, ok := dbStruct.Users[id]
	if !ok {
		return User{}, ErrNotExist
	}

	user.Email = email
	user.Password = hashedPassword
	dbStruct.Users[id] = user

	err = db.writeDB(dbStruct)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

// func (db *DB) SaveRefreshToken(id int, refreshToken string) error {
// 	dbStruct, err := db.loadDB()
// 	if err != nil {
// 		return err
// 	}
// 	user := dbStruct.Users[id]
// 	user.RefreshToken = refreshToken
// 	expirationTime := time.Now().Add(time.Duration(60 * 24 * 60 * 60 * time.Second)).UTC().Format(time.RFC3339)
// 	user.RefreshTokenExpirationTime = expirationTime
// 	dbStruct.Users[id] = user
// 	err = db.writeDB(dbStruct)
// 	if err != nil {
// 		return err
// 	}
// 	return err
// }

// func (db *DB) GetUserFromRefreshToken(refreshToken string) (User, error) {
// 	dbStruct, err := db.loadDB()
// 	if err != nil {
// 		return User{}, nil
// 	}
// 	for _, user := range dbStruct.Users {
// 		if user.RefreshToken == refreshToken {
// 			return user, nil
// 		}
// 	}
// 	return User{}, err
// }

// func (db *DB) UpdateUserFromRefreshToken(refreshToken string) error {
// 	dbStruct, err := db.loadDB()
// 	if err != nil {
// 		return err
// 	}
// 	for _, user := range dbStruct.Users {
// 		if user.RefreshToken == refreshToken {
// 			parseTime, err := time.Parse(time.RFC3339, user.RefreshTokenExpirationTime)
// 			if err != nil {
// 				return err
// 			}
// 			if time.Now().UTC().After(parseTime) {
// 				user.RefreshToken = ""
// 				user.RefreshTokenExpirationTime = ""
// 				dbStruct.Users[user.Id] = user
// 				return errors.New("refresh token expired")
// 			}
// 		}
// 	}
// 	err = db.writeDB(dbStruct)
// 	return err
// }
