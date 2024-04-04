package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/RealMotz/chirpy/internal/database"
)

func (cfg *apiConfig) createChirp(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	decodedChirp := database.Chirp{}
	err := decoder.Decode(&decodedChirp)
	if err != nil {
		handleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	err = validateChirp(decodedChirp)
	if err != nil {
		handleErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	err = cfg.db.CreateIfNotExits()
	if err != nil {
		handleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	chirp, err := cfg.db.CreateChirp(cleanBody(decodedChirp.Body))
	if err != nil {
		handleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}
	handleJsonResponse(w, http.StatusCreated, chirp)
}

func validateChirp(chirp database.Chirp) error {
	if len(chirp.Body) > 140 {
		return errors.New("chirp is too long")
	}
	return nil
}

func cleanBody(body string) string {
	profanities := []string{
		"kerfuffle",
		"sharbert",
		"fornax",
	}

	words := strings.Split(body, " ")
	for i := 0; i < len(words); i++ {
		for _, profanity := range profanities {
			if profanity == strings.ToLower(words[i]) {
				words[i] = "****"
				break
			}
		}
	}
	return strings.Join(words, " ")
}
