package repositories

import (
	"example.com/m/initializers"
	"example.com/m/models"
	"example.com/m/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindByID(id int) (models.User, *gorm.DB)
	FindByEmail(email string) (models.User, *gorm.DB)
	FindAll(c *gin.Context) ([]models.User, *utils.Pagination, error)
	Create(models.CreateUserInput) (models.User, error)
}

type userRepository struct{}

func NewUserRepository() UserRepository {
	return &userRepository{}
}

func (r *userRepository) FindByID(id int) (models.User, *gorm.DB) {
	var user models.User
	result := initializers.DB.First(&user, "id = ?", id)
	return user, result
}

func (r *userRepository) FindByEmail(email string) (models.User, *gorm.DB) {
	var user models.User
	result := initializers.DB.Find(&user, "email = ?", email)
	return user, result
}

func (r *userRepository) FindAll(c *gin.Context) ([]models.User, *utils.Pagination, error) {
	pg, err := utils.Paginate(c, &models.User{})
	if err != nil {
		return nil, nil, err
	}

	var users []models.User
	result := initializers.DB.Limit(pg.Limit).Offset(pg.Offset).Find(&users)
	return users, pg, result.Error
}

func (r *userRepository) Create(input models.CreateUserInput) (models.User, error) {
	user := models.User{
		Email:    input.Email,
		Password: input.Password,
	}

	result := initializers.DB.Create(&user)
	return user, result.Error
}
