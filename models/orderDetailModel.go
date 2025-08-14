package models

import (
	"gorm.io/gorm"
)

type OrderID uint
type UnitPrice float64
type Total float64

func (OrderDetails) TableName() string {
	return "order_details"
}

type OrderDetails struct {
	gorm.Model
	OrderID   *uint     `gorm:"not null"`
	Order     Order     `gorm:"foreignKey:OrderID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ProductID *uint     `gorm:"not null"`
	Product   Product   `gorm:"foreignKey:ProductID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"` //relation to product table
	Qty       StockQty  `gorm:"type:int;default:0"`
	UnitPrice UnitPrice `gorm:"type:decimal(10,2);default=0"`
	Total     Total     `gorm:"type:decimal(10,2);default=0"`
}
