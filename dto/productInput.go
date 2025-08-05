package dto

type CreateProductInput struct {
	Code string `json:"code" binding:"required,min=3"`
	Name string `json:"name" binding:"required,min=3"`
	Desc string `json:"desc" binding:"required,min=8"`
}
