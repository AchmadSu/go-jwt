package dto

type CreateProductInput struct {
	Code string `json:"code" binding:"required,min=3"`
	Name string `json:"name" binding:"required,min=3"`
	Desc string `json:"desc" binding:"required,min=8"`
}

type UpdateProductInput struct {
	Code     string `json:"code" binding:"omitempty,min=3"`
	Name     string `json:"name" binding:"omitempty,min=3"`
	Desc     string `json:"desc" binding:"omitempty,min=8"`
	IsActive *int   `json:"is_active" binding:"number"`
}
