package repositories

import (
	"example.com/m/config"
	"example.com/m/dto"
	"example.com/m/helpers"
	"example.com/m/initializers"
	"example.com/m/models"
	"example.com/m/utils"
	"gorm.io/gorm"
)

const UserTable config.TableName = "users"

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
	result := initializers.DB.First(&user, "id = ?", id).Debug()
	return user, result
}

func (r *userRepository) FindByEmail(email string) (models.User, *gorm.DB) {
	var user models.User
	result := initializers.DB.Find(&user, "email = ?", email).Debug()
	return user, result
}

func (r *userRepository) FindAll(request *dto.PaginationRequest) (*dto.PaginationResponse[dto.PublicUser], error) {
	query := initializers.DB.Model(&models.User{}).
		Select("id", "name", "email", "is_active", "created_at", "updated_at")
	query = utils.FilterQuery(request, query, string(UserTable)).Debug()
	allowedSortFields := []string{"id", "name", "email", "is_active", "created_at", "updated_at"}
	searchFields := []string{"name", "email"}
	defaultOrder := "created_at desc"
	pageResult, err := utils.Paginate[dto.PublicUser](request, query, allowedSortFields, defaultOrder, searchFields)
	if err != nil {
		return nil, err
	}
	helpers.SetEntityStatusLabel(pageResult.Data,
		func(item *dto.PublicUser) int {
			return int(item.IsActive)
		},
		func(item *dto.PublicUser, label string) {
			item.Status = label
		})
	return pageResult, nil
}

func (r *userRepository) Create(input *dto.CreateUserInput) (dto.PublicUser, error) {
	user := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
	}

	result := initializers.DB.Create(&user).Debug()
	return utils.ToPublicUser(user), result.Error
}
