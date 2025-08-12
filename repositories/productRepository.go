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

const ProductTable config.TableName = "products"

type ProductRepository interface {
	FindByProductID(id int) (models.Product, *gorm.DB)
	FindByProductCode(code string) (models.Product, *gorm.DB)
	FindByProductName(name string) (models.Product, *gorm.DB)
	FindAllProducts(paginate *dto.PaginationRequest) (*dto.PaginationResponse[dto.PublicProduct], error)
	CreateProduct(input *dto.CreateProductInput, creatorId uint) (dto.PublicProduct, error)
	UpdateProduct(id int, input *dto.UpdateProductInput, modifierId uint) (dto.PublicProduct, error)
}

type productRepository struct{}

func NewProductRepository() ProductRepository {
	return &productRepository{}
}

func (r *productRepository) FindByProductID(id int) (models.Product, *gorm.DB) {
	var product models.Product
	result := initializers.DB.First(&product, "id = ?", id)
	var productWithUser models.Product
	err := helpers.PreloadRelationByID(&productWithUser, product.ID, []string{"Creator", "Modifier"})
	if err != nil {
		return models.Product{}, &gorm.DB{Error: err}
	}

	return productWithUser, result
}

func (r *productRepository) FindByProductCode(code string) (models.Product, *gorm.DB) {
	var product models.Product
	result := initializers.DB.First(&product, "code = ?", code)
	var productWithUser models.Product
	err := helpers.PreloadRelationByID(&productWithUser, product.ID, []string{"Creator", "Modifier"})
	if err != nil {
		return models.Product{}, &gorm.DB{Error: err}
	}

	return productWithUser, result
}

func (r *productRepository) FindByProductName(name string) (models.Product, *gorm.DB) {
	var product models.Product
	result := initializers.DB.First(&product, "name = ?", name)
	var productWithUser models.Product
	err := helpers.PreloadRelationByID(&productWithUser, product.ID, []string{"Creator", "Modifier"})
	if err != nil {
		return models.Product{}, &gorm.DB{Error: err}
	}

	return productWithUser, result
}

func (r *productRepository) FindAllProducts(request *dto.PaginationRequest) (*dto.PaginationResponse[dto.PublicProduct], error) {
	query := initializers.DB.Model(&models.Product{}).
		Joins("LEFT JOIN users AS creator ON creator.id = products.created_by").
		Joins("LEFT JOIN users AS modifier ON modifier.id = products.modified_by").
		Select([]string{
			"products.id AS id",
			"products.code",
			"products.name AS name",
			"products.desc AS desc",
			"products.is_active",
			"products.created_by AS creator_id",
			"products.modified_by AS modifier_id",
			"products.created_at",
			"products.updated_at",
			"creator.name AS creator_name",
			"modifier.name AS modifier_name",
		})

	query = utils.FilterQuery(request, query, string(ProductTable)).Debug()

	allowedSortFields := []string{
		`id`,
		`name`,
		`code`,
		`is_active`,
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
	pageResult, err := utils.Paginate[dto.PublicProduct](request, query, allowedSortFields, defaultOrder, searchFields)

	if err != nil {
		return nil, err
	}

	helpers.SetEntityStatusLabel(pageResult.Data,
		func(item *dto.PublicProduct) int {
			return int(item.IsActive)
		},
		func(item *dto.PublicProduct, label string) {
			item.Status = label
		})

	return pageResult, nil
}

func (r *productRepository) CreateProduct(input *dto.CreateProductInput, creatorId uint) (dto.PublicProduct, error) {
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
	err := helpers.PreloadRelationByID(&productWithUser, product.ID, []string{"Creator", "Modifier"})
	if err != nil {
		return dto.PublicProduct{}, err
	}

	return utils.ToPublicProduct(productWithUser), nil
}

func (r *productRepository) UpdateProduct(id int, input *dto.UpdateProductInput, modifierId uint) (dto.PublicProduct, error) {
	var product models.Product
	trx := initializers.DB.Begin()
	if trx.Error != nil {
		trx.Rollback()
		return dto.PublicProduct{}, trx.Error
	}

	if err := trx.First(&product, id).Error; err != nil {
		trx.Rollback()
		return dto.PublicProduct{}, err
	}

	data := map[string]any{
		"Code":       input.Code,
		"Name":       input.Name,
		"Desc":       input.Desc,
		"IsActive":   *input.IsActive,
		"ModifiedBy": modifierId,
	}

	if err := utils.AssignedKeyModel(&product, data); err != nil {
		trx.Rollback()
		return dto.PublicProduct{}, err
	}

	if err := trx.Save(&product).Error; err != nil {
		trx.Rollback()
		return dto.PublicProduct{}, err
	}

	if err := trx.Commit().Error; err != nil {
		return dto.PublicProduct{}, err
	}

	var productWithUser models.Product
	err := helpers.PreloadRelationByID(&productWithUser, product.ID, []string{"Creator", "Modifier"})
	if err != nil {
		return dto.PublicProduct{}, err
	}

	return utils.ToPublicProduct(productWithUser), nil
}
