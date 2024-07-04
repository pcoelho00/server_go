package database

import (
	"encoding/json"
	"os"
	"sync"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DBStructure struct {
	Chirps        map[int]ChirpsMsg       `json:"chirps"`
	Users         map[int]User            `json:"users"`
	RefreshTokens map[string]RefreshToken `json:"refresh_token"`
}

func (db *DB) ensureDB() error {
	err := os.WriteFile(db.path, []byte("{}"), 0777)
	if err != nil {
		return err
	}

	return nil
}

func NewDB(path string) (*DB, error) {

	// Create a new DB instance
	db := &DB{
		path: path,
		mux:  &sync.RWMutex{},
	}
	// Check if the file exists
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		err := db.ensureDB()
		if err != nil {
			// File does not exist, handle the error according
			return nil, err
		}
	}

	return db, nil
}

func (db *DB) LoadDB() (DBStructure, error) {

	data, err := os.ReadFile(db.path)
	if err != nil {
		return DBStructure{}, err
	}

	var dbStructure DBStructure
	err = json.Unmarshal(data, &dbStructure)
	if err != nil {
		return DBStructure{}, err
	}

	if dbStructure.Chirps == nil {
		dbStructure.Chirps = make(map[int]ChirpsMsg)
	}

	if dbStructure.Users == nil {
		dbStructure.Users = make(map[int]User)
	}

	if dbStructure.RefreshTokens == nil {
		dbStructure.RefreshTokens = make(map[string]RefreshToken)
	}

	return dbStructure, nil

}

func (db *DB) WriteDB(dbStructure DBStructure) error {
	db.mux.Lock()
	defer db.mux.Unlock()

	b, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}

	os.WriteFile(db.path, b, 0777)
	return nil

}
