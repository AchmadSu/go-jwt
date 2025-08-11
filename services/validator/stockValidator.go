package validator

import (
	"net/http"

	"example.com/m/dto"
	"example.com/m/errs"
	"example.com/m/repositories"
	"example.com/m/utils"
)

type StockValidatorService interface {
	ValidateInsertStock(input *dto.CreateStockInput) (bool, error)
	ValidateUpdateStock(id int, input *dto.UpdateStockInput) (bool, error)
}

type stockValidatorService struct {
	repo        repositories.StockRepository
	productRepo repositories.ProductRepository
}

func NewStockValidatorService(repo repositories.StockRepository, productRepo repositories.ProductRepository) *stockValidatorService {
	return &stockValidatorService{repo: repo, productRepo: productRepo}
}

func (v *stockValidatorService) ValidateInsertStock(input *dto.CreateStockInput) (bool, error) {
	if input.Qty < 0 || input.Price <= 0 || input.Date == "" || input.Time == "" {
		return false, errs.New("invalid create stock request ", http.StatusBadRequest)
	}
	_, result := v.productRepo.FindByProductID(int(input.ProductId))
	if result.RowsAffected == 0 {
		return false, errs.New("product not found", http.StatusNotFound)
	}

	if _, err := utils.MergeDateTime(input.Date, input.Time); err != nil {
		return false, errs.New("invalid date entry format", http.StatusBadRequest)
	}

	return true, nil
}

func (v *stockValidatorService) ValidateUpdateStock(id int, input *dto.UpdateStockInput) (bool, error) {
	_, result := v.repo.FindByStockID(id)
	if result.RowsAffected == 0 {
		return false, errs.New("stock not found", http.StatusNotFound)
	}
	if input.ProductId != nil {
		_, result := v.productRepo.FindByProductID(int(*input.ProductId))
		if result.RowsAffected == 0 {
			return false, errs.New("product not found", http.StatusNotFound)
		}
	}
	if input.Qty != nil && *input.Qty < 0 {
		return false, errs.New("qty must contains greater equal than zero", http.StatusBadRequest)
	}
	if input.Price != nil && *input.Price <= 0 {
		return false, errs.New("price must contains greater than zero", http.StatusBadRequest)
	}
	if input.IsActive != nil && *input.IsActive < 0 {
		return false, errs.New("invalid create stock request ", http.StatusBadRequest)
	}

	if input.Date != "" && input.Time != "" {
		if _, err := utils.MergeDateTime(input.Date, input.Time); err != nil {
			return false, errs.New("invalid date entry format", http.StatusBadRequest)
		}
	}

	return true, nil
}
