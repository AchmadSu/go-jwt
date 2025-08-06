package utils

import (
	"math"
	"net/http"
	"strconv"
	"strings"

	"example.com/m/dto"
	"example.com/m/errs"
	"gorm.io/gorm"
)

func Paginate[T any](
	request *dto.PaginationRequest,
	query *gorm.DB,
	allowedSortFields []string,
	defaultOrder string,
	searchFields []string,
) (*dto.PaginationResponse[T], error) {
	page, pageErr := strconv.Atoi(request.Page)
	limit, limitErr := strconv.Atoi(request.Limit)
	if page < 1 || pageErr != nil {
		page = 1
	}

	if limit < 10 || limitErr != nil {
		limit = 10
	}

	if limit > 100 {
		limit = 100
	}

	var orConditions []string
	var args []any

	if request.Search != "" && len(searchFields) > 0 {
		searchTerm := "%" + request.Search + "%"
		for _, searchItem := range searchFields {
			orConditions = append(orConditions, searchItem+" ILIKE ?")
			args = append(args, searchTerm)
		}
		query = query.Where(strings.Join(orConditions, " OR "), args...)
	}

	if len(request.SortBy) > 0 && len(allowedSortFields) > 0 {
		for _, sortItem := range request.SortBy {
			parts := strings.Split(sortItem, ":")
			if len(parts) != 2 {
				continue
			}
			field, direction := parts[0], strings.ToLower(parts[1])
			if !ContainsString(allowedSortFields, field) {
				return nil, errs.New("Sort key '"+field+"' is not allowed", http.StatusBadRequest)
			}
			if direction != "asc" && direction != "desc" {
				direction = "asc"
			}
			query = query.Order(field + " " + direction)
		}
	} else if defaultOrder != "" {
		query = query.Order(defaultOrder)
	} else {
		query = query.Order("id asc")
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
