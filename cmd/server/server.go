package main

import (
	"fmt"

	"github.com/dj-ph3luy/go-playground/internal/controllers"
	"github.com/dj-ph3luy/go-playground/internal/db"
	"github.com/dj-ph3luy/go-playground/internal/entities"
	"github.com/dj-ph3luy/go-playground/internal/services"
	"github.com/dj-ph3luy/go-playground/internal/services/user"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var userService services.IUserService 

func main(){
	database := db.Connect(db.Config{
		ConnectionString: "postgresql://postgres:pw123@localhost:5432/postgres?sslmode=disable",
	})
	db.Migrate(database)
	initServices(database)

	fmt.Println("starting app ...")
	router := gin.Default()
	registerControllers(router);
	
	router.Run(":8090")
}

func registerControllers(r *gin.Engine) {
	controllers := []controllers.IController{
		&controllers.UserController{
			Service: userService,
		},
		&controllers.LoginController{
			Service: userService,
		},
	}
	for _,controller := range controllers {
		controller.RegisterRoutes(r)
	}
}

func initServices(database *gorm.DB) {
	userService = user.New(db.New[entities.User](database))
}