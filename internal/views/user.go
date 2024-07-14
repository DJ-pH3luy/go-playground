package views

import "github.com/golang-jwt/jwt/v5"

type User struct {
	Id      uint   `json:"id"`
	Name    string `json:"username"`
	Email   string `json:"email"`
	IsAdmin bool   `json:"isAdmin"`
}

type UserClaims struct {
	User User `json:"user"`
	jwt.RegisteredClaims
}
