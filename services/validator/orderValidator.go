package validator

import (
	"fmt"
	"net/http"

	"example.com/m/dto"
	"example.com/m/errs"
	"example.com/m/repositories"
	"example.com/m/utils"
)

type OrderValidatorService interface {
	ValidateCreateOrder(input *dto.CreateOrderInput) error
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

	stockMap, err := v.stockRepo.GetGrandStockPerProductIDs(productIDs)
	if err != nil {
		return err
	}

	insufficientIDs := make([]uint, 0)
	for _, d := range *details {
		grandQty := stockMap[*d.ProductID]
		if grandQty < int(d.Qty) {
			insufficientIDs = append(insufficientIDs, *d.ProductID)
		}
	}

	if len(insufficientIDs) == 0 {
		return nil
	}

	productMap, err := v.productRepo.FindByProductIDs(insufficientIDs)
	if err != nil {
		return err
	}

	for _, d := range *details {
		grandQty := stockMap[*d.ProductID]
		if grandQty < int(d.Qty) {
			if p, ok := productMap[*d.ProductID]; ok && p.Code != "" && p.Name != "" {
				return errs.New(
					fmt.Sprintf("qty of %s (%s) is greater than balance qty", p.Name, p.Code),
					http.StatusBadRequest,
				)
			}
			return errs.New(
				fmt.Sprintf("qty of product id: %d is greater than balance qty", *d.ProductID),
				http.StatusBadRequest,
			)
		}
	}
	return nil
}
