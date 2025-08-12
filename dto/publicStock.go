package dto

import "time"

type ProductID uint
type StockQty int
type StockPrice float64
type StockStatus int

type PublicStock struct {
	ID           uint        `json:"id"`
	ProductID    *uint       `json:"product_id"`
	ProductCode  string      `json:"product_code"`
	ProductName  string      `json:"product_name"`
	IsActive     StockStatus `json:"is_active"`
	Status       string      `json:"status"`
	Qty          int         `json:"qty"`
	Price        float64     `json:"price"`
	DateEntry    time.Time   `json:"date_entry"`
	CreatorId    *uint       `json:"creator_id"`
	CreatorName  string      `json:"creator_name"`
	ModifierId   *uint       `json:"modifier_id"`
	ModifierName string      `json:"modifier_name"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
}
