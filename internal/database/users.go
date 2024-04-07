package database

import (
	"errors"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
}

type UserError struct {
	Error error
	Code  int
}

func (db *DataBase) GetUser(email string, pwd []byte) (UserResponse, error) {
	dbData, err := db.Read()
	if err != nil {
		log.Fatal(err)
		return UserResponse{}, err
	}

	user := getUserByEmail(email, dbData.Users)
	if user == nil {
		return UserResponse{}, errors.New("user doesn't exist")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), pwd)
	if err != nil {
		return UserResponse{}, errors.New("passwords don't match")
	}

	return UserResponse{
		Id:    user.Id,
		Email: user.Email,
	}, nil
}

func (db *DataBase) CreateUser(email string, pwd []byte) (UserResponse, UserError) {
	dbData, err := db.Read()
	if err != nil {
		log.Fatal(err)
		return UserResponse{}, UserError{
			Error: err,
			Code:  http.StatusInternalServerError,
		}
	}

	user := getUserByEmail(email, dbData.Users)
	if user != nil {
		return UserResponse{}, UserError{
			Error: errors.New("email already exists"),
			Code:  http.StatusConflict,
		}
	}

	newId := len(dbData.Users) + 1
	dbData.Users[newId] = User{
		Id:       newId,
		Email:    email,
		Password: string(pwd),
	}
	db.Write(dbData)
	return UserResponse{
			Id:    newId,
			Email: email,
		}, UserError{
			Error: nil,
		}
}

func getUserByEmail(email string, users map[int]User) *User {
	for _, user := range users {
		if user.Email == email {
			return &user
		}
	}
	return nil
}
