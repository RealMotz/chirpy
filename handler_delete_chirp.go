package main

import (
	"errors"
	"net/http"

	"github.com/RealMotz/chirpy/internal/auth"
)

func (cfg *apiConfig) deleteChirp(w http.ResponseWriter, r *http.Request) {
	// authenticate the user
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

	chirp, err := cfg.db.GetChirp(r.PathValue("id"))
	if err != nil {
		handleErrorResponse(w, http.StatusNotFound, err)
		return
	}

	// if the user is not the author, return 403
	if chirp.AuthorId != authorId {
		handleErrorResponse(w, http.StatusForbidden, errors.New("forbidden resource"))
		return
	}

	err = cfg.db.DeleteChirp(r.PathValue("id"))
	if err != nil {
		handleErrorResponse(w, http.StatusNotFound, err)
		return
	}

	// return 200 after successfully deleting a chirp
	handleJsonResponse(w, http.StatusOK, struct{}{})
}
