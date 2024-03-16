package middleware

import (
	"net/http"

	"github.com/dj-ph3luy/go-playground/internal/models"
	"github.com/gin-gonic/gin"
)

func BasicAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		username, password, ok := ctx.Request.BasicAuth()
		if !ok {
			ctx.String(http.StatusUnauthorized, "Authorization required")
			ctx.Abort()
			return
		}
		_, err := models.CheckLogin(username, password)
		if err != nil {
			ctx.String(http.StatusUnauthorized, "Unauthorized")
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}