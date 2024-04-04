package main

import (
	"encoding/json"
	"net/http"
	"net/mail"
	"strings"

	"github.com/RealMotz/chirpy/internal/database"
)

func (cfg *apiConfig) createUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	decodedUser := database.User{}
	err := decoder.Decode(&decodedUser)
	if err != nil {
		handleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	trimmedEmail := strings.Trim(decodedUser.Email, " ")

	if !isEmailValid(trimmedEmail) {
		handleErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	chirp, err := cfg.db.CreateUser(trimmedEmail)
	if err != nil {
		handleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}
	handleJsonResponse(w, http.StatusCreated, chirp)
}

func isEmailValid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
