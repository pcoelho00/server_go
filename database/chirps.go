package database

import "fmt"

type ChirpsMsg struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}

func (db *DB) WriteChirpsToDB(msg ChirpsMsg) (DBStructure, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	dbStructure, err := db.LoadDB()
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

func (db *DB) GetChirp(id int) (ChirpsMsg, error) {

	dbStructure, err := db.LoadDB()
	if err != nil {
		return ChirpsMsg{}, err
	}

	msg, ok := dbStructure.Chirps[id]
	if !ok {
		return ChirpsMsg{}, fmt.Errorf("chirp with id %d not found", id)
	}

	return msg, nil
}
