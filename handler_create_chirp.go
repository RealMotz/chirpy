package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/RealMotz/chirpy/internal/auth"
	"github.com/RealMotz/chirpy/internal/database"
)

func (cfg *apiConfig) createChirp(w http.ResponseWriter, r *http.Request) {
	token, err := auth.FetchAuthHeader(r.Header.Get("Authorization"), "Bearer")
	if err != nil {
		handleErrorResponse(w, http.StatusUnauthorized, err)
		return
	}

	parsedToken, err := cfg.parseToken(token)
	if err != nil {
		handleErrorResponse(w, http.StatusUnauthorized, err)
		return
	}

	err = cfg.verifyIssuer(parsedToken, auth.AccessToken.String())
	if err != nil {
		handleErrorResponse(w, http.StatusUnauthorized, err)
		return
	}

	authorId, err := cfg.getSubject(parsedToken)
	if err != nil {
		handleErrorResponse(w, http.StatusUnauthorized, err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	decodedChirp := database.Chirp{}
	err = decoder.Decode(&decodedChirp)
	if err != nil {
		handleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	err = validateChirp(decodedChirp)
	if err != nil {
		handleErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	chirp, err := cfg.db.CreateChirp(cleanBody(decodedChirp.Body), authorId)
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
