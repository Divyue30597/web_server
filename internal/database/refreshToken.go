package database

import "time"

type RefreshToken struct {
	UserId    int       `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

func (db *DB) SaveRefreshToken(userId int, token string) error {
	dbStruct, err := db.loadDB()
	if err != nil {
		return err
	}

	refreshToken := RefreshToken{
		UserId:    userId,
		Token:     token,
		ExpiresAt: time.Now().Add(60 * 24 * time.Hour).UTC(),
	}
	dbStruct.RefreshToken[token] = refreshToken

	err = db.writeDB(dbStruct)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) RevokeRefreshToken(token string) error {
	dbStruct, err := db.loadDB()
	if err != nil {
		return err
	}

	// deleting refresh token
	delete(dbStruct.RefreshToken, token)

	err = db.writeDB(dbStruct)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) GetUserFromRefreshToken(token string) (User, error) {
	dbStruct, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	// check for refresh token
	refreshToken, ok := dbStruct.RefreshToken[token]
	if !ok {
		return User{}, ErrNotExist
	}

	if refreshToken.ExpiresAt.Before(time.Now().UTC()) {
		return User{}, ErrNotExist
	}

	user, err := db.GetUserById(refreshToken.UserId)
	if err != nil {
		return User{}, err
	}

	return user, nil

}