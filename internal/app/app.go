package app

import (
	"fmt"

	"github.com/dj-ph3luy/go-playground/internal/controllers"
	"github.com/dj-ph3luy/go-playground/internal/models"
	"github.com/gin-gonic/gin"
)

func StartApp() {
	fmt.Println("starting app ...")
	models.ConnectDataBase()
	router := gin.Default()
	registerControllers(router);
	
	router.Run(":8090")
}

func registerControllers(r *gin.Engine) {
	controllers := []controllers.IController{
		&controllers.BasicController{},
		&controllers.UserController{},
	}
	for _,controller := range controllers {
		controller.RegisterRoutes(r)
	}
}
