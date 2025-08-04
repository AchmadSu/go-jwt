package middleware

import (
	"example.com/m/helpers"
	"example.com/m/repositories"
	"example.com/m/services/token"
	"example.com/m/utils"
	"github.com/gin-gonic/gin"
)

func RequireAuth(c *gin.Context) {

	messageError := "Unauthorized. You have no permission to access this endpoint!"
	resp := utils.NewResponse()
	tokenString, err := helpers.ExtractToken(c)
	if err != nil {
		errResp := utils.PrintErrorResponse(resp, err, messageError)
		errResp.Send(c)
		c.AbortWithStatus(utils.GetSafeStatusCode(errResp.StatusCode))
		return
	}
	token := token.NewJwtTokenService(repositories.NewUserRepository())
	user, err := token.ValidateToken(tokenString)
	if err != nil {
		errResp := utils.PrintErrorResponse(resp, err, messageError)
		errResp.Send(c)
		c.AbortWithStatus(utils.GetSafeStatusCode(errResp.StatusCode))
		return
	}

	c.Set("user", utils.ToPublicUser(user))
	c.Next()
}
