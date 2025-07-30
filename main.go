package main

import (
	"example.com/m/controllers"
	"example.com/m/initializers"
	"example.com/m/middleware"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()
	initializers.SyncDatabase()
}

func main() {
	// fmt.Println("Hello")
	r := gin.Default()
	// r.GET("/ping", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{
	// 		"message": "pong",
	// 	})
	// })
	r.POST("/signup", controllers.SignUp)
	r.POST("/auth", controllers.Login)
	r.POST("/logout", middleware.RequireAuth, controllers.Logout)
	r.GET("/validate", middleware.RequireAuth, controllers.Validate)
	r.GET("/users", middleware.RequireAuth, controllers.GetUsers)
	r.Run()
}
