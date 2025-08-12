package utils

import (
	"fmt"
	"strconv"
	"time"

	"example.com/m/config"
	"example.com/m/dto"
	"gorm.io/gorm"
)

func CompareNumberQuery(query *gorm.DB, field string, operator string, value any) *gorm.DB {
	valStr, ok := value.(string)
	if !ok {
		return query
	}

	if number, err := strconv.Atoi(valStr); err == nil {
		return query.Where(fmt.Sprintf("%s %s ?", field, operator), number)
	}

	if number, err := strconv.ParseFloat(valStr, 64); err == nil {
		return query.Where(fmt.Sprintf("%s %s ?", field, operator), number)
	}

	return query
}

func OnDateQuery(query *gorm.DB, field string, value string) *gorm.DB {
	if value != "" {
		if _, err := time.ParseInLocation(string(config.LayoutDate), value, time.Local); err == nil {
			return query.Where(fmt.Sprintf("%s >= ?", field), value+" 00:00:00").Where(fmt.Sprintf("%s <= ?", field), value+" 23:59:59")
		}
	}
	return query
}

func BetweenDateQuery(query *gorm.DB, field string, value map[string]string) *gorm.DB {
	if value["start_date"] != "" && value["end_date"] != "" {
		if _, err := time.ParseInLocation(string(config.LayoutDate), value["start_date"], time.Local); err == nil {
			if _, err := time.ParseInLocation(string(config.LayoutDate), value["end_date"], time.Local); err == nil {
				start := value["start_date"] + " 00:00:00"
				end := value["end_date"] + " 23:59:59"
				return query.Where(fmt.Sprintf("%s >= ?", field), start).
					Where(fmt.Sprintf("%s <= ?", field), end)
			}
		}
	}
	return query
}

func ApplyDateRange(param map[string]any, query *gorm.DB, keyStart, keyEnd, dbField string) *gorm.DB {
	startDate, okStart := param[keyStart].(string)
	endDate, okEnd := param[keyEnd].(string)

	if okStart && okEnd && startDate != "" && endDate != "" {
		return BetweenDateQuery(query, dbField, map[string]string{
			"start_date": startDate,
			"end_date":   endDate,
		})
	}
	return query
}

func FilterQuery(param *dto.PaginationRequest, query *gorm.DB, mainTable string) *gorm.DB {
	if param.IsActive != nil {
		if *param.IsActive == 0 || *param.IsActive == 1 || *param.IsActive == 2 {
			query = query.Where(fmt.Sprintf("%s.is_active = ?", mainTable), *param.IsActive)
		}
	}

	if param.CreateDateStart != "" && param.CreateDateEnd != "" {
		dateMap := map[string]any{
			"create_date_start": param.CreateDateStart,
			"create_date_end":   param.CreateDateEnd,
		}
		query = ApplyDateRange(dateMap, query, "create_date_start", "create_date_end", fmt.Sprintf("%s.created_at", mainTable))
	}

	// Handle updated_at date range
	if param.UpdateDateStart != "" && param.UpdateDateEnd != "" {
		dateMap := map[string]any{
			"update_date_start": param.UpdateDateStart,
			"update_date_end":   param.UpdateDateEnd,
		}
		query = ApplyDateRange(dateMap, query, "update_date_start", "update_date_end", fmt.Sprintf("%s.updated_at", mainTable))
	}

	// Handle created_by
	if param.CreatorId != nil && *param.CreatorId > 0 {
		query = query.Where(fmt.Sprintf("%s.created_by = ?", mainTable), *param.CreatorId)
	}

	// Handle modified_by
	if param.ModifierId != nil && *param.ModifierId > 0 {
		query = query.Where(fmt.Sprintf("%s.modified_by = ?", mainTable), *param.ModifierId)
	}

	return query
}
