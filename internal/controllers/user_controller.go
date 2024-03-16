package controllers

import (
	"net/http"

	"github.com/dj-ph3luy/go-playground/internal/middleware"
	"github.com/dj-ph3luy/go-playground/internal/models"
	"github.com/gin-gonic/gin"
)

type UserController struct {
}

type RegisterInput struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (c *UserController) RegisterRoutes(router *gin.Engine) {
	userGroup := router.Group("v1/user")
	userGroup.POST("/register", c.registerHandler)

	userGroupProtected := userGroup.Group("/")
	userGroupProtected.Use(middleware.BasicAuthMiddleware())
	userGroupProtected.GET("/", c.getUsers)
	userGroupProtected.GET("/:id", c.getUser)
}

func (c *UserController) getUser(ctx *gin.Context) {
	id := ctx.Param("id")
	user, err := models.GetUser(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message":"could get user", "error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, user)
}

func (c *UserController) getUsers(ctx *gin.Context) {
	users, err := models.GetUsers()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message":"could get users", "error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, users)
}

func (c *UserController) registerHandler(ctx *gin.Context) {
	var input RegisterInput

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := models.User{}

	user.Username = input.Username
	user.Email = input.Email
	user.Password = input.Password

	err := user.Save()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message":"could not save user", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "registration success", "id": user.ID})
}
