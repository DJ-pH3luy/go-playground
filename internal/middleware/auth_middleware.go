package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/dj-ph3luy/go-playground/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var Secret = []byte("some_secret") // TODO: config setup


type CustomClaims struct {
    User models.UserViewModel `json:"user"`
    jwt.RegisteredClaims
}

func TokenOrBasicAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if _, _, ok := ctx.Request.BasicAuth(); ok {
			basicAuthProcedure(ctx)
		} else {
			tokenAuthProcedure(ctx)
		}
	}
}

func BasicAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		basicAuthProcedure(ctx)
	}
}

func basicAuthProcedure(ctx *gin.Context) {
	username, password, ok := ctx.Request.BasicAuth()
	if !ok {
		ctx.String(http.StatusUnauthorized, "Authorization header required")
		ctx.Abort()
		return
	}
	user, err := models.CheckLogin(username, password)
	if err != nil {
		ctx.String(http.StatusUnauthorized, "Password or username incorrect")
		ctx.Abort()
		return
	}
	ctx.Set("user", user)
	ctx.Next()
}

func TokenAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenAuthProcedure(ctx)
	}
}

func tokenAuthProcedure(ctx *gin.Context) {
	authHeader := ctx.Request.Header.Get("Authorization");
	if authHeader == "" {
		ctx.String(http.StatusUnauthorized, "Authorization header is required")
		ctx.Abort()
		return
	}
	if !strings.HasPrefix(authHeader, "Bearer ") {
		ctx.String(http.StatusUnauthorized, "Bearer token is required")
		ctx.Abort()
		return
	}

	token, err := verifyToken(authHeader[len("Bearer "):])
	switch {
		case token.Valid:
			claims, ok := token.Claims.(*CustomClaims)
			if !ok {
				ctx.String(http.StatusUnauthorized, "Could not parse claims")
				ctx.Abort()
				return
			}
			ctx.Set("user", claims.User)
			ctx.Next()

		case errors.Is(err, jwt.ErrTokenMalformed):
			ctx.String(http.StatusUnauthorized, "Could not parse token")
			ctx.Abort()
			return

		case errors.Is(err, jwt.ErrTokenSignatureInvalid):
			ctx.String(http.StatusUnauthorized, "Could not verify token")
			ctx.Abort()
			return

		case errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet):
			ctx.String(http.StatusUnauthorized, "Token expired")
			ctx.Abort()
			return

		default:
			ctx.String(http.StatusUnauthorized, "Could not authorize")
			ctx.Abort()
			return
	}
}

func verifyToken(tokenString string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return Secret, nil
	})
 }


