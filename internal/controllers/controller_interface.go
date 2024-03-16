package controllers

import "github.com/gin-gonic/gin"

type IController interface {
	RegisterRoutes(r *gin.Engine)
}

