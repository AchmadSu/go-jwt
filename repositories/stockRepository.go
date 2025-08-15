package repositories

import (
	"fmt"
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
	GetStockByProductIDs(productIDs []uint) (map[uint][]models.Stock, error)
	UpdateStockQtyByIDs(trx *gorm.DB, stockMap map[uint]int, operator string, modifierID uint) error
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

func (r *stockRepository) GetStockByProductIDs(productIDs []uint) (map[uint][]models.Stock, error) {
	var results []models.Stock

	err := initializers.DB.
		Table("stocks").
		Where("is_active = 1").
		Where("product_id IN ?", productIDs).
		Where("qty > 0").
		Order("product_id ASC").
		Order("date_entry ASC").
		Find(&results).Error

	if err != nil {
		return nil, err
	}

	stockMap := make(map[uint][]models.Stock)
	for _, stock := range results {
		stockMap[*stock.ProductID] = append(stockMap[*stock.ProductID], stock)
	}

	return stockMap, nil
}

func (r *stockRepository) GetGrandStockPerProductID(productID uint) (int, error) {
	var res result
	err := initializers.DB.
		Model(&models.Stock{}).
		Where("product_id = ?", productID).
		Where("is_active = 1").
		Select("COALESCE(SUM(qty), 0) as total").
		Scan(&res).Error
	if err != nil {
		return 0, errs.New(utils.GetSafeErrorMessage(err, "Out of stock for this product"), http.StatusNotFound)
	}
	return res.Total, nil
}

func (r *stockRepository) GetGrandStockPerProductIDs(productIDs []uint) (map[uint]int, error) {
	var results []struct {
		ProductID uint `gorm:"column:product_id"`
		TotalQty  int  `gorm:"column:total_qty"`
	}

	err := initializers.DB.Table("stocks").
		Select("product_id, SUM(qty) as total_qty").
		Where("is_active = 1").
		Where("product_id IN ?", productIDs).
		Group("product_id").
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

func (r *stockRepository) UpdateStockQtyByIDs(trx *gorm.DB, stockMap map[uint]int, operator string, modifierID uint) error {
	if len(stockMap) == 0 {
		return errs.New("no stock to update", http.StatusBadRequest)
	}

	allowedOperator := []string{"+", "-"}
	if !utils.ContainsString(allowedOperator, operator) {
		return errs.New("invalid operator qty stock. only '+' or '-'", http.StatusBadRequest)
	}

	stockIDs := make([]uint, 0, len(stockMap))
	for id := range stockMap {
		if id == 0 {
			return errs.New("stock ID cannot be 0", http.StatusBadRequest)
		}
		stockIDs = append(stockIDs, id)
	}

	var existingIDs []uint
	if err := trx.Table("stocks").Where("id IN ?", stockIDs).Pluck("id", &existingIDs).Error; err != nil {
		return errs.New("failed to check stock IDs", http.StatusInternalServerError)
	}

	fmt.Printf("stockIDS: %v .existingIDS: %v", stockIDs, existingIDs)

	if len(existingIDs) != len(stockIDs) {
		missing := utils.DifferenceUint(stockIDs, existingIDs)
		return errs.New(fmt.Sprintf("stock IDs not found: %v", missing), http.StatusNotFound)
	}

	caseExpr := "CASE"
	for id, qty := range stockMap {
		caseExpr += fmt.Sprintf(" WHEN id = %d THEN qty %s %d", id, operator, qty)
	}
	caseExpr += " END"

	result := trx.Exec(
		fmt.Sprintf("UPDATE stocks SET qty = %s, modified_by = ? WHERE id IN ?", caseExpr),
		modifierID, stockIDs,
	).Debug()
	if err := result.Error; err != nil {
		return errs.New(utils.GetSafeErrorMessage(err, "unknown update stock qty by ids error occurred"), http.StatusInternalServerError)
	}

	return nil
}
