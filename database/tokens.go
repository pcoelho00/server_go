package database

import "time"

type RefreshToken struct {
	UserID         int       `json:"user_id"`
	Token          string    `json:"refresh_token"`
	ExpirationTime time.Time `json:"expiration_time"`
}

func (db *DB) SaveRefreshToken(token string, user_id int) error {
	dbStructure, err := db.LoadDB()
	if err != nil {
		return err
	}

	NewRefreshToken := RefreshToken{
		UserID:         user_id,
		Token:          token,
		ExpirationTime: time.Now().AddDate(0, 0, 60),
	}

	dbStructure.RefreshTokens[token] = NewRefreshToken
	err = db.WriteDB(dbStructure)
	if err != nil {
		return err
	}

	return nil

}

func (db *DB) FindRefreshToken(token string) (int, error) {

	dbStructure, err := db.LoadDB()
	if err != nil {
		return 0, err
	}

	now := time.Now()
	refreshToken, ok := dbStructure.RefreshTokens[token]
	if !ok {
		return 0, nil
	}

	if now.After(refreshToken.ExpirationTime) {
		return 0, nil
	}

	return refreshToken.UserID, nil

}

func (db *DB) RevokeRefreshToken(token string) error {

	dbStructure, err := db.LoadDB()
	if err != nil {
		return err
	}

	delete(dbStructure.RefreshTokens, token)

	err = db.WriteDB(dbStructure)
	if err != nil {
		return err
	}

	return nil

}
