package repositories

import (
	"errors"
	"net/http"
	"time"

	"example.com/m/config"
	"example.com/m/dto"
	"example.com/m/errs"
	"example.com/m/initializers"
	"example.com/m/models"
	"example.com/m/utils"
	"gorm.io/gorm"
)

const OrderTable config.TableName = "orders"

type OrderRepository interface {
	FindByOrderID(id int) (dto.PublicOrderWithDetail, error)
	FindByOrderCode(code string) (dto.PublicOrderWithDetail, error)
	FindAllOrders(paginate *dto.PaginationRequest, orderRequest *dto.PaginationOrderRequest) (*dto.PaginationResponse[dto.PublicOrder], error)
	CreateOrder(input *dto.CreateOrderInput, updateStockMap map[uint]int, creatorId uint) (dto.PublicOrderWithDetail, error)
}

type orderRepository struct{}

func NewOrderRepository() OrderRepository {
	return &orderRepository{}
}

func (r *orderRepository) FindByOrderID(id int) (dto.PublicOrderWithDetail, error) {
	var order models.Order
	result := initializers.DB.Preload("Creator").Preload("OrderDetails.Product").First(&order, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return dto.PublicOrderWithDetail{}, errs.New("order not found", http.StatusNotFound)
		}
		return dto.PublicOrderWithDetail{}, errs.New("unknown order error occurred", http.StatusBadRequest)
	}

	return utils.ToPublicOrder(order, order.OrderDetails), nil
}

func (r *orderRepository) FindByOrderCode(code string) (dto.PublicOrderWithDetail, error) {
	var order models.Order
	result := initializers.DB.Preload("Creator").Preload("OrderDetails.Product").First(&order, "code = ?", code)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return dto.PublicOrderWithDetail{}, errs.New("order not found", http.StatusNotFound)
		}
		return dto.PublicOrderWithDetail{}, errs.New("unknown order error occurred", http.StatusBadRequest)
	}

	return utils.ToPublicOrder(order, order.OrderDetails), nil
}

func (r *orderRepository) FindAllOrders(request *dto.PaginationRequest, orderRequest *dto.PaginationOrderRequest) (*dto.PaginationResponse[dto.PublicOrder], error) {
	query := initializers.DB.Model(&models.Order{}).
		Joins("LEFT JOIN users AS creator ON creator.id = orders.created_by").
		Select([]string{
			"orders.id AS id",
			"orders.total_qty",
			"orders.grand_total",
			"orders.date_entry",
			"orders.created_by AS creator_id",
			"orders.created_at",
			"creator.name AS creator_name",
		})

	query = utils.OrderFilterQuery(orderRequest, query)
	query = utils.FilterQuery(request, query, string(OrderTable)).Debug()

	allowedSortFields := []string{
		`id`,
		`code`,
		`total_qty`,
		`grand_total`,
		`creator_name`,
		`created_at`,
	}
	searchFields := []string{
		`creator.name`,
		`order.code`,
	}
	defaultOrder := "created_at desc"
	pageResult, err := utils.Paginate[dto.PublicOrder](request, query, allowedSortFields, defaultOrder, searchFields)

	if err != nil {
		return nil, err
	}

	return pageResult, nil
}

func (r *orderRepository) CreateOrder(input *dto.CreateOrderInput, updateStockMap map[uint]int, creatorId uint) (dto.PublicOrderWithDetail, error) {
	var grandTotal float64
	var totalQty int

	var returnValue dto.PublicOrderWithDetail

	err := initializers.DB.Transaction(func(trx *gorm.DB) error {
		productIDs := make([]uint, 0, len(input.Details))
		for _, d := range input.Details {
			productIDs = append(productIDs, *d.ProductID)
		}

		var count int64
		if err := trx.Model(&models.Product{}).
			Where("id IN ?", productIDs).
			Count(&count).Error; err != nil {
			return err
		}

		if count != int64(len(productIDs)) {
			return errs.New("one or some of these products not found", http.StatusBadRequest)
		}

		orderDetails := make([]models.OrderDetails, 0, len(input.Details))
		for _, detail := range input.Details {
			total := float64(detail.Qty) * float64(detail.UnitPrice)
			orderDetails = append(orderDetails, models.OrderDetails{
				ProductID: detail.ProductID,
				Qty:       models.StockQty(detail.Qty),
				UnitPrice: models.UnitPrice(detail.UnitPrice),
				Total:     models.Total(total),
			})
			totalQty += int(detail.Qty)
			grandTotal += total
		}

		code, err := utils.GenerateCodeOrder(trx, time.Now())
		if err != nil {
			return err
		}

		order := models.Order{
			Code:       code,
			DateEntry:  input.DateEntry,
			TotalQty:   models.StockQty(totalQty),
			GrandTotal: models.Total(grandTotal),
			CreatedBy:  &creatorId,
		}
		if err := trx.Create(&order).Error; err != nil {
			return err
		}

		for i := range orderDetails {
			orderDetails[i].OrderID = &order.ID
		}

		if err := trx.CreateInBatches(orderDetails, len(orderDetails)).Error; err != nil {
			return err
		}

		var stockRepo stockRepository

		if err := stockRepo.UpdateStockQtyByIDs(trx, updateStockMap, "-", creatorId); err != nil {
			return err
		}

		var finalOrder models.Order
		if err := trx.Preload("Creator").
			Preload("OrderDetails.Product").
			First(&finalOrder, order.ID).Error; err != nil {
			return err
		}

		returnValue = utils.ToPublicOrder(finalOrder, finalOrder.OrderDetails)
		return nil
	})

	return returnValue, err
}
