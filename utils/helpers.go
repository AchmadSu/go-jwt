package utils

import (
	"log"
	"net/http"

	"example.com/m/dto"
	"example.com/m/errs"
	"github.com/gin-gonic/gin"
)

func GetUserFromContext(c *gin.Context) (dto.PublicUser, error) {
	userAny, exist := c.Get("user")
	if !exist {
		return dto.PublicUser{}, errs.New("user context not found", http.StatusNotFound)
	}

	user, ok := userAny.(dto.PublicUser)
	if !ok {
		return dto.PublicUser{}, errs.New("invalid user context type", http.StatusInternalServerError)
	}
	return user, nil
}

func PrintErrorResponse(resp *Response, err error, message string) *Response {
	log.Printf("[ERROR] %v", err)
	if httpErr, ok := err.(*errs.HTTPError); ok {
		resp.SetStatus(httpErr.StatusCode).
			SetMessage(message).
			SetError(httpErr.Message)
		return resp
	}
	resp.SetStatus(http.StatusInternalServerError).
		SetMessage(message).
		SetError(GetSafeErrorMessage(err, "Unknown error occurred"))
	return resp
}

func ContainsString(list []string, val string) bool {
	for _, v := range list {
		if v == val {
			return true
		}
	}
	return false
}
