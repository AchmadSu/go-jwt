package dto

type PaginationResponse[T any] struct {
	Data       []T
	Limit      int
	Page       int
	Offset     int
	Total      int64
	TotalPages int
}
