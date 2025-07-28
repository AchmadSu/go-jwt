package utils

import (
	"math"
	"strconv"

	"example.com/m/dto"
	"gorm.io/gorm"
)

func Paginate[T any](request *dto.PaginationRequest, query *gorm.DB) (*dto.PaginationResponse[T], error) {
	page, _ := strconv.Atoi(request.Page)
	limit, _ := strconv.Atoi(request.Limit)

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	offset := (page - 1) * limit

	var results []T
	if err := query.Limit(limit).Offset(offset).Find(&results).Error; err != nil {
		return nil, err
	}

	return &dto.PaginationResponse[T]{
		Data:       results,
		Limit:      limit,
		Page:       page,
		Offset:     offset,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}
