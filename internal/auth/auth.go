package auth

import (
	"errors"
	"strings"
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

var ErrorInvalidAuthHeader error = errors.New("malformed authorization header")

type TokenType int

const (
	RefreshToken TokenType = iota + 1
	AccessToken
)

func (i TokenType) String() string {
	return [...]string{"chirpy-refresh", "chirpy-access"}[i-1]
}

func FetchAuthHeader(header string) (string, error) {
	splitAuth := strings.Split(header, " ")
	if len(splitAuth) < 2 || splitAuth[0] != "Bearer" {
		return "", ErrorInvalidAuthHeader
	}
	return splitAuth[1], nil
}
