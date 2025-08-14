package repositories

import (
	"net/http"

	"example.com/m/config"
	"example.com/m/dto"
	"example.com/m/errs"
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
	GetGrandStockPerProductID(productID uint) (int, error)
	GetGrandStockPerProductIDs(productIDs []uint) (map[uint]int, error)
	GetStockByProductID(productID uint) ([]models.Stock, error)
	GetStockByProductIDs(productIDs []uint) (map[*uint][]models.Stock, error)
	UpdateStockByIDs(trx *gorm.DB, stockMap map[uint]int, modifierID uint) error
}

type stockRepository struct{}

type result struct {
	Total int
}

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

func (r *stockRepository) GetStockByProductID(productID uint) ([]models.Stock, error) {
	var stocks []models.Stock
	result := initializers.DB.Where("product_id = ?", productID).Find(&stocks)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return []models.Stock{}, errs.New("Stock not found", http.StatusNotFound)
		}
		return []models.Stock{}, errs.New(utils.GetSafeErrorMessage(result.Error, "Unknown stock error occurred"), http.StatusInternalServerError)
	}

	return stocks, nil
}

func (r *stockRepository) GetStockByProductIDs(productIDs []uint) (map[*uint][]models.Stock, error) {
	var results []models.Stock

	err := initializers.DB.
		Where("id_product IN ?", productIDs).
		Where("stock > 0").
		Order("id_product ASC").
		Order("date_entry ASC").
		Find(&results).Error

	if err != nil {
		return nil, err
	}

	stockMap := make(map[*uint][]models.Stock)
	for _, stock := range results {
		stockMap[stock.ProductID] = append(stockMap[stock.ProductID], stock)
	}

	return stockMap, nil
}

func (r *stockRepository) GetGrandStockPerProductID(productID uint) (int, error) {
	var res result
	err := initializers.DB.
		Model(&models.Stock{}).
		Where("id_product = ?", productID).
		Select("COALESCE(SUM(qty), 0) as total").
		Scan(&res).Error
	if err != nil {
		return 0, errs.New(utils.GetSafeErrorMessage(err, "Out of stock for this product"), http.StatusNotFound)
	}
	return res.Total, nil
}

func (r *stockRepository) GetGrandStockPerProductIDs(productIDs []uint) (map[uint]int, error) {
	var results []struct {
		ProductID uint `gorm:"column:id_product"`
		TotalQty  int  `gorm:"column:total_qty"`
	}

	err := initializers.DB.Table("stocks").
		Select("id_product, SUM(qty) as total_qty").
		Where("id_product IN ?", productIDs).
		Group("id_product").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	stockMap := make(map[uint]int)
	for _, r := range results {
		stockMap[r.ProductID] = r.TotalQty
	}

	return stockMap, nil
}

func (r *stockRepository) FindAllStocks(request *dto.PaginationRequest, stockRequest *dto.PaginationStockRequest) (*dto.PaginationResponse[dto.PublicStock], error) {
	query := initializers.DB.Model(&models.Stock{}).
		Joins("LEFT JOIN products AS product ON product.id = stocks.product_id").
		Joins("LEFT JOIN users AS creator ON creator.id = stocks.created_by").
		Joins("LEFT JOIN users AS modifier ON modifier.id = stocks.modified_by").
		Select([]string{
			"stocks.id AS id",
			"product.id AS product_id",
			"product.code AS product_code",
			"product.name AS product_name",
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
		`product_id`,
		`product_name`,
		`product_code`,
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
		`creator.name`,
		`modifier.name`,
	}
	defaultOrder := "product.name asc"
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
		ProductID: &input.ProductID,
		Qty:       models.StockQty(input.Qty),
		Price:     models.StockPrice(input.Price),
		DateEntry: input.DateEntry,
		CreatedBy: &creatorId,
	}
	result := initializers.DB.Create(&stock)
	if result.Error != nil {
		return dto.PublicStock{}, result.Error
	}

	var finalStock models.Stock
	err := helpers.PreloadRelationByID(&finalStock, stock.ID, []string{"Creator", "Modifier", "Product"})
	if err != nil {
		return dto.PublicStock{}, err
	}

	return utils.ToPublicStock(finalStock), nil
}

func (r *stockRepository) UpdateStock(id int, input *dto.UpdateStockInput, modifierId uint) (dto.PublicStock, error) {
	var stock models.Stock

	trx := initializers.DB.Begin()
	if trx.Error != nil {
		trx.Rollback()
		return dto.PublicStock{}, trx.Error
	}

	if err := trx.First(&stock, id).Error; err != nil {
		trx.Rollback()
		return dto.PublicStock{}, err
	}

	data := map[string]any{
		"ProductID":  input.ProductID,
		"Qty":        input.Qty,
		"Price":      input.Price,
		"IsActive":   input.IsActive,
		"DateEntry":  input.DateEntry,
		"ModifiedBy": &modifierId,
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

	var finalStock models.Stock
	err := helpers.PreloadRelationByID(&finalStock, stock.ID, []string{"Creator", "Modifier", "Product"})
	if err != nil {
		return dto.PublicStock{}, err
	}

	return utils.ToPublicStock(finalStock), nil
}

func (r *stockRepository) UpdateStockByIDs(trx *gorm.DB, stockMap map[uint]int, modifierID uint) error {
	if len(stockMap) == 0 {
		return errs.New("no stock to update", http.StatusBadRequest)
	}

	stockIDs := make([]uint, 0, len(stockMap))
	for id := range stockMap {
		if id == 0 {
			return errs.New("stock ID cannot be nil", http.StatusBadRequest)
		}
		stockIDs = append(stockIDs, id)
	}

	var count int64
	if err := trx.Table("stocks").
		Where("id IN ?", stockIDs).
		Count(&count).Error; err != nil {
		return errs.New("failed to check stock IDs", http.StatusInternalServerError)
	}
	if count != int64(len(stockIDs)) {
		return errs.New("one or some stock IDs not found", http.StatusNotFound)
	}

	for id, qty := range stockMap {
		if err := trx.Table("stocks").
			Where("id = ?", id).
			Updates(map[string]any{
				"qty":        gorm.Expr("qty - ?", qty),
				"updated_by": modifierID,
			}).Error; err != nil {
			return errs.New("failed to update one or more stock qty", http.StatusNotModified)
		}
	}

	return nil
}
