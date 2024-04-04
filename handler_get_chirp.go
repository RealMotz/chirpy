package main

import "net/http"

func (cfg *apiConfig) getChirps(w http.ResponseWriter, r *http.Request) {
	err := cfg.db.CreateIfNotExits()
	if err != nil {
		handleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	chirps, err := cfg.db.GetChirps()
	if err != nil {
		handleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	handleJsonResponse(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) getChirp(w http.ResponseWriter, r *http.Request) {
	err := cfg.db.CreateIfNotExits()
	if err != nil {
		handleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	chirp, err := cfg.db.GetChirp(r.PathValue("id"))
	if err != nil {
		handleErrorResponse(w, http.StatusNotFound, err)
		return
	}

	handleJsonResponse(w, http.StatusOK, chirp)
}
