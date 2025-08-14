package dto

type PaginationOrderRequest struct {
	ProductID             *uint  `form:"product_id"`
	ProductCode           string `form:"product_code"`
	ProductName           string `form:"product_name"`
	TotalQty              *uint  `form:"total_qty"`
	TotalQtyLessThan      *uint  `form:"total_qty_less_than"`
	TotalQtyGreaterThan   *uint  `form:"total_qty_greater_than"`
	GrandTotal            *uint  `form:"grand_total"`
	GrandTotalLessThan    *uint  `form:"grand_total_less_than"`
	GrandTotalGreaterThan *uint  `form:"grand_total_greater_than"`
	DateEntry             string `form:"date_entry"`
	DateEntryStart        string `form:"date_entry_start"`
	DateEntryEnd          string `form:"date_entry_end"`
}
