package repositories

import (
	"example.com/m/dto"
	"example.com/m/initializers"
	"example.com/m/models"
	"example.com/m/utils"
	"gorm.io/gorm"
)

type ProductRepository interface {
	FindByID(id int) (models.Product, *gorm.DB)
	FindByCode(code string) (models.Product, *gorm.DB)
	FindByName(name string) (models.Product, *gorm.DB)
	FindAll(paginate *dto.PaginationRequest, creatorId uint, modifierId uint) (*dto.PaginationResponse[dto.PublicProduct], error)
	Create(input *dto.CreateProductInput, creatorId uint) (dto.PublicProduct, error)
}

type productRepository struct{}

func NewProductRepository() ProductRepository {
	return &productRepository{}
}

func (r *productRepository) FindByID(id int) (models.Product, *gorm.DB) {
	var product models.Product
	result := initializers.DB.First(&product, "id = ?", id)
	return product, result
}

func (r *productRepository) FindByCode(code string) (models.Product, *gorm.DB) {
	var product models.Product
	result := initializers.DB.First(&product, "code = ?", code)
	return product, result
}

func (r *productRepository) FindByName(name string) (models.Product, *gorm.DB) {
	var product models.Product
	result := initializers.DB.First(&product, "name = ?", name)
	return product, result
}

func (r *productRepository) FindAll(request *dto.PaginationRequest, creatorId uint, modifierId uint) (*dto.PaginationResponse[dto.PublicProduct], error) {
	query := initializers.DB.Model(&models.Product{}).
		Joins("LEFT JOIN users AS creator ON creator.id = products.created_by").
		Joins("LEFT JOIN users AS modifier ON modifier.id = products.modified_by").
		Select([]string{
			"products.id AS id",
			"products.code",
			"products.name AS name",
			"products.desc AS description",
			"products.created_by AS creator_id",
			"products.modified_by AS modifier_id",
			"products.created_at",
			"products.updated_at",
			"creator.name AS creator_name",
			"modifier.name AS modifer_name",
		})
	if creatorId > 0 {
		query = query.Where("products.created_by = ?", creatorId)
	} else if modifierId > 0 {
		query = query.Where("products.modified_by = ?", modifierId)
	}
	allowedSortFields := []string{
		`id`,
		`name`,
		`code`,
		`creator_name`,
		`modifier_name`,
		`created_at`,
		`updated_at`,
	}
	searchFields := []string{
		`products.name`,
		`products.code`,
		`products.desc`,
		`creator.name`,
		`modifier.name`,
	}
	defaultOrder := "products.name asc"
	return utils.Paginate[dto.PublicProduct](request, query, allowedSortFields, defaultOrder, searchFields)
}

func (r *productRepository) Create(input *dto.CreateProductInput, creatorId uint) (dto.PublicProduct, error) {
	product := models.Product{
		Code:      input.Code,
		Name:      input.Name,
		Desc:      input.Desc,
		CreatedBy: &creatorId,
	}
	result := initializers.DB.Create(&product)
	if result.Error != nil {
		return dto.PublicProduct{}, result.Error
	}

	var productWithUser models.Product
	err := initializers.DB.
		Preload("Creator").
		Preload("Modifier").
		First(&productWithUser, product.ID).Error
	if err != nil {
		return dto.PublicProduct{}, err
	}

	return utils.ToPublicProduct(productWithUser), nil
}
