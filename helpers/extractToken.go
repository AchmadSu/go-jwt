package helpers

import (
	"net/http"
	"strings"

	"example.com/m/errs"
	"github.com/gin-gonic/gin"
)

func ExtractToken(c *gin.Context) (string, error) {
	cookieToken, err := c.Cookie("Authorization")
	authHeader := c.GetHeader("Authorization")

	if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimPrefix(authHeader, "Bearer "), nil
	} else if err == nil {
		return cookieToken, nil
	}

	return "", errs.New("Token not found in Header and Cookies", http.StatusBadRequest)
}
