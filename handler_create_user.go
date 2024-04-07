package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/mail"
	"strings"

	"github.com/RealMotz/chirpy/internal/database"
	"golang.org/x/crypto/bcrypt"
)

func (cfg *apiConfig) createUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	decodedUser := database.User{}
	err := decoder.Decode(&decodedUser)
	if err != nil {
		handleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	email := strings.Trim(decodedUser.Email, " ")
	if !isEmailValid(email) {
		handleErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	pwd, err := bcrypt.GenerateFromPassword([]byte(decodedUser.Password), 10)
	if err != nil {
		log.Printf("Error processing password: %v", err)
		handleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	user, userErr := cfg.db.CreateUser(email, pwd)
	if userErr.Error != nil {
		handleErrorResponse(w, userErr.Code, userErr.Error)
		return
	}
	handleJsonResponse(w, http.StatusCreated, user)
}

func isEmailValid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
