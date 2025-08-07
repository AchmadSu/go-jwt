package services

import (
	"context"
	"net/http"
	"strconv"

	"example.com/m/config"
	"example.com/m/dto"
	"example.com/m/errs"
	"example.com/m/repositories"
	"example.com/m/services/validator"
	"example.com/m/utils"
)

type ProductService interface {
	GetProduct(data *dto.PaginationRequest) (dto.PublicProduct, error)
	GetAllProducts(paginate *dto.PaginationRequest) (*dto.PaginationResponse[dto.PublicProduct], error)
	Create(ctx context.Context, input *dto.CreateProductInput) (dto.PublicProduct, error)
	Update(id int, ctx context.Context, input *dto.UpdateProductInput) (dto.PublicProduct, error)
}

type productService struct {
	repo      repositories.ProductRepository
	validator validator.ProductValidatorService
}

func NewProductService(
	repo repositories.ProductRepository,
	validator validator.ProductValidatorService,
) ProductService {
	return &productService{repo: repo, validator: validator}
}

func (s *productService) GetProduct(data *dto.PaginationRequest) (dto.PublicProduct, error) {
	var publicProduct dto.PublicProduct
	var errResult error

	if *data.ID > 0 {
		product, result := s.repo.FindByID(*data.ID)
		publicProduct = utils.ToPublicProduct(product)
		errResult = result.Error
	} else if data.Name != "" {
		product, result := s.repo.FindByName(data.Name)
		publicProduct = utils.ToPublicProduct(product)
		errResult = result.Error
	} else {
		product, result := s.repo.FindByCode(data.Code)
		publicProduct = utils.ToPublicProduct(product)
		errResult = result.Error
	}

	if utils.IsEmptyProduct(publicProduct) {
		return dto.PublicProduct{}, errs.New("Product not found", http.StatusNotFound)
	}

	return publicProduct, errResult
}

func (s *productService) GetAllProducts(request *dto.PaginationRequest) (*dto.PaginationResponse[dto.PublicProduct], error) {

	if request.CreateDateStart != "" && request.CreateDateEnd != "" {
		err := utils.ValidateDateRange(request.CreateDateStart, request.CreateDateEnd)
		if err != nil {
			return nil, err
		}
	}

	if request.UpdateDateStart != "" && request.UpdateDateEnd != "" {
		err := utils.ValidateDateRange(request.CreateDateStart, request.CreateDateEnd)
		if err != nil {
			return nil, err
		}
	}

	pg, err := s.repo.FindAll(request)
	if err != nil {
		return nil, err
	}

	if len(pg.Data) == 0 {
		return nil, errs.New("Products not found", http.StatusNotFound)
	}

	return pg, nil
}

func (s *productService) Create(ctx context.Context, input *dto.CreateProductInput) (dto.PublicProduct, error) {
	creatorId, ok := ctx.Value(config.UserIDKey).(uint)
	if !ok {
		return dto.PublicProduct{}, errs.New("invalid context session user ID", http.StatusInternalServerError)
	}

	mapValidator := map[string]string{
		"code": input.Code,
		"name": input.Name,
	}

	isValid, err := s.validator.ValidateInsertProduct(mapValidator)
	if !isValid {
		return dto.PublicProduct{}, err
	}
	return s.repo.Create(input, creatorId)
}

func (s *productService) Update(id int, ctx context.Context, input *dto.UpdateProductInput) (dto.PublicProduct, error) {
	modifierId, ok := ctx.Value(config.UserIDKey).(uint)
	if !ok {
		return dto.PublicProduct{}, errs.New("invalid context session user ID", http.StatusBadRequest)
	}

	mapValidator := map[string]string{
		"code":      input.Code,
		"name":      input.Name,
		"is_active": strconv.Itoa(input.IsActive),
	}

	isValid, err := s.validator.ValidateInsertProduct(mapValidator)
	if !isValid {
		return dto.PublicProduct{}, err
	}

	return s.repo.Update(id, input, modifierId)
}
