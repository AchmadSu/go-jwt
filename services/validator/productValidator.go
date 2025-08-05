package validator

import (
	"net/http"
	"strconv"

	"example.com/m/errs"
	"example.com/m/repositories"
)

type ProductValidatorService interface {
	ValidateInsertProduct(input map[string]string) (bool, error)
	// ValidateUserLogin(input *dto.LoginUserInput) (dto.PublicUser, error)
}

type productValidatorService struct {
	productRepo repositories.ProductRepository
}

func NewProductValidatorService(repo repositories.ProductRepository) *productValidatorService {
	return &productValidatorService{productRepo: repo}
}

func (v *productValidatorService) ValidateInsertProduct(input map[string]string) (bool, error) {
	if input["code"] == "" && input["name"] == "" && input["desc"] == "" && input["is_active"] == "" {
		return false, errs.New("content must have at least one field", http.StatusBadRequest)
	}
	_, result := v.productRepo.FindByCode(input["code"])
	if result.RowsAffected > 0 {
		return false, errs.New("code is already exists. Please try another code!", http.StatusBadRequest)
	}

	_, result = v.productRepo.FindByName(input["name"])
	if result.RowsAffected > 0 {
		return false, errs.New("name is already exists. Please try another name!", http.StatusBadRequest)
	}

	if input["is_active"] != "" {
		parsedIsActive, err := strconv.Atoi(input["is_active"])
		if err != nil || (parsedIsActive != 0 && parsedIsActive != 1) {
			return false, errs.New("status must contain 1 or 0", http.StatusBadRequest)
		}
	}

	return true, nil
}
