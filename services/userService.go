package services

import (
	"example.com/m/models"
	"example.com/m/repositories"
	"example.com/m/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserService interface {
	GetUser(id, email string) (models.User, *gorm.DB)
	GetAllUsers(c *gin.Context) ([]models.User, *utils.Pagination, error)
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
