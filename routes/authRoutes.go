package routes

import (
	"example.com/m/controllers"
	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(r *gin.Engine) {
	auth := r.Group("/")
	{
		auth.POST("/signup", controllers.SignUp)
		auth.POST("/login", controllers.Login)
	}
}
