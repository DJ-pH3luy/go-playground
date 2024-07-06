package controllers

import (
	"net/http"
	"time"

	"github.com/dj-ph3luy/go-playground/internal/middleware"
	"github.com/dj-ph3luy/go-playground/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type LoginController struct {
}

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (c *LoginController) RegisterRoutes(router *gin.Engine) {
	loginGroup := router.Group("v1/login")
	loginGroup.POST("/", c.loginHandler)
}

func (c *LoginController) loginHandler(ctx *gin.Context) {
	var input LoginInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "bad request", "error": err.Error()})
		return
	}

	user, err := models.CheckLogin(input.Username, input.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "login failed", "error": "password or username incorrect"})
		return
	}

	tokenString, err := generateJWT(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "login failed", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func generateJWT(user models.UserViewModel) (string, error) {
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &middleware.CustomClaims{
		User: user,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(middleware.Secret)
}
