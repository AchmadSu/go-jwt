package bootstrap

import (
	"example.com/m/repositories"
	"example.com/m/services"
	"example.com/m/services/validator"
)

var StockService services.StockService

func InitStockService() {
	stockRepo := repositories.NewStockRepository()
	productRepo := repositories.NewProductRepository()
	stockValidator := validator.NewStockValidatorService(stockRepo, productRepo)
	StockService = services.NewStockService(stockRepo, productRepo, stockValidator)
}
