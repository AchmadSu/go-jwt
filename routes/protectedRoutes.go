package routes

import (
	"example.com/m/controllers"
	"example.com/m/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterProtectedRoutes(r *gin.Engine) {
	protected := r.Group("/")
	protected.Use(middleware.RequireAuth)
	{
		//users
		protected.GET("/users", controllers.GetUsers)
		protected.POST("/logout", controllers.Logout)
		protected.GET("/products", controllers.GetProducts)
		protected.POST("/products", controllers.CreateProduct)
	}
}
