package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthRequest struct {
	Password         string `json:"password"`
	Email            string `json:"email"`
	ExpiresInSeconds int    `json:"expires_in_seconds"`
}

type AuthResponse struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
	Token string `json:"token"`
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

	expiration := time.Now().Add(time.Duration(req.ExpiresInSeconds) * time.Second)
	if req.ExpiresInSeconds <= 0 || req.ExpiresInSeconds > 24 {
		expiration = time.Now().Add(24 * time.Hour)
	}

	// Create the claims
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expiration),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    "chirpy",
		Subject:   strconv.Itoa(user.Id),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(cfg.jwtSecret)
	if err != nil {
		fmt.Printf("Error signing jwt token")
		handleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}
	handleJsonResponse(w, http.StatusOK, AuthResponse{
		Id:    user.Id,
		Email: user.Email,
		Token: ss,
	})
}
