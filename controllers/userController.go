package controllers

import (
	"net/http"

	"example.com/m/dto"
	"example.com/m/errs"
	"example.com/m/helpers"
	"example.com/m/repositories"
	"example.com/m/services"
	"example.com/m/utils"
	"github.com/gin-gonic/gin"
)

var userService = services.NewUserService(repositories.NewUserRepository())

func SignUp(c *gin.Context) {
	var input dto.CreateUserInput
	resp := utils.NewResponse()
	message := "Failed to create user"
	if c.Bind(&input) != nil {
		resp.SetStatus(http.StatusBadRequest).
			SetMessage(message).
			SetError("Failed to read body").
			Send(c)
		return
	}
	user, err := userService.Register(&input)
	if err != nil {
		if httpErr, ok := err.(*errs.HTTPError); ok {
			resp.SetStatus(httpErr.StatusCode).
				SetMessage(message).
				SetError(httpErr.Message).
				Send(c)
			return
		}
		resp.SetStatus(http.StatusInternalServerError).
			SetMessage(message).
			SetError(utils.GetSafeErrorMessage(err, "Unknown error occurred")).
			Send(c)
		return
	}
	message = "User has registered successfully"
	resp.SetStatus(http.StatusOK).
		SetMessage(message).
		SetPayload(user).
		Send(c)
}

func Login(c *gin.Context) {
	var input dto.LoginUserInput
	resp := utils.NewResponse()
	message := "Unauthorized. Failed to login"

	if c.Bind(&input) != nil {
		resp.SetStatus(http.StatusBadRequest).
			SetMessage(message).
			SetError("Failed to read body").
			Send(c)
		return
	}

	result := userService.Login(&input)
	if result.Err != nil {
		if httpErr, ok := result.Err.(*errs.HTTPError); ok {
			resp.SetStatus(httpErr.StatusCode).
				SetMessage(message).
				SetError(httpErr.Message).
				Send(c)
			return
		}
		resp.SetStatus(http.StatusInternalServerError).
			SetMessage(message).
			SetError(utils.GetSafeErrorMessage(result.Err, "Unknown error occurred")).
			Send(c)
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", result.Token, result.Exp, "", "", false, true)

	message = "Login successfully"
	resp.SetStatus(http.StatusOK).
		SetMessage(message).
		SetPayload(utils.LoginResponse{
			User:  result.User,
			Token: result.Token,
		}).
		Send(c)
}

func GetUsers(c *gin.Context) {
	resp := utils.NewResponse()
	id := c.Query("id")
	email := c.Query("email")

	if email != "" || id != "" {
		user, err := userService.GetUser(id, email)

		if err != nil {
			if httpErr, ok := err.(*errs.HTTPError); ok {
				resp.SetStatus(httpErr.StatusCode).
					SetMessage("Failed to fetch user data").
					SetError(httpErr.Message).
					Send(c)
				return
			}
			resp.SetStatus(http.StatusInternalServerError).
				SetMessage("Failed to get user data").
				SetError(utils.GetSafeErrorMessage(err, "Unknown error occurred")).
				Send(c)
			return
		}

		resp.SetMessage("Get user by ID or Email successfully").
			SetPayload(user).
			Send(c)
		return
	}

	var pagination dto.PaginationRequest
	if err := c.ShouldBindQuery(&pagination); err != nil {
		resp.SetStatus(http.StatusBadRequest).
			SetMessage("Failed to get user data").
			SetError(utils.GetSafeErrorMessage(err, "Unknown error occurred")).
			Send(c)
		return
	}

	pg, err := userService.GetAllUsers(&pagination)
	if err != nil {
		if httpErr, ok := err.(*errs.HTTPError); ok {
			resp.SetStatus(httpErr.StatusCode).
				SetMessage("Failed to fetch user data").
				SetError(httpErr.Message).
				Send(c)
			return
		}
		resp.SetStatus(http.StatusInternalServerError).
			SetMessage("Failed to fetch user data").
			SetError(utils.GetSafeErrorMessage(err, "Unknown error occurred")).
			Send(c)
		return
	}

	resp.SetMessage("Get users successfully").
		SetPayload(pg.Data).
		SetMeta(gin.H{
			"page":      pg.Page,
			"totalPage": pg.TotalPages,
			"totalData": pg.Total,
		}).
		Send(c)
}

func Logout(c *gin.Context) {
	message := "Logout failed!"
	resp := utils.NewResponse()
	tokenString, _ := helpers.ExtractToken(c)
	err := userService.Logout(tokenString)
	if err != nil {
		if httpErr, ok := err.(*errs.HTTPError); ok {
			resp.SetStatus(httpErr.StatusCode).
				SetMessage(message).
				SetError(httpErr.Message).
				Send(c)
			return
		}
		resp.SetStatus(http.StatusInternalServerError).
			SetMessage(message).
			SetError(utils.GetSafeErrorMessage(err, "Unknown error occurred")).
			Send(c)
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", "", -1, "", "", false, true)
	c.Set("user", nil)

	message = "Logout successfuly!"
	resp.SetStatus(http.StatusOK).
		SetMessage(message).
		Send(c)
}
