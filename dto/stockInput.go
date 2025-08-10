package dto

import "time"

type CreateStockInput struct {
	ProductId uint      `json:"product_id" binding:"required,number"`
	Qty       int       `json:"qty" binding:"required,gte=0"`
	Price     float64   `json:"price" binding:"required,gte=0"`
	DateEntry time.Time `json:"date_entry" binding:"required"`
}

type UpdateStockInput struct {
	ProductId *uint     `json:"product_id" binding:"omitempty,number"`
	Qty       *int      `json:"qty" binding:"omitempty,gte=0"`
	Price     *float64  `json:"price" binding:"omitempty,gte=0"`
	DateEntry time.Time `json:"date_entry" binding:"omitempty"`
	IsActive  *int      `json:"is_active" binding:"omitempty"`
}
