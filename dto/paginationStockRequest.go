package dto

type PaginationStockRequest struct {
	ProductID        *uint  `form:"product_id"`
	ProductCode      string `form:"product_code"`
	ProductName      string `form:"product_name"`
	Qty              *uint  `form:"stock_qty"`
	QtyLessThan      *uint  `form:"stock_less_than"`
	QtyGreaterThan   *uint  `form:"stock_greater_than"`
	Price            *uint  `form:"price"`
	PriceLessThan    *uint  `form:"price_less_than"`
	PriceGreaterThan *uint  `form:"price_greater_than"`
	DateEntry        string `form:"date_entry"`
	DateEntryStart   string `form:"date_entry_start"`
	DateEntryEnd     string `form:"date_entry_end"`
}
