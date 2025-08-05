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
	GetProduct(data map[string]string) (dto.PublicProduct, error)
	GetAllProducts(paginate *dto.PaginationRequest, data map[string]string) (*dto.PaginationResponse[dto.PublicProduct], error)
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

func (s *productService) GetProduct(data map[string]string) (dto.PublicProduct, error) {
	var publicProduct dto.PublicProduct
	var errResult error

	if data["id"] != "" {
		parsedID, err := strconv.Atoi(data["id"])
		if err != nil {
			return dto.PublicProduct{}, errs.New("Product ID is not a number!", http.StatusBadRequest)
		}
		product, result := s.repo.FindByID(parsedID)
		publicProduct = utils.ToPublicProduct(product)
		errResult = result.Error
	} else if data["name"] != "" {
		product, result := s.repo.FindByName(data["name"])
		publicProduct = utils.ToPublicProduct(product)
		errResult = result.Error
	} else {
		product, result := s.repo.FindByCode(data["code"])
		publicProduct = utils.ToPublicProduct(product)
		errResult = result.Error
	}

	if utils.IsEmptyProduct(publicProduct) {
		return dto.PublicProduct{}, errs.New("Product not found", http.StatusNotFound)
	}

	return publicProduct, errResult
}

func (s *productService) GetAllProducts(request *dto.PaginationRequest, data map[string]string) (*dto.PaginationResponse[dto.PublicProduct], error) {
	var creatorId uint
	var modifierId uint
	var statusProduct int
	switch data["is_active"] {
	case "true":
		statusProduct = 1
	case "false":
		statusProduct = 0
	default:
		statusProduct = 2
	}
	uIntMap := make(map[string]uint)
	for key, strVal := range data {
		if key == "creator_id" || key == "modifier_id" {
			parsedUint64, err := strconv.ParseUint(strVal, 10, 64)
			if err != nil {
				parsedUint64 = 0
			}
			uIntMap[key] = uint(parsedUint64)
		}
	}
	if uIntMap["creator_id"] > 0 {
		creatorId = uIntMap["creator_id"]
	}
	if uIntMap["modifier_id"] > 0 {
		modifierId = uIntMap["modifier_id"]
	}

	pg, err := s.repo.FindAll(request, statusProduct, creatorId, modifierId)
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
