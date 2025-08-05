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
	GetProduct(id, code string) (dto.PublicProduct, error)
	GetAllProducts(paginate *dto.PaginationRequest) (*dto.PaginationResponse[dto.PublicProduct], error)
	Create(ctx context.Context, input *dto.CreateProductInput) (dto.PublicProduct, error)
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

func (s *productService) GetProduct(id, code string) (dto.PublicProduct, error) {
	var publicProduct dto.PublicProduct
	var errResult error

	if id != "" {
		parsedID, err := strconv.Atoi(id)
		if err != nil {
			return dto.PublicProduct{}, errs.New("Product ID is not a number!", http.StatusBadRequest)
		}
		product, result := s.repo.FindByID(parsedID)
		publicProduct = utils.ToPublicProduct(product)
		errResult = result.Error
	} else {
		product, result := s.repo.FindByCode(code)
		publicProduct = utils.ToPublicProduct(product)
		errResult = result.Error
	}

	if utils.IsEmptyProduct(publicProduct) {
		return dto.PublicProduct{}, errs.New("Product not found", http.StatusNotFound)
	}

	return publicProduct, errResult
}

func (s *productService) GetAllProducts(request *dto.PaginationRequest) (*dto.PaginationResponse[dto.PublicProduct], error) {
	pg, err := s.repo.FindAll(request)
	if err != nil {
		return nil, err
	}

	if len(pg.Data) == 0 {
		return nil, errs.New("Products not found", http.StatusNotFound)
	}

	return pg, err
}

func (s *productService) Create(ctx context.Context, input *dto.CreateProductInput) (dto.PublicProduct, error) {
	creatorId, ok := ctx.Value(config.UserIDKey).(uint)
	if !ok {
		return dto.PublicProduct{}, errs.New("invalid context session user ID", http.StatusInternalServerError)
	}

	isValid, err := s.validator.ValidateInsertProduct(input)
	if !isValid {
		return dto.PublicProduct{}, err
	}
	return s.repo.Create(input, creatorId)
}
