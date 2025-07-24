package controllers

import (
	"net/http"

	"example.com/m/errs"
	"example.com/m/models"
	"example.com/m/repositories"
	"example.com/m/services"
	"example.com/m/utils"
	"github.com/gin-gonic/gin"
)

var userService = services.NewUserService(repositories.NewUserRepository())

func SignUp(c *gin.Context) {
	var input models.CreateUserInput
	resp := utils.NewResponse()
	message := "Failed to create user"
	if c.Bind(&input) != nil {
		resp.SetStatus(http.StatusBadRequest).
			SetMessage(message).
			SetError("Failed to read body").
			Send(c)
		return
	}
	user, err := userService.Register(input)
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
		SetPayload(utils.ToPublicUser(user)).
		Send(c)
}

func Login(c *gin.Context) {
	var input models.LoginUserInput
	resp := utils.NewResponse()
	message := "Unauthorized. Failed to login"

	if c.Bind(&input) != nil {
		resp.SetStatus(http.StatusBadRequest).
			SetMessage(message).
			SetError("Failed to read body").
			Send(c)
		return
	}

	user, tokenString, err := userService.Login(c, input)
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

	message = "Login successfully"
	resp.SetStatus(http.StatusOK).
		SetMessage(message).
		SetPayload(utils.LoginResponse{
			User:  utils.ToPublicUser(user),
			Token: tokenString,
		}).
		Send(c)
}

func GetUsers(c *gin.Context) {
	resp := utils.NewResponse()
	id := c.Query("id")
	email := c.Query("email")

	if email != "" || id != "" {
		user, result := userService.GetUser(id, email)

		if result.Error != nil {
			resp.SetStatus(http.StatusInternalServerError).
				SetMessage("Failed to get user data").
				SetError(utils.GetSafeErrorMessage(result.Error, "Unknown error occurred")).
				Send(c)
			return
		}

		if result.RowsAffected == 0 {
			resp.SetStatus(http.StatusNotFound).
				SetMessage("User not found").
				Send(c)
			return
		}

		publicUser := utils.ToPublicUser(user)
		resp.SetMessage("Get user successfully").
			SetPayload(publicUser).
			Send(c)
		return
	}

	users, pg, err := userService.GetAllUsers(c)
	if err != nil {
		resp.SetStatus(http.StatusInternalServerError).
			SetMessage("Failed to fetch user data").
			SetError(utils.GetSafeErrorMessage(err, "Unknown error occurred")).
			Send(c)
		return
	}

	publicUsers := make([]models.PublicUser, 0, len(users))
	for _, u := range users {
		publicUsers = append(publicUsers, utils.ToPublicUser(u))
	}

	resp.SetMessage("Get users successfully").
		SetPayload(publicUsers).
		SetMeta(gin.H{
			"total":        pg.Total,
			"current_page": pg.Page,
			"total_pages":  pg.TotalPages,
		}).
		Send(c)
}

func Validate(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "I'm logged in",
	})
}
