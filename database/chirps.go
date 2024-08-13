package database

import (
	"fmt"
	"sort"
)

type ChirpsMsg struct {
	Id       int    `json:"id"`
	AuthorId int    `json:"author_id"`
	Body     string `json:"body"`
}

type ById []ChirpsMsg

func (a ById) Len() int           { return len(a) }
func (a ById) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ById) Less(i, j int) bool { return a[i].Id < a[j].Id }

func (db *DB) WriteChirpsToDB(msg ChirpsMsg) (DBStructure, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	dbStructure, err := db.LoadDB()
	if err != nil {
		return DBStructure{}, err
	}

	dbStructure.Chirps[msg.Id] = msg
	dbStructure.DBInfo.LastChirpID = msg.Id
	return dbStructure, nil

}

func (db *DB) CreateChirp(author_id int, body string) (ChirpsMsg, error) {

	dbStructure, err := db.LoadDB()
	if err != nil {
		println("error loading")
		return ChirpsMsg{}, err
	}

	last_id := dbStructure.DBInfo.LastChirpID
	Chirp := ChirpsMsg{
		Id:       last_id + 1,
		AuthorId: author_id,
		Body:     body,
	}

	return Chirp, nil

}

func (db *DB) GetChirps(id int) ([]ChirpsMsg, error) {
	dbStructure, err := db.LoadDB()

	if err != nil {
		return []ChirpsMsg{}, err
	}

	msgs := make([]ChirpsMsg, 0)

	if id < 0 {
		return []ChirpsMsg{}, fmt.Errorf("id must be greater than zero")
	}

	if id == 0 {
		for _, msg := range dbStructure.Chirps {
			msgs = append(msgs, msg)
		}
	} else {
		for _, msg := range dbStructure.Chirps {
			if msg.AuthorId == id {
				msgs = append(msgs, msg)
			}
		}
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

func (db *DB) DeleteChirp(author_id, id int) error {

	dbStructure, err := db.LoadDB()
	if err != nil {
		return err
	}

	msg, ok := dbStructure.Chirps[id]
	if !ok {
		return fmt.Errorf("msg id %v doesn't exist", id)
	} else if msg.AuthorId != author_id {
		return fmt.Errorf("author_id %v doesn't match", author_id)
	} else {
		delete(dbStructure.Chirps, id)
		return nil
	}

}

func (db *DB) SortMsgs(msgs []ChirpsMsg, by string) ([]ChirpsMsg, error) {
	if by == "asc" {
		sort.Sort(ById(msgs))
	} else if by == "desc" {
		sort.Sort(sort.Reverse(ById(msgs)))
	} else {
		return []ChirpsMsg{}, fmt.Errorf("invalid sorting type")
	}
	return msgs, nil
}
