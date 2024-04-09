package main

import (
	"net/http"
	"strconv"
)

func (cfg *apiConfig) getChirps(w http.ResponseWriter, r *http.Request) {
	authorIdParam := r.URL.Query().Get("author_id")
	authorId, err := getAuthorId(authorIdParam)
	if err != nil {
		handleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}
	sortParam := r.URL.Query().Get("sort")
	var descSort bool = false
	if sortParam == "desc" {
		descSort = true
	}

	chirps, err := cfg.db.GetChirps(authorId, descSort)
	if err != nil {
		handleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	handleJsonResponse(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) getChirp(w http.ResponseWriter, r *http.Request) {
	chirp, err := cfg.db.GetChirp(r.PathValue("id"))
	if err != nil {
		handleErrorResponse(w, http.StatusNotFound, err)
		return
	}

	handleJsonResponse(w, http.StatusOK, chirp)
}

func getAuthorId(query string) (int, error) {
	if query == "" {
		return 0, nil
	}
	id, err := strconv.Atoi(query)
	if err != nil {
		return 0, err
	}
	return id, nil
}
