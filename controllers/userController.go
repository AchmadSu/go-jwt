package controllers

import (
	"net/http"

	"example.com/m/bootstrap"
	"example.com/m/dto"
	"example.com/m/errs"
	"example.com/m/helpers"
	"example.com/m/utils"
	"github.com/gin-gonic/gin"
)

func SignUp(c *gin.Context) {
	var input dto.CreateUserInput
	resp := utils.NewResponse()
	message := "Failed to create user"
	if c.Bind(&input) != nil {
		err := errs.New("Body request invalid", http.StatusBadRequest)
		errResp := utils.PrintErrorResponse(resp, err, message)
		errResp.Send(c)
		return
	}
	user, err := bootstrap.UserService.Register(&input)
	if err != nil {
		errResp := utils.PrintErrorResponse(resp, err, message)
		errResp.Send(c)
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
		err := errs.New("Body request invalid", http.StatusBadRequest)
		errResp := utils.PrintErrorResponse(resp, err, message)
		errResp.Send(c)
		return
	}

	result := bootstrap.UserService.Login(&input)
	if result.Err != nil {
		errResp := utils.PrintErrorResponse(resp, result.Err, message)
		errResp.Send(c)
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
	var pagination dto.PaginationRequest
	resp := utils.NewResponse()
	message := "Failed to fetch user data"

	if err := c.ShouldBindQuery(&pagination); err != nil {
		err = errs.New(utils.GetSafeErrorMessage(err, "Body request pagination invalid"), http.StatusBadRequest)
		errResp := utils.PrintErrorResponse(resp, err, message)
		errResp.Send(c)
		return
	}

	if pagination.ID != nil || pagination.Email != "" {
		user, err := bootstrap.UserService.GetUser(&pagination)
		if err != nil {
			errResp := utils.PrintErrorResponse(resp, err, message)
			errResp.Send(c)
			return
		}
		message = "Get user by ID or Email successfully"
		resp.SetMessage(message).
			SetPayload(user).
			Send(c)
		return
	}

	pg, err := bootstrap.UserService.GetAllUsers(&pagination)
	if err != nil {
		errResp := utils.PrintErrorResponse(resp, err, message)
		errResp.Send(c)
		return
	}

	message = "Get users successfully"
	resp.SetMessage(message).
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
	err := bootstrap.UserService.Logout(tokenString)
	if err != nil {
		errResp := utils.PrintErrorResponse(resp, err, message)
		errResp.Send(c)
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
