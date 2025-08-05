package dto

import "time"

type ProductStatus int

type PublicProduct struct {
	ID           uint          `json:"id"`
	Code         string        `json:"code"`
	Name         string        `json:"name"`
	Desc         string        `json:"description"`
	IsActive     ProductStatus `json:"is_active"`
	Status       string        `json:"status"`
	CreatorId    *uint         `json:"creator_id"`
	CreatorName  string        `json:"creator_name"`
	ModifierId   *uint         `json:"modifier_id"`
	ModifierName string        `json:"modifier_name"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
}
