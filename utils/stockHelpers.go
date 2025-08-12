package utils

import (
	"fmt"

	"example.com/m/config"
	"example.com/m/dto"
	"gorm.io/gorm"
)

const StockTable config.TableName = "stocks"

func StockFilterQuery(param *dto.PaginationStockRequest, query *gorm.DB) *gorm.DB {
	if param.ProductID != nil {
		query = query.Where("product_id = ?", *param.ProductID)
	}
	if param.ProductCode != "" {
		query = query.Where("product.code = ?", param.ProductCode)
	}
	if param.ProductName != "" {
		searchTerm := "%" + param.ProductName + "%"
		query = query.Where("product.name ILIKE ?", searchTerm)
	}
	if param.Qty != nil {
		query = query.Where("qty = ?", *param.Qty)
	}
	if param.QtyGreaterThan != nil {
		query = query.Where("qty > ?", *param.QtyGreaterThan)
	}
	if param.QtyLessThan != nil {
		query = query.Where("qty < ?", *param.QtyLessThan)
	}
	if param.Price != nil {
		query = query.Where("price = ?", *param.Price)
	}
	if param.PriceGreaterThan != nil {
		query = query.Where("price > ?", *param.PriceGreaterThan)
	}
	if param.PriceLessThan != nil {
		query = query.Where("price < ?", *param.PriceLessThan)
	}
	if param.DateEntry != "" {
		query = OnDateQuery(query, "date_entry", param.DateEntry)
	}
	if param.DateEntryStart != "" && param.DateEntryEnd != "" {
		dateMap := map[string]any{
			"date_entry_start": param.DateEntryStart,
			"date_entry_end":   param.DateEntryEnd,
		}
		query = ApplyDateRange(dateMap, query, "date_entry_start", "date_entry_end", fmt.Sprintf("%s.date_entry", StockTable))
	}
	return query
}
