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

const StockTable config.TableName = "stocks"

type StockRepository interface {
	FindByStockID(id int) (models.Stock, *gorm.DB)
	FindAllStocks(paginate *dto.PaginationRequest, stockRequest *dto.PaginationStockRequest) (*dto.PaginationResponse[dto.PublicStock], error)
	CreateStock(input *dto.CreateStockInput, creatorId uint) (dto.PublicStock, error)
	UpdateStock(id int, input *dto.UpdateStockInput, modifierId uint) (dto.PublicStock, error)
}

type stockRepository struct{}

func NewStockRepository() StockRepository {
	return &stockRepository{}
}

func (r *stockRepository) FindByStockID(id int) (models.Stock, *gorm.DB) {
	var stock models.Stock
	result := initializers.DB.First(&stock, "id = ?", id)
	var stockWithUser models.Stock
	err := helpers.PreloadRelationByID(&stockWithUser, stock.ID, []string{"Creator", "Modifier"})
	if err != nil {
		return models.Stock{}, &gorm.DB{Error: err}
	}

	return stockWithUser, result
}

func (r *stockRepository) FindAllStocks(request *dto.PaginationRequest, stockRequest *dto.PaginationStockRequest) (*dto.PaginationResponse[dto.PublicStock], error) {
	query := initializers.DB.Model(&models.Stock{}).
		Joins("LEFT JOIN products AS product ON product.id = stocks.product_id").
		Joins("LEFT JOIN users AS creator ON creator.id = stocks.created_by").
		Joins("LEFT JOIN users AS modifier ON modifier.id = stocks.modified_by").
		Select([]string{
			"stocks.id AS id",
			"product.code",
			"product.name AS name",
			"stocks.qty",
			"stocks.price",
			"stocks.date_entry",
			"stocks.is_active",
			"stocks.created_by AS creator_id",
			"stocks.modified_by AS modifier_id",
			"stocks.created_at",
			"stocks.updated_at",
			"creator.name AS creator_name",
			"modifier.name AS modifier_name",
		})

	query = utils.StockFilterQuery(stockRequest, query)
	query = utils.FilterQuery(request, query, string(StockTable)).Debug()

	allowedSortFields := []string{
		`id`,
		`name`,
		`code`,
		`is_active`,
		`qty`,
		`price`,
		`date_entry`,
		`creator_name`,
		`modifier_name`,
		`created_at`,
		`updated_at`,
	}
	searchFields := []string{
		`products.name`,
		`products.code`,
		`creator.name`,
		`modifier.name`,
	}
	defaultOrder := "products.name asc"
	pageResult, err := utils.Paginate[dto.PublicStock](request, query, allowedSortFields, defaultOrder, searchFields)

	if err != nil {
		return nil, err
	}

	helpers.SetEntityStatusLabel(pageResult.Data,
		func(item *dto.PublicStock) int {
			return int(item.IsActive)
		},
		func(item *dto.PublicStock, label string) {
			item.Status = label
		})

	return pageResult, nil
}

func (r *stockRepository) CreateStock(input *dto.CreateStockInput, creatorId uint) (dto.PublicStock, error) {
	stock := models.Stock{
		ProductId: &input.ProductId,
		Qty:       models.StockQty(input.Qty),
		Price:     models.StockPrice(input.Price),
		DateEntry: input.DateEntry,
		CreatedBy: &creatorId,
	}
	result := initializers.DB.Create(&stock)
	if result.Error != nil {
		return dto.PublicStock{}, result.Error
	}

	var stockWithUser models.Stock
	err := helpers.PreloadRelationByID(&stockWithUser, stock.ID, []string{"Creator", "Modifier"})
	if err != nil {
		return dto.PublicStock{}, err
	}

	return utils.ToPublicStock(stockWithUser), nil
}

func (r *stockRepository) UpdateStock(id int, input *dto.UpdateStockInput, modifierId uint) (dto.PublicStock, error) {
	var stock models.Stock

	trx := initializers.DB.Begin()
	if trx.Error != nil {
		return dto.PublicStock{}, trx.Error
	}

	if err := trx.First(&stock, id).Error; err != nil {
		trx.Rollback()
		return dto.PublicStock{}, err
	}

	data := map[string]any{
		"ProductId":  input.ProductId,
		"Qty":        input.Qty,
		"Price":      input.Price,
		"DateEntry":  input.DateEntry,
		"IsActive":   input.IsActive,
		"ModifiedBy": modifierId,
	}

	if err := utils.AssignedKeyModel(&stock, data); err != nil {
		trx.Rollback()
		return dto.PublicStock{}, err
	}

	// Save into DB
	if err := trx.Save(&stock).Error; err != nil {
		trx.Rollback()
		return dto.PublicStock{}, err
	}

	// Commit transaction
	if err := trx.Commit().Error; err != nil {
		return dto.PublicStock{}, err
	}

	var stockWithUser models.Stock
	if err := helpers.PreloadRelationByID(&stockWithUser, stock.ID, []string{"Creator", "Modifier"}); err != nil {
		return dto.PublicStock{}, err
	}

	return utils.ToPublicStock(stockWithUser), nil
}
