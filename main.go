package main

import (
	"ticketon-auth-service/api/controllers"
	"ticketon-auth-service/api/middlewares"
	"ticketon-auth-service/api/repository"
	"github.com/gin-gonic/gin"
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
			apiUser.PUT("/:id", controllers.Ping).Use(middlewares.Auth())
		}

		accountApi := api.Group("/account")
		{
			accountApi.GET("", controllers.FindAccount).Use(middlewares.Auth())
		}

	}
	return router
}
