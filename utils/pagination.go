package utils

import (
	"math"
	"strconv"

	"example.com/m/initializers"
	"github.com/gin-gonic/gin"
)

type Pagination struct {
	Limit      int
	Page       int
	Offset     int
	Total      int64
	TotalPages int
}

func Paginate(c *gin.Context, model interface{}) (*Pagination, error) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 1
	}

	var total int64
	if err := initializers.DB.Model(model).Count(&total).Error; err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	offset := (page - 1) * limit

	return &Pagination{
		Limit:      limit,
		Page:       page,
		Offset:     offset,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}
