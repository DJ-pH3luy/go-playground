package controllers

import (
	"net/http"

	"github.com/dj-ph3luy/go-playground/internal/middleware"
	"github.com/dj-ph3luy/go-playground/internal/models"
	"github.com/gin-gonic/gin"
)

type UserController struct {
}

type RegisterUserInput struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UpdateUserPasswordInput struct {
	Password string `json:"password" binding:"required"`
}

func (c *UserController) RegisterRoutes(router *gin.Engine) {
	userGroup := router.Group("v1/user")
	userGroup.POST("/register", c.registerHandler)

	userGroupProtected := userGroup.Group("/")
	userGroupProtected.Use(middleware.BasicAuthMiddleware())
	userGroupProtected.GET("/", c.getUsers)
	userGroupProtected.GET("/:id", c.getUser)
	userGroupProtected.PUT("/:id/password", c.updateUserPassword)
}

func (c *UserController) getUser(ctx *gin.Context) {
	id := ctx.Param("id")
	user, err := models.GetUser(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message":"could get user", "error": err.Error()})
		return
	}
	if user.Username != LoggedInUser(ctx).Username && !LoggedInUser(ctx).IsAdmin{
		ctx.JSON(http.StatusForbidden, gin.H{"message":"forbidden"})
		return
	}
	ctx.JSON(http.StatusOK, user)
}

func (c *UserController) getUsers(ctx *gin.Context) {
	users, err := models.GetUsers()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message":"could not get users", "error": err.Error()})
		return
	}
	if !LoggedInUser(ctx).IsAdmin{
		ctx.JSON(http.StatusForbidden, gin.H{"message":"forbidden"})
		return
	}
	ctx.JSON(http.StatusOK, users)
}

func (c *UserController) registerHandler(ctx *gin.Context) {
	var input RegisterUserInput

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message":"bad request", "error": err.Error()})
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

	ctx.JSON(http.StatusOK, gin.H{"message": "registration success", "user": user.MapToView()})
}

func (c *UserController) updateUserPassword(ctx *gin.Context) {
	var input UpdateUserPasswordInput

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message":"bad request","error": err.Error()})
		return
	}
	id := ctx.Param("id")
	user, err := models.GetUser(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message":"could not get user", "error": err.Error()})
		return
	}
	if user.Username != LoggedInUser(ctx).Username {
		ctx.JSON(http.StatusForbidden, gin.H{"message":"you can only change your own password"})
		return
	}

	updatedUser := models.User{
		Username: user.Username,
		Email: user.Email,
		Password: input.Password,
		IsAdmin: user.IsAdmin,
	}
	updatedUser.ID = user.Id
	err = updatedUser.Update()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message":"could not update user", "error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message":"updated user", "user": updatedUser.MapToView()})
}

func LoggedInUser(ctx *gin.Context) (models.UserViewModel) {
	user, ok := ctx.Get("user")
	if (!ok) {
		return models.UserViewModel{}
	}
	return user.(models.UserViewModel)
}


 