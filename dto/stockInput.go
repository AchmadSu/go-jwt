package dto

import "time"

type CreateStockInput struct {
	ProductId uint    `json:"product_id" binding:"required,number"`
	Qty       int     `json:"qty" binding:"required,gte=0"`
	Price     float64 `json:"price" binding:"required,gte=0"`
	Date      string  `json:"date" binding:"required"`
	Time      string  `json:"time" binding:"required"`
	DateEntry time.Time
}

type UpdateStockInput struct {
	ProductId *uint    `json:"product_id" binding:"omitempty,number"`
	Qty       *int     `json:"qty" binding:"omitempty,gte=0"`
	Price     *float64 `json:"price" binding:"omitempty,gte=0"`
	Date      string   `json:"date" binding:"omitempty"`
	Time      string   `json:"time" binding:"omitempty"`
	IsActive  *int     `json:"is_active" binding:"omitempty"`
	DateEntry time.Time
}
