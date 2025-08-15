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
		// users
		protected.GET("/users", controllers.GetUsers)
		protected.POST("/logout", controllers.Logout)
		// products
		protected.GET("/products", controllers.GetProducts)
		protected.POST("/products", controllers.CreateProduct)
		protected.PUT("/products", controllers.UpdateProduct)
		//stocks
		protected.GET("/stocks", controllers.GetStocks)
		protected.POST("/stocks", controllers.CreateStock)
		protected.PUT("/stocks", controllers.UpdateStock)
		//orders
		protected.GET("/orders", controllers.GetOrders)
		protected.POST("/orders", controllers.CreateOrder)
	}
}
