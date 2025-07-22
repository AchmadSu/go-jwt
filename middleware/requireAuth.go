package middleware

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"example.com/m/initializers"
	"example.com/m/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func RequireAuth(c *gin.Context) {

	tokenString, err := c.Cookie("Authorization")
	messageError := "Unauthorized. You have no permission to access this endpoint!"

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": messageError,
		})
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(os.Getenv("SECRET_API_KEY")), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		expUnix, ok := claims["exp"].(float64)
		if !ok {
			c.JSON(http.StatusNotFound, gin.H{
				"error": messageError,
			})
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		tokenExpired := float64(time.Now().Unix()) > expUnix
		if tokenExpired {
			c.JSON(http.StatusNotFound, gin.H{
				"error": messageError,
			})
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		var user models.User
		initializers.DB.First(&user, claims["sub"])

		if user.ID == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"error": messageError,
			})
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		c.Set("user", user)
		c.Next()

	} else {
		c.JSON(http.StatusNotFound, gin.H{
			"error": messageError,
		})
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}
