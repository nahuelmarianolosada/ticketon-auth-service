package main

import (
	"github.com/gin-gonic/gin"
	"ticketon-auth-service/api/controllers"
	"ticketon-auth-service/api/middlewares/auth"
	"ticketon-auth-service/api/repository"
)

func main() {
	// Initialize Database
	repository.Connect()
	repository.Migrate()
	// Initialize Router
	router := initRouter()
	router.Run(":8080")
}

func initRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/ping", controllers.Ping)
	api := router.Group("/api")
	{
		api.POST("/login", controllers.GenerateToken)

		apiUser := api.Group("/users")
		{
			apiUser.POST("", controllers.RegisterUser)
			apiUser.PUT("/:id", auth.AuthMiddleware(), controllers.UpdateUser)
		}

		accountApi := api.Group("/account")
		{
			accountApi.GET("", auth.AuthMiddleware(), controllers.FindAccount)
		}

	}
	return router
}
