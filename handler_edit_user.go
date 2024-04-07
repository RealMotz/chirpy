package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/RealMotz/chirpy/internal/database"
	"github.com/golang-jwt/jwt/v5"
)

func (cfg *apiConfig) updateUser(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	token = strings.Split(token, " ")[1]

	parsedToken, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return cfg.jwtSecret, nil
	})
	if err != nil {
		fmt.Printf("Error parsing jwt token")
		handleErrorResponse(w, http.StatusUnauthorized, err)
		return
	}

	userId, err := parsedToken.Claims.GetSubject()
	if err != nil {
		fmt.Printf("Error parsing jwt token")
		handleErrorResponse(w, http.StatusUnauthorized, err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	req := database.User{}
	err = decoder.Decode(&req)
	if err != nil {
		handleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	user, httpError := AreEmailAndPasswordValid(req.Email, req.Password)
	if httpError.Error != nil {
		handleErrorResponse(w, httpError.Code, httpError.Error)
		return
	}

	id, err := strconv.Atoi(userId)
	if err != nil {
		handleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}
	updatedUser := database.User{
		Id:       id,
		Email:    user.Email,
		Password: user.Password,
	}

	res, httpError := cfg.db.UpdateUser(updatedUser)
	if httpError.Error != nil {
		handleErrorResponse(w, httpError.Code, httpError.Error)
		return
	}
	handleJsonResponse(w, http.StatusOK, res)
}
