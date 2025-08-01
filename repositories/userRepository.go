package repositories

import (
	"example.com/m/dto"
	"example.com/m/initializers"
	"example.com/m/models"
	"example.com/m/utils"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindByID(id int) (models.User, *gorm.DB)
	FindByEmail(email string) (models.User, *gorm.DB)
	FindAll(paginate *dto.PaginationRequest) (*dto.PaginationResponse[dto.PublicUser], error)
	Create(input *dto.CreateUserInput) (dto.PublicUser, error)
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

func (r *userRepository) FindAll(request *dto.PaginationRequest) (*dto.PaginationResponse[dto.PublicUser], error) {
	query := initializers.DB.Model(&models.User{}).
		Select("id, email", "created_at", "updated_at")
	allowedSortFields := []string{"id", "email", "created_at"}
	searchFields := []string{"email"}
	defaultOrder := "created_at desc"
	return utils.Paginate[dto.PublicUser](request, query, allowedSortFields, defaultOrder, searchFields)
}

func (r *userRepository) Create(input *dto.CreateUserInput) (dto.PublicUser, error) {
	user := models.User{
		Email:    input.Email,
		Password: input.Password,
	}

	result := initializers.DB.Create(&user)
	return utils.ToPublicUser(user), result.Error
}
