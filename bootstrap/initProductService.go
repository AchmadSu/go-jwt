package bootstrap

import (
	"example.com/m/repositories"
	"example.com/m/services"
	"example.com/m/services/validator"
)

var ProductService services.ProductService

func InitProductService() {
	productRepo := repositories.NewProductRepository()
	userValidator := validator.NewProductValidatorService(productRepo)
	ProductService = services.NewProductService(productRepo, userValidator)
}
