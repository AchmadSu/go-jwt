package utils

import (
	"fmt"

	"example.com/m/config"
	"example.com/m/dto"
	"gorm.io/gorm"
)

const OrderTable config.TableName = "stocks"

func OrderFilterQuery(param *dto.PaginationOrderRequest, query *gorm.DB) *gorm.DB {
	// if param.ProductID != nil {
	// 	query = query.Where("product_id = ?", *param.ProductID)
	// }
	// if param.ProductCode != "" {
	// 	query = query.Where("product.code = ?", param.ProductCode)
	// }
	// if param.ProductName != "" {
	// 	searchTerm := "%" + param.ProductName + "%"
	// 	query = query.Where("product.name ILIKE ?", searchTerm)
	// }
	if param.TotalQty != nil {
		query = query.Where("total_qty = ?", *param.TotalQty)
	}
	if param.TotalQtyGreaterThan != nil {
		query = query.Where("total_qty > ?", *param.TotalQtyGreaterThan)
	}
	if param.TotalQtyLessThan != nil {
		query = query.Where("total_qty < ?", *param.TotalQtyLessThan)
	}
	if param.GrandTotal != nil {
		query = query.Where("grand_total = ?", *param.GrandTotal)
	}
	if param.GrandTotalGreaterThan != nil {
		query = query.Where("grand_total > ?", *param.GrandTotalGreaterThan)
	}
	if param.GrandTotalLessThan != nil {
		query = query.Where("grand_total < ?", *param.GrandTotalLessThan)
	}
	if param.DateEntry != "" {
		query = OnDateQuery(query, "date_entry", param.DateEntry)
	}
	if param.DateEntryStart != "" && param.DateEntryEnd != "" {
		dateMap := map[string]any{
			"date_entry_start": param.DateEntryStart,
			"date_entry_end":   param.DateEntryEnd,
		}
		query = ApplyDateRange(dateMap, query, "date_entry_start", "date_entry_end", fmt.Sprintf("%s.date_entry", OrderTable))
	}
	return query
}
