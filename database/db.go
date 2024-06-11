package database

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type ChirpsMsg struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]ChirpsMsg `json:"chirps"`
}

func (db *DB) ensureDB() error {
	err := os.WriteFile(db.path, []byte("{}"), 0777)
	if err != nil {
		return err
	}

	return nil
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

func (db *DB) UpdateDB(msg ChirpsMsg) (DBStructure, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	dbStructure, err := db.LoadDB()
	fmt.Printf("%v", dbStructure)
	if err != nil {
		return DBStructure{}, err
	}

	dbStructure.Chirps[msg.Id] = msg
	return dbStructure, nil

}

func (db *DB) CreateChirp(body string) (ChirpsMsg, error) {

	dbStructure, err := db.LoadDB()
	if err != nil {
		println("error loading")
		return ChirpsMsg{}, err
	}

	last_id := len(dbStructure.Chirps)
	Chirp := ChirpsMsg{
		Id:   last_id + 1,
		Body: body,
	}

	return Chirp, nil

}

func (db *DB) GetChirps() ([]ChirpsMsg, error) {
	dbStructure, err := db.LoadDB()

	if err != nil {
		return []ChirpsMsg{}, err
	}

	msgs := make([]ChirpsMsg, 0)
	for _, msg := range dbStructure.Chirps {
		msgs = append(msgs, msg)
	}

	return msgs, nil
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

	// Perform any additional initialization logic here

	return db, nil
}
