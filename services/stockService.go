package services

import (
	"context"
	"net/http"

	"example.com/m/config"
	"example.com/m/dto"
	"example.com/m/errs"
	"example.com/m/repositories"
	"example.com/m/services/validator"
	"example.com/m/utils"
)

type StockService interface {
	GetStock(data *dto.PaginationRequest) (dto.PublicStock, error)
	GetAllStock(paginate *dto.PaginationRequest, stockData *dto.PaginationStockRequest) (*dto.PaginationResponse[dto.PublicStock], error)
	CreateStock(ctx context.Context, input *dto.CreateStockInput) (dto.PublicStock, error)
	UpdateStock(id int, ctx context.Context, input *dto.UpdateStockInput) (dto.PublicStock, error)
}

type stockService struct {
	repo        repositories.StockRepository
	productRepo repositories.ProductRepository
	validator   validator.StockValidatorService
}

func NewStockService(
	repo repositories.StockRepository,
	productRepo repositories.ProductRepository,
	validator validator.StockValidatorService,
) StockService {
	return &stockService{
		repo:        repo,
		productRepo: productRepo,
		validator:   validator,
	}
}

func (s *stockService) GetStock(data *dto.PaginationRequest) (dto.PublicStock, error) {
	var publicStock dto.PublicStock
	var errResult error

	if data.ID != nil && *data.ID > 0 {
		product, result := s.repo.FindByStockID(*data.ID)
		publicStock = utils.ToPublicStock(product)
		errResult = result.Error
	} else {
		return dto.PublicStock{}, errs.New("invalid id request", http.StatusBadRequest)
	}

	if utils.IsEmptyStock(publicStock) {
		return dto.PublicStock{}, errs.New("stock not found", http.StatusNotFound)
	}

	return publicStock, errResult
}

func (s *stockService) GetAllStock(request *dto.PaginationRequest, stockData *dto.PaginationStockRequest) (*dto.PaginationResponse[dto.PublicStock], error) {

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

	pg, err := s.repo.FindAllStocks(request, stockData)
	if err != nil {
		return nil, err
	}

	if len(pg.Data) == 0 {
		return nil, errs.New("Products not found", http.StatusNotFound)
	}

	return pg, nil
}

func (s *stockService) CreateStock(ctx context.Context, input *dto.CreateStockInput) (dto.PublicStock, error) {
	creatorId, ok := ctx.Value(config.UserIDKey).(uint)
	if !ok {
		return dto.PublicStock{}, errs.New("invalid context session user ID", http.StatusInternalServerError)
	}

	isValid, err := s.validator.ValidateInsertStock(input)
	if !isValid {
		return dto.PublicStock{}, err
	}
	dateEntry, _ := utils.MergeDateTime(input.Date, input.Time)
	input.DateEntry = dateEntry
	return s.repo.CreateStock(input, creatorId)
}

func (s *stockService) UpdateStock(id int, ctx context.Context, input *dto.UpdateStockInput) (dto.PublicStock, error) {
	modifierId, ok := ctx.Value(config.UserIDKey).(uint)
	if !ok {
		return dto.PublicStock{}, errs.New("invalid context session user ID", http.StatusBadRequest)
	}

	isValid, err := s.validator.ValidateUpdateStock(id, input)
	if !isValid {
		return dto.PublicStock{}, err
	}
	if input.Date != "" && input.Time != "" {
		dateEntry, _ := utils.MergeDateTime(input.Date, input.Time)
		input.DateEntry = dateEntry
	}
	return s.repo.UpdateStock(id, input, modifierId)
}
