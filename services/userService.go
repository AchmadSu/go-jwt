package services

import (
	"example.com/m/models"
	"example.com/m/repositories"
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
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo}
}

func (s *userService) GetUser(id, email string) (models.User, *gorm.DB) {
	if id != "" {
		return s.repo.FindByID(id)
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
