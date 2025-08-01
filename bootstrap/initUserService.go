package bootstrap

import (
	"example.com/m/repositories"
	"example.com/m/services"
	"example.com/m/services/token"
	"example.com/m/services/validator"
)

var UserService services.UserService

func InitUserService() {
	userRepo := repositories.NewUserRepository()
	userValidator := validator.NewUserValidatorService(userRepo)
	userToken := token.NewJwtTokenService(userRepo)
	UserService = services.NewUserService(userRepo, userValidator, userToken)
}
