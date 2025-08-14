package dto

import "time"

type Total float64

type OrderDetail struct {
	ProductID uint       `json:"product_id"`
	Code      string     `json:"product_code"`
	Name      string     `json:"product_name"`
	Qty       StockQty   `json:"qty"`
	UnitPrice StockPrice `json:"unit_price"`
	Total     Total      `json:"total"`
}

type PublicOrder struct {
	ID          uint      `json:"id"`
	Code        string    `json:"code"`
	DateEntry   time.Time `json:"date_entry"`
	TotalQty    int       `json:"total_qty"`
	GrandTotal  Total     `json:"grand_total"`
	CreatorId   *uint     `json:"creator_id"`
	CreatorName string    `json:"creator_name"`
	CreatedAt   time.Time `json:"created_at"`
}

type PublicOrderWithDetail struct {
	ID          uint          `json:"id"`
	Code        string        `json:"code"`
	DateEntry   time.Time     `json:"date_entry"`
	TotalQty    int           `json:"total_qty"`
	GrandTotal  Total         `json:"grand_total"`
	CreatorId   *uint         `json:"creator_id"`
	CreatorName string        `json:"creator_name"`
	Details     []OrderDetail `json:"details"`
	CreatedAt   time.Time     `json:"created_at"`
}
