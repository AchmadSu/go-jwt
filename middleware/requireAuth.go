package middleware

import (
	"net/http"

	"example.com/m/errs"
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
		if httpErr, ok := err.(*errs.HTTPError); ok {
			resp.SetStatus(httpErr.StatusCode).
				SetMessage(messageError).
				SetError(httpErr.Message).
				Send(c)
			c.AbortWithStatus(httpErr.StatusCode)
			return
		}
		resp.SetStatus(http.StatusInternalServerError).
			SetMessage(messageError).
			SetError(utils.GetSafeErrorMessage(err, "Unknown error occurred")).
			Send(c)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	token := token.NewJwtTokenService(repositories.NewUserRepository())
	user, err := token.ValidateToken(tokenString)
	if err != nil {
		if httpErr, ok := err.(*errs.HTTPError); ok {
			resp.SetStatus(httpErr.StatusCode).
				SetMessage(messageError).
				SetError(httpErr.Message).
				Send(c)
			c.AbortWithStatus(httpErr.StatusCode)
			return
		}
		resp.SetStatus(http.StatusInternalServerError).
			SetMessage(messageError).
			SetError(utils.GetSafeErrorMessage(err, "Unknown error occurred")).
			Send(c)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Set("user", utils.ToPublicUser(user))
	c.Next()
}
