package bootstrap

import (
	"example.com/m/repositories"
	"example.com/m/services"
	"example.com/m/services/validator"
)

var OrderService services.OrderService

func InitialOrderService() {
	repo := repositories.NewOrderRepository()
	stockRepo := repositories.NewStockRepository()
	productRepo := repositories.NewProductRepository()
	orderValidator := validator.NewOrderValidatorService(stockRepo, productRepo)
	OrderService = services.NewOrderService(repo, stockRepo, orderValidator)
}
