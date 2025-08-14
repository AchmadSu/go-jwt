package services

import (
	"context"
	"fmt"
	"net/http"

	"example.com/m/config"
	"example.com/m/dto"
	"example.com/m/errs"
	"example.com/m/repositories"
	"example.com/m/services/validator"
	"example.com/m/utils"
)

type OrderService interface {
	GetOrder(data *dto.PaginationRequest) (dto.PublicOrderWithDetail, error)
	GetAllOrder(paginate *dto.PaginationRequest, orderData *dto.PaginationOrderRequest) (*dto.PaginationResponse[dto.PublicOrder], error)
	CreateOrder(ctx context.Context, input *dto.CreateOrderInput) (dto.PublicOrderWithDetail, error)
}

type orderService struct {
	repo        repositories.OrderRepository
	stockRepo   repositories.StockRepository
	productRepo repositories.ProductRepository
	validator   validator.OrderValidatorService
}

func NewOrderService(
	repo repositories.OrderRepository,
	stockRepo repositories.StockRepository,
	productRepo repositories.ProductRepository,
	validator validator.OrderValidatorService,
) OrderService {
	return &orderService{
		repo:        repo,
		stockRepo:   stockRepo,
		productRepo: productRepo,
		validator:   validator,
	}
}

func (s *orderService) GetOrder(data *dto.PaginationRequest) (dto.PublicOrderWithDetail, error) {
	var publicOrder dto.PublicOrderWithDetail
	var errResult error

	if data.ID != nil && *data.ID > 0 {
		order, err := s.repo.FindByOrderID(*data.ID)
		publicOrder = order
		errResult = err
	} else if data.Code != "" {
		order, err := s.repo.FindByOrderCode(data.Code)
		publicOrder = order
		errResult = err
	} else {
		return dto.PublicOrderWithDetail{}, errs.New("invalid request", http.StatusBadRequest)
	}

	if utils.IsEmptyOrder(publicOrder) {
		return dto.PublicOrderWithDetail{}, errs.New("order not found", http.StatusNotFound)
	}

	return publicOrder, errResult
}

func (s *orderService) GetAllOrder(request *dto.PaginationRequest, orderData *dto.PaginationOrderRequest) (*dto.PaginationResponse[dto.PublicOrder], error) {

	if request.CreateDateStart != "" && request.CreateDateEnd != "" {
		err := utils.ValidateDateRange(request.CreateDateStart, request.CreateDateEnd)
		if err != nil {
			return nil, err
		}
	}

	if request.UpdateDateStart != "" && request.UpdateDateEnd != "" {
		err := utils.ValidateDateRange(request.UpdateDateStart, request.UpdateDateEnd)
		if err != nil {
			return nil, err
		}
	}

	pg, err := s.repo.FindAllOrders(request, orderData)
	if err != nil {
		return nil, err
	}

	if len(pg.Data) == 0 {
		return nil, errs.New("Products not found", http.StatusNotFound)
	}

	return pg, nil
}

func (s *orderService) CreateOrder(ctx context.Context, input *dto.CreateOrderInput) (dto.PublicOrderWithDetail, error) {
	creatorId, ok := ctx.Value(config.UserIDKey).(uint)
	if !ok {
		return dto.PublicOrderWithDetail{}, errs.New("missing or invalid context session user ID", http.StatusInternalServerError)
	}

	err := s.validator.ValidateCreateOrder(input)
	if err != nil {
		return dto.PublicOrderWithDetail{}, err
	}

	detailResult, updateStockMap, err := s.AssignFifoOrderDetail(&input.Details)
	if err != nil {
		return dto.PublicOrderWithDetail{}, err
	}

	dateEntry, err := utils.MergeDateTime(input.Date, input.Time)
	if err != nil {
		return dto.PublicOrderWithDetail{}, err
	}

	input.Details = *detailResult
	input.DateEntry = dateEntry

	return s.repo.CreateOrder(input, updateStockMap, creatorId)
}

func (s *orderService) AssignFifoOrderDetail(details *[]dto.CreateOrderDetail) (*[]dto.CreateOrderDetail, map[uint]int, error) {
	if details == nil || len(*details) == 0 {
		return nil, nil, errs.New("order must have at least one detail product", http.StatusBadRequest)
	}
	productIDs := make([]uint, 0)
	updateStockMap := make(map[uint]int)
	for _, d := range *details {
		if d.ProductID == nil {
			return nil, nil, errs.New("product ID cannot be nil", http.StatusBadRequest)
		}
		productIDs = append(productIDs, *d.ProductID)
	}

	stockMap, err := s.stockRepo.GetStockByProductIDs(productIDs)
	if err != nil {
		return nil, nil, err
	}

	for i, d := range *details {
		if d.ProductID == nil {
			return nil, nil, errs.New("product ID cannot be nil", http.StatusBadRequest)
		}

		stocks, ok := stockMap[d.ProductID]
		if !ok || len(stocks) == 0 {
			return nil, nil, errs.New(fmt.Sprintf("no stock found for product ID: %d", *d.ProductID), http.StatusBadRequest)
		}

		requestQty := d.Qty
		var total float64

		for _, stock := range stocks {
			if requestQty <= 0 {
				break
			}

			if requestQty > dto.StockQty(stock.Qty) {
				total += float64(stock.Price) * float64(stock.Qty)
				requestQty -= dto.StockQty(stock.Qty)
				updateStockMap[stock.ID] = int(stock.Qty)
			} else {
				total += float64(stock.Price) * float64(requestQty)
				requestQty = 0
				updateStockMap[stock.ID] = int(requestQty)
				break
			}
		}

		if d.Qty > 0 && total > 0 {
			(*details)[i].UnitPrice = dto.StockPrice(float64(total) / float64(d.Qty))
		} else {
			(*details)[i].UnitPrice = 0
		}
	}

	return details, updateStockMap, nil
}
