package main

import (
	"example.com/m/initializers"
	"example.com/m/routes"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()
	initializers.SyncDatabase()
}

func main() {
	r := gin.Default()
	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"statusCode": 404, "message": "Endpoint not found"})
	})
	routes.RegisterAuthRoutes(r)
	routes.RegisterProtectedRoutes(r)
	r.Run()
}
