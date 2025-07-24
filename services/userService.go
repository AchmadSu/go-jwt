package services

import (
	"net/http"
	"strconv"

	"example.com/m/models"
	"example.com/m/repositories"
	"example.com/m/services/token"
	"example.com/m/services/validator"
	"example.com/m/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService interface {
	GetUser(id, email string) (models.User, *gorm.DB)
	GetAllUsers(c *gin.Context) ([]models.User, *utils.Pagination, error)
	Register(models.CreateUserInput) (models.User, error)
	Login(*gin.Context, models.LoginUserInput) (models.User, string, error)
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo}
}

func (s *userService) GetUser(id, email string) (models.User, *gorm.DB) {
	if id != "" {
		parsedID, err := strconv.Atoi(id)
		if err != nil {
			return models.User{}, &gorm.DB{Error: err}
		}
		return s.repo.FindByID(parsedID)
	}
	return s.repo.FindByEmail(email)
}

func (s *userService) GetAllUsers(c *gin.Context) ([]models.User, *utils.Pagination, error) {
	return s.repo.FindAll(c)
}

func (s *userService) Register(input models.CreateUserInput) (models.User, error) {
	validator := validator.NewUserValidatorService(s.repo)
	isValid, err := validator.ValidateUserRegister(input.Email)
	if !isValid {
		return models.User{}, err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), 10)
	if err != nil {
		return models.User{}, err
	}
	input.Password = string(hash)
	return s.repo.Create(input)
}

func (s *userService) Login(c *gin.Context, input models.LoginUserInput) (models.User, string, error) {
	validator := validator.NewUserValidatorService(s.repo)
	user, err := validator.ValidateUserLogin(input)
	if err != nil {
		return models.User{}, "", err
	}
	token := token.NewJwtTokenService(s.repo)
	tokenString, exp, err := token.CreateToken(int(user.ID))
	if err != nil {
		return models.User{}, "", err
	}
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, exp, "", "", false, true)
	return user, tokenString, nil
}
