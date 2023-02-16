package auth

import (
	"errors"

	"github.com/golang-jwt/jwt/v4"
)

var ErrUserNotFound = errors.New("user not found")

type User struct {
	ID    int
	Name  string
	Email string
}

type Claims struct {
	UserID int
	Name   string
	Email  string
	jwt.RegisteredClaims
}
