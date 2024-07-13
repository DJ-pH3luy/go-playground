package controllers

import (
	"net/http"

	"github.com/dj-ph3luy/go-playground/internal/dto"
	"github.com/dj-ph3luy/go-playground/internal/middleware"
	"github.com/dj-ph3luy/go-playground/internal/services"
	"github.com/dj-ph3luy/go-playground/internal/views"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	Service services.IUserService
}

// type RegisterUserInput struct {
// 	Username string `json:"username" binding:"required"`
// 	Email    string `json:"email" binding:"required"`
// 	Password string `json:"password" binding:"required"`
// }

// type UpdateUserPasswordInput struct {
// 	Password string `json:"password" binding:"required"`
// }

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
	user, err := c.Service.GetById(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "could get user", "error": err.Error()})
		return
	}
	if user.Name != loggedInUser(ctx).Name && !loggedInUser(ctx).IsAdmin {
		ctx.JSON(http.StatusForbidden, gin.H{"message": "forbidden"})
		return
	}
	ctx.JSON(http.StatusOK, user)
}

func (c *UserController) getUsers(ctx *gin.Context) {
	users, err := c.Service.GetAll(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "could not get users", "error": err.Error()})
		return
	}
	if !loggedInUser(ctx).IsAdmin {
		ctx.JSON(http.StatusForbidden, gin.H{"message": "forbidden"})
		return
	}
	ctx.JSON(http.StatusOK, users)
}

func (c *UserController) registerHandler(ctx *gin.Context) {
	var input dto.CreateUser

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "bad request", "error": err.Error()})
		return
	}

	id, err := c.Service.Create(ctx, input)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "bad request", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "registration success", "id": id})
}

func (c *UserController) updateUserPassword(ctx *gin.Context) {
	var input dto.UpdateUser

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "bad request", "error": err.Error()})
		return
	}

	input.Id = ctx.Param("id")
	id, err := c.Service.Update(ctx, input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "could not update user", "error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "updated user", "id": id})
}

// loggedInUser returns the currently logged in user from the context
func loggedInUser(ctx *gin.Context) views.User {
	user, ok := ctx.Get("user")
	if !ok {
		return views.User{}
	}
	return user.(views.User)
}
