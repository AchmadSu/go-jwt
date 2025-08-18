package validator

import (
	"fmt"
	"net/http"
	"strings"

	"example.com/m/dto"
	"example.com/m/errs"
	"example.com/m/repositories"
	"example.com/m/utils"
)

type OrderValidatorService interface {
	ValidateInsertOrder(input *dto.CreateOrderInput) error
	ValidateDetailOrder(details *[]dto.CreateOrderDetail) error
}

type orderValidatorService struct {
	stockRepo   repositories.StockRepository
	productRepo repositories.ProductRepository
}

func NewOrderValidatorService(stockRepo repositories.StockRepository, productRepo repositories.ProductRepository) *orderValidatorService {
	return &orderValidatorService{stockRepo: stockRepo, productRepo: productRepo}
}

func (v *orderValidatorService) ValidateInsertOrder(input *dto.CreateOrderInput) error {
	if err := v.ValidateDetailOrder(&input.Details); err != nil {
		return err
	}
	if _, err := utils.MergeDateTime(input.Date, input.Time); err != nil {
		return errs.New("invalid order date entry format", http.StatusBadRequest)
	}
	return nil
}

func (v *orderValidatorService) ValidateDetailOrder(details *[]dto.CreateOrderDetail) error {
	if details == nil || len(*details) == 0 {
		return errs.New("order must have at least one detail product", http.StatusBadRequest)
	}

	productIDs := make([]uint, 0)
	seen := make(map[uint]struct{})
	for _, d := range *details {
		if d.ProductID == nil || *d.ProductID == 0 {
			return errs.New("invalid product id", http.StatusBadRequest)
		}
		if d.Qty <= 0 {
			return errs.New("qty of the product must greater than 0", http.StatusBadRequest)
		}

		if _, exists := seen[*d.ProductID]; !exists {
			productIDs = append(productIDs, *d.ProductID)
			seen[*d.ProductID] = struct{}{}
		}
	}

	productMap, err := v.productRepo.FindByProductIDs(productIDs)
	if err != nil {
		return err
	}

	var inactiveMessage []string
	for _, p := range productMap {
		if p.IsActive != 1 {
			inactiveMessage = append(inactiveMessage,
				fmt.Sprintf("product: %s (%s) is inactive", p.Name, p.Code))
		}
	}
	if len(inactiveMessage) > 0 {
		return errs.New(strings.Join(inactiveMessage, "\n"), http.StatusBadRequest)
	}

	stockMap, err := v.stockRepo.GetGrandStockPerProductIDs(productIDs)
	if err != nil {
		return err
	}

	var notFoundIDs []string
	var insufficientMessage []string
	for _, d := range *details {
		grandQty, ok := stockMap[*d.ProductID]
		if !ok {
			notFoundIDs = append(notFoundIDs, fmt.Sprintf("%d", *d.ProductID))
			continue
		}
		if grandQty < int(d.Qty) {
			if p, ok := productMap[*d.ProductID]; ok {
				insufficientMessage = append(insufficientMessage,
					fmt.Sprintf("qty of %s (%s) is greater than balance qty", p.Name, p.Code))
			} else {
				insufficientMessage = append(insufficientMessage,
					fmt.Sprintf("qty of product id: %d is greater than balance qty", *d.ProductID))
			}
		}
	}

	if len(notFoundIDs) > 0 {
		return errs.New(fmt.Sprintf("product ids not found: %s",
			strings.Join(notFoundIDs, ", ")), http.StatusBadRequest)
	}

	if len(insufficientMessage) > 0 {
		return errs.New(strings.Join(insufficientMessage, "\n"), http.StatusBadRequest)
	}

	return nil
}
