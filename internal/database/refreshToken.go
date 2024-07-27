package database

import "time"

type RefreshToken struct {
	UserID    int       `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

func (db *DB) SaveRefreshToken(userId int, token string) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	refreshToken := RefreshToken{
		UserID:    userId,
		Token:     token,
		ExpiresAt: time.Now().Add(60 * 24 * time.Hour),
	}
	dbStructure.RefreshTokens[token] = refreshToken

	err = db.writeDB(dbStructure)
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
	delete(dbStruct.RefreshTokens, token)

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
	refreshToken, ok := dbStruct.RefreshTokens[token]
	if !ok {
		return User{}, ErrNotExist
	}

	if refreshToken.ExpiresAt.Before(time.Now().UTC()) {
		return User{}, ErrNotExist
	}

	user, err := db.GetUserById(refreshToken.UserID)
	if err != nil {
		return User{}, err
	}

	return user, nil
}
