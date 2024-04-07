package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/mail"
	"strings"

	"github.com/RealMotz/chirpy/internal/database"
	"golang.org/x/crypto/bcrypt"
)

func (cfg *apiConfig) createUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	req := database.User{}
	err := decoder.Decode(&req)
	if err != nil {
		handleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	user, httpError := AreEmailAndPasswordValid(req.Email, req.Password)
	if httpError.Error != nil {
		handleErrorResponse(w, httpError.Code, httpError.Error)
		return
	}

	res, httpError := cfg.db.CreateUser(user.Email, user.Password)
	if httpError.Error != nil {
		handleErrorResponse(w, httpError.Code, httpError.Error)
		return
	}
	handleJsonResponse(w, http.StatusCreated, res)
}

func AreEmailAndPasswordValid(email string, pwd string) (database.User, database.HttpError) {
	email = strings.Trim(email, " ")
	if !IsEmailValid(email) {
		return database.User{}, database.HttpError{
			Error: errors.New("invalid email"),
			Code:  http.StatusBadRequest,
		}
	}

	encryptedPwd, err := bcrypt.GenerateFromPassword([]byte(pwd), 10)
	if err != nil {
		log.Printf("Error encrypting password: %v", err)
		return database.User{}, database.HttpError{
			Error: err,
			Code:  http.StatusInternalServerError,
		}
	}

	return database.User{
			Email:    email,
			Password: string(encryptedPwd),
		}, database.HttpError{
			Error: nil,
		}
}

func IsEmailValid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
