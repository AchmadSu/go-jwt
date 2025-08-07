package dto

type PaginationRequest struct {
	ID              *int     `form:"id"`
	Name            string   `form:"name"`
	Code            string   `form:"code"`
	Email           string   `form:"email"`
	IsActive        *int     `form:"is_active"`
	Status          string   `form:"status"`
	CreatorId       *uint    `form:"creator_id"`
	ModifierId      *uint    `form:"modifier_id"`
	CreateDateStart string   `form:"create_date_start"`
	CreateDateEnd   string   `form:"create_date_end"`
	UpdateDateStart string   `form:"update_date_start"`
	UpdateDateEnd   string   `form:"update_date_end"`
	Limit           string   `form:"limit"`
	Page            string   `form:"page"`
	SortBy          []string `form:"sort_by[]"`
	Search          string   `form:"search"`
}
