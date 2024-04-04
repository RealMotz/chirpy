package main

import "net/http"

func (cfg *apiConfig) getChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetChirps()
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
