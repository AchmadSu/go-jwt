package dto

import "time"

type CreateOrderDetail struct {
	ProductID *uint      `json:"product_id" binding:"required"`
	Qty       StockQty   `json:"qty" binding:"required,gte=1"`
	UnitPrice StockPrice `json:"-"`
}

type CreateOrderInput struct {
	Date      string              `json:"date" binding:"required"`
	Time      string              `json:"time" binding:"required"`
	Details   []CreateOrderDetail `json:"details" binding:"required,min=1,dive"`
	DateEntry time.Time           `json:"-"`
}
