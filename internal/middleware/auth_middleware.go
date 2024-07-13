package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/dj-ph3luy/go-playground/internal/services"
	"github.com/dj-ph3luy/go-playground/internal/views"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var UserService services.IUserService

var Secret = []byte("some_secret") //  TODO use viper for config

type CustomClaims struct {
	User views.User `json:"user"`
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
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization is required"})
		return
	}
	user, err := UserService.Login(ctx, username, password)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "username or password incorrect"})
		return
	}
	ctx.Set("user", user.ToView())
	ctx.Next()
}

func TokenAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenAuthProcedure(ctx)
	}
}

func tokenAuthProcedure(ctx *gin.Context) {
	authHeader := ctx.Request.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization is required"})
		return
	}

	token, err := verifyToken(authHeader[len("Bearer "):])
	switch {
	case token.Valid:
		claims, ok := token.Claims.(*CustomClaims)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "could not parse token"})
			return
		}
		ctx.Set("user", claims.User)
		ctx.Next()

	case errors.Is(err, jwt.ErrTokenMalformed):
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "could not parse token"})
		return

	case errors.Is(err, jwt.ErrTokenSignatureInvalid):
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "could not verify token"})
		return

	case errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet):
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token expired"})
		return

	default:
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "could not authorize"})
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
