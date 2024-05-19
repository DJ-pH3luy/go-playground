package controllers

import (
	"fmt"
	"net/http"

	"github.com/dj-ph3luy/go-playground/internal/middleware"
	"github.com/gin-gonic/gin"
)

type BasicController struct {
	
}

func (c *BasicController) RegisterRoutes(router *gin.Engine) {
	basicGroup := router.Group("/v1/foo")
	basicGroup.Use(middleware.TokenAuthMiddleware())
	basicGroup.GET("/", c.getHandler);
	basicGroup.POST("/", c.postHandler);
}

func (c *BasicController) getHandler(ctx *gin.Context) {
	fmt.Println("GET basic test")
	ctx.JSON(http.StatusOK, gin.H{"user": LoggedInUser(ctx)})
}

func (c *BasicController) postHandler(ctx *gin.Context) {
	fmt.Println("POST basic test")
	ctx.String(http.StatusOK, fmt.Sprintln("POST OK"))
}