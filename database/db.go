package database

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/pcoelho00/server_go/auth"
)

type ChirpsMsg struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}

type DB struct {
	path string
	mux  *sync.RWMutex
}

type User struct {
	Id       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RefreshToken struct {
	UserID         int       `json:"user_id"`
	Token          string    `json:"refresh_token"`
	ExpirationTime time.Time `json:"expiration_time"`
}

type PublicUser struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
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

func (db *DB) CreateUser(email, password string) (User, error) {
	dbStructure, err := db.LoadDB()
	if err != nil {
		log.Println("error loading")
		return User{}, err
	}

	last_id := len(dbStructure.Users)

	passhash, err := auth.HashPassword(password)
	if err != nil {
		log.Println("error creating password hash")
		return User{}, err
	}

	NewUser := User{
		Id:       last_id + 1,
		Email:    email,
		Password: passhash,
	}

	dbStructure.Users[NewUser.Id] = NewUser

	err = db.WriteDB(dbStructure)
	if err != nil {
		log.Println("couldn't save the database")
		return User{}, err
	}

	return NewUser, nil
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

func (db *DB) GetUsers() ([]PublicUser, error) {
	dbStructure, err := db.LoadDB()

	if err != nil {
		return []PublicUser{}, err
	}

	users := make([]PublicUser, 0)
	for _, user := range dbStructure.Users {
		p_user := PublicUser{
			Id:    user.Id,
			Email: user.Email,
		}
		users = append(users, p_user)
	}

	return users, nil

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

func (db *DB) GetUser(id int) (User, error) {

	dbStructure, err := db.LoadDB()
	if err != nil {
		return User{}, err
	}

	user, ok := dbStructure.Users[id]
	if !ok {
		return User{}, fmt.Errorf("User with id %d not found", id)
	}

	return user, nil
}

func (db *DB) GetPublicUser(id int) (PublicUser, error) {

	dbStructure, err := db.LoadDB()
	if err != nil {
		return PublicUser{}, err
	}

	user, ok := dbStructure.Users[id]
	if !ok {
		return PublicUser{}, fmt.Errorf("User with id %d not found", id)
	}

	return PublicUser{Id: user.Id, Email: user.Email}, nil
}

func (db *DB) GetUserFromLogin(email, password string) (User, error) {

	dbStructure, err := db.LoadDB()
	if err != nil {
		return User{}, err
	}

	search_id := 0

	for id, user := range dbStructure.Users {
		if email == user.Email {
			search_id = id
			break
		}
	}
	if search_id == 0 {
		return User{}, fmt.Errorf("User %s doesn't exist", email)
	}

	user, ok := dbStructure.Users[search_id]
	if !ok {
		return User{}, fmt.Errorf("User with id %d not found", search_id)
	}

	pass_check := auth.CheckPasswordHash(password, user.Password)

	if pass_check {
		return user, nil
	} else {
		return User{}, fmt.Errorf("password doesn't match")
	}

}

func (db *DB) UpdateUser(id int, email, password string) (PublicUser, error) {

	dbStructure, err := db.LoadDB()
	if err != nil {
		return PublicUser{}, err
	}

	user, ok := dbStructure.Users[id]
	if !ok {
		return PublicUser{}, fmt.Errorf("User with id %d not found", id)
	}

	passhash, err := auth.HashPassword(password)
	if err != nil {
		println("error creating password hash")
		return PublicUser{}, err
	}

	UpdatedUser := User{
		Id:       user.Id,
		Email:    email,
		Password: passhash,
	}

	dbStructure.Users[id] = UpdatedUser

	err = db.WriteDB(dbStructure)
	if err != nil {
		return PublicUser{}, err
	}

	return PublicUser{Id: UpdatedUser.Id, Email: UpdatedUser.Email}, nil

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
