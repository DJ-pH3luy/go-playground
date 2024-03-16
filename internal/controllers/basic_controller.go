package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BasicController struct {
	
}

func (c *BasicController) RegisterRoutes(router *gin.Engine) {
	router.GET("/v1/foo", c.getHandler);
	router.POST("/v1/foo", c.postHandler);
}

func (c *BasicController) getHandler(ctx *gin.Context) {
	fmt.Println("GET basic test")
	ctx.String(http.StatusOK, fmt.Sprintln("GET OK"))
}

func (c *BasicController) postHandler(ctx *gin.Context) {
	fmt.Println("POST basic test")
	ctx.String(http.StatusOK, fmt.Sprintln("POST OK"))
}