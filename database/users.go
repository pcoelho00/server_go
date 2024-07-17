package database

import (
	"fmt"
	"log"

	"github.com/pcoelho00/server_go/auth"
)

type User struct {
	Id       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type PublicUser struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
}

func (db *DB) CreateUser(email, password string) (User, error) {
	dbStructure, err := db.LoadDB()
	if err != nil {
		log.Println("error loading")
		return User{}, err
	}

	passhash, err := auth.HashPassword(password)
	if err != nil {
		log.Println("error creating password hash")
		return User{}, err
	}

	last_id := dbStructure.DBInfo.LastUserID
	NewUser := User{
		Id:       last_id + 1,
		Email:    email,
		Password: passhash,
	}

	dbStructure.Users[NewUser.Id] = NewUser
	dbStructure.DBInfo.LastUserID = NewUser.Id

	err = db.WriteDB(dbStructure)
	if err != nil {
		log.Println("couldn't save the database")
		return User{}, err
	}

	return NewUser, nil
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
