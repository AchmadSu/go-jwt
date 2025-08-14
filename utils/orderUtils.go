package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"example.com/m/dto"
	"example.com/m/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const OrderPrefix = "INV-"

func ToPublicOrder(order models.Order, details []models.OrderDetails) dto.PublicOrderWithDetail {
	return dto.PublicOrderWithDetail{
		ID:          order.ID,
		Code:        order.Code,
		DateEntry:   order.DateEntry,
		TotalQty:    int(order.TotalQty),
		GrandTotal:  dto.Total(order.GrandTotal),
		CreatorId:   &order.Creator.ID,
		CreatorName: order.Creator.Name,
		CreatedAt:   order.CreatedAt,
		Details:     ToPublicOrderDetail(details),
	}
}

func ToPublicOrderDetail(details []models.OrderDetails) []dto.OrderDetail {
	result := make([]dto.OrderDetail, 0, len(details))
	for _, detail := range details {
		result = append(result, dto.OrderDetail{
			ProductID: detail.Product.ID,
			Code:      detail.Product.Code,
			Name:      detail.Product.Name,
			Qty:       dto.StockQty(detail.Qty),
			UnitPrice: dto.StockPrice(detail.UnitPrice),
			Total:     dto.Total(detail.Total),
		})
	}
	return result
}

func IsEmptyOrder(o dto.PublicOrderWithDetail) bool {
	return o.ID == 0 && o.Code == ""
}

func formatOrderCode(t time.Time, seq int) string {
	return fmt.Sprintf("%s%s%04d", OrderPrefix, t.Format("200601"), seq)
}

func GenerateCodeOrder(tx *gorm.DB, now time.Time) (string, error) {
	prefixMonth := now.Format("200601")
	prefix := OrderPrefix + prefixMonth
	likePattern := prefix + "%"

	var lastCode string

	if err := tx.
		Model(&models.Order{}).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("code ILIKE ?", likePattern).
		Order("code DESC").
		Limit(1).
		Pluck("code", &lastCode).Error; err != nil {
		return "", err
	}

	seq := 1
	if lastCode != "" {
		re := regexp.MustCompile(`(\d{4})$`)
		m := re.FindStringSubmatch(lastCode)

		if len(m) == 2 {
			if n, err := strconv.Atoi(m[1]); err == nil {
				seq = n + 1
			}
		}
	}

	return formatOrderCode(now, seq), nil
}
