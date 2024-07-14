package controllers

import (
	"net/http"

	"github.com/dj-ph3luy/go-playground/internal/services"
	"github.com/gin-gonic/gin"
)

type LoginController struct {
	Service services.IUserService
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

	user, err := c.Service.Login(ctx, input.Username, input.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "login failed", "error": "password or username incorrect"})
		return
	}

	tokenString, err := c.Service.GenerateJWT(user.ToView())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "login failed", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": tokenString})
}
