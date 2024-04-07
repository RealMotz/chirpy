package main

import (
	"encoding/json"
	"net/http"
)

type AuthRequest struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

func (cfg *apiConfig) login(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	req := AuthRequest{}
	err := decoder.Decode(&req)
	if err != nil {
		handleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	user, err := cfg.db.GetUser(req.Email, []byte(req.Password))
	if err != nil {
		handleErrorResponse(w, http.StatusUnauthorized, err)
		return
	}

	handleJsonResponse(w, http.StatusOK, user)
}
