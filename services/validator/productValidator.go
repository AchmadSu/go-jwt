package validator

import (
	"net/http"

	"example.com/m/dto"
	"example.com/m/errs"
	"example.com/m/repositories"
)

type ProductValidatorService interface {
	ValidateInsertProduct(input *dto.CreateProductInput) (bool, error)
	// ValidateUserLogin(input *dto.LoginUserInput) (dto.PublicUser, error)
}

type productValidatorService struct {
	productRepo repositories.ProductRepository
}

func NewProductValidatorService(repo repositories.ProductRepository) *productValidatorService {
	return &productValidatorService{productRepo: repo}
}

func (v *productValidatorService) ValidateInsertProduct(input *dto.CreateProductInput) (bool, error) {
	_, result := v.productRepo.FindByCode(input.Code)
	if result.RowsAffected > 0 {
		return false, errs.New("Code is already exists. Please try another code!", http.StatusBadRequest)
	}

	_, result = v.productRepo.FindByName(input.Name)
	if result.RowsAffected > 0 {
		return false, errs.New("Name is already exists. Please try another name!", http.StatusBadRequest)
	}

	return true, nil
}
