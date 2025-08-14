package models

import (
	"time"

	"gorm.io/gorm"
)

func (Order) TableName() string {
	return "orders"
}

type Order struct {
	gorm.Model
	Code         string    `gorm:"uniqueIndex;size:20;not null"`
	TotalQty     StockQty  `gorm:"type:int;default:0"`
	GrandTotal   Total     `gorm:"type:decimal(10,2);default=0"`
	DateEntry    time.Time `gorm:"not null"`
	OrderDetails []OrderDetails
	CreatedBy    *uint `gorm:"null"`
	ModifiedBy   *uint `gorm:"null"`
	Creator      User  `gorm:"foreignKey:CreatedBy;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`  //relation to user table
	Modifier     User  `gorm:"foreignKey:ModifiedBy;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"` //relation to user table                                                              // 1 = Active & 0 = Inactive
}
