package utils

import (
	"example.com/m/dto"
	"example.com/m/models"
)

func ToPublicStock(s models.Stock) dto.PublicStock {
	return dto.PublicStock{
		ID:           s.ID,
		ProductId:    s.ProductId,
		Qty:          int(s.Qty),
		Price:        float64(s.Price),
		DateEntry:    s.DateEntry,
		IsActive:     dto.StockStatus(s.IsActive),
		Status:       s.IsActive.String(),
		CreatorId:    &s.Creator.ID,
		CreatorName:  s.Creator.Name,
		CreatedAt:    s.CreatedAt,
		ModifierId:   &s.Modifier.ID,
		ModifierName: s.Modifier.Name,
		UpdatedAt:    s.UpdatedAt,
	}
}

func IsEmptyStock(s dto.PublicStock) bool {
	return s.ID == 0 && s.ProductId == nil
}
