package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthRequest struct {
	Password         string `json:"password"`
	Email            string `json:"email"`
	ExpiresInSeconds int    `json:"expires_in_seconds"`
}

type AuthResponse struct {
	Id           int    `json:"id"`
	Email        string `json:"email"`
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

type TokenType int

const (
	RefreshToken TokenType = iota + 1
	AccessToken
)

func (i TokenType) String() string {
	return [...]string{"chirpy-refresh", "chirpy-access"}[i-1]
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

	accessToken, err := cfg.createToken(AccessToken.String(), user.Id, time.Duration(time.Hour))
	if err != nil {
		handleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}
	oneDay := time.Duration(time.Hour * 24 * 60)
	refreshToken, err := cfg.createToken(RefreshToken.String(), user.Id, oneDay)
	if err != nil {
		handleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	handleJsonResponse(w, http.StatusOK, AuthResponse{
		Id:           user.Id,
		Email:        user.Email,
		Token:        accessToken,
		RefreshToken: refreshToken,
	})
}

func (cfg *apiConfig) createToken(issuerName string, subject int, expiration time.Duration) (string, error) {
	// Create the claims
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    issuerName,
		Subject:   strconv.Itoa(subject),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(cfg.jwtSecret)
	if err != nil {
		fmt.Printf("Error signing jwt token")
		return "", err
	}

	return ss, nil
}

func (cfg *apiConfig) refreshLoginToken(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	splitAuth := strings.Split(auth, " ")
	if len(splitAuth) < 2 || splitAuth[0] != "Bearer" {
		handleErrorResponse(w, http.StatusUnauthorized, errors.New("malformed authorization header"))
		return
	}
	token := splitAuth[1]

	parsedToken, err := cfg.parseToken(token)
	if err != nil {
		handleErrorResponse(w, http.StatusUnauthorized, err)
		return
	}

	err = cfg.verifyIssuer(parsedToken, RefreshToken.String())
	if err != nil {
		handleErrorResponse(w, http.StatusUnauthorized, err)
		return
	}

	_, err = cfg.db.GetRevokedToken(token)
	if err == nil {
		handleErrorResponse(w, http.StatusUnauthorized, errors.New("token has been revoked"))
		return
	}

	id, err := cfg.getSubject(parsedToken)
	if err != nil {
		handleErrorResponse(w, http.StatusUnauthorized, err)
		return
	}

	accessToken, err := cfg.createToken(AccessToken.String(), id, time.Duration(time.Hour))
	if err != nil {
		handleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	type customResponse struct {
		Token string `json:"token"`
	}

	handleJsonResponse(w, http.StatusOK, customResponse{
		Token: accessToken,
	})
}

func (cfg *apiConfig) revokeLoginToken(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	splitAuth := strings.Split(auth, " ")
	if len(splitAuth) < 2 || splitAuth[0] != "Bearer" {
		handleErrorResponse(w, http.StatusUnauthorized, errors.New("malformed authorization header"))
		return
	}
	token := splitAuth[1]

	parsedToken, err := cfg.parseToken(token)
	if err != nil {
		handleErrorResponse(w, http.StatusUnauthorized, err)
		return
	}

	err = cfg.verifyIssuer(parsedToken, RefreshToken.String())
	if err != nil {
		handleErrorResponse(w, http.StatusUnauthorized, err)
		return
	}

	err = cfg.db.Revoke(token)
	if err != nil {
		handleErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	handleJsonResponse(w, http.StatusOK, struct{}{})
}

func (cfg *apiConfig) parseToken(token string) (*jwt.Token, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return cfg.jwtSecret, nil
	})
	if err != nil {
		fmt.Printf("Error parsing jwt token")
		return &jwt.Token{}, err
	}

	return parsedToken, nil
}

func (cfg *apiConfig) verifyIssuer(token *jwt.Token, tokenIssuer string) error {
	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		fmt.Printf("Error retrieving jwt issuer")
		return err
	}

	if issuer != tokenIssuer {
		return errors.New("token not authorized")
	}

	return nil
}

func (cfg *apiConfig) getSubject(token *jwt.Token) (int, error) {
	subject, err := token.Claims.GetSubject()
	if err != nil {
		fmt.Printf("Error retrieving jwt subject")
		return 0, err
	}

	parsedSubject, err := strconv.Atoi(subject)
	if err != nil {
		fmt.Printf("Error parsing jwt subject")
		return 0, err
	}

	return parsedSubject, nil
}
