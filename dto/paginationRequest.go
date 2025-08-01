package dto

type PaginationRequest struct {
	Limit  string   `form:"limit"`
	Page   string   `form:"page"`
	SortBy []string `form:"sort_by[]"`
	Search string   `form:"search"`
}
