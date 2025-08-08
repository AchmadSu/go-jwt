package utils

import (
	"example.com/m/dto"
	"example.com/m/models"
)

func ToPublicProduct(p models.Product) dto.PublicProduct {
	return dto.PublicProduct{
		ID:           p.ID,
		Code:         p.Code,
		Name:         p.Name,
		Desc:         p.Desc,
		IsActive:     dto.ProductStatus(p.IsActive),
		Status:       p.IsActive.String(),
		CreatorId:    &p.Creator.ID,
		CreatorName:  p.Creator.Name,
		CreatedAt:    p.CreatedAt,
		ModifierId:   &p.Modifier.ID,
		ModifierName: p.Modifier.Name,
		UpdatedAt:    p.UpdatedAt,
	}
}

func IsEmptyProduct(p dto.PublicProduct) bool {
	return p.ID == 0 && p.Code == ""
}
