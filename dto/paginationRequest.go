package dto

type PaginationRequest struct {
	Limit string `form:"limit"`
	Page  string `form:"page"`
}
