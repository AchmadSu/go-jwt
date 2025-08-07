package models

import (
	"time"

	"example.com/m/config"
	"gorm.io/gorm"
)

type ProductId uint
type StockQty int
type StockPrice float64
type StockStatus int

func (Stock) TableName() string {
	return "stocks"
}

func (ss StockStatus) String() string {
	switch ss {
	case config.Active:
		return "Active"
	case config.Draft:
		return "Draft"
	default:
		return "Inactive"
	}
}

type Stock struct {
	gorm.Model
	ProductId  *uint       `gorm:"not null"`
	Product    Product     `gorm:"foreignKey:ProductID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"` //relation to product table
	Qty        StockQty    `gorm:"type:int;default:0"`
	Price      StockPrice  `gorm:"type:decimal(10,2);default=0"`
	DateEntry  time.Time   `gorm:"not null"`
	IsActive   StockStatus `gorm:"type:int;default:1"`
	CreatedBy  *uint       `gorm:"null"`
	ModifiedBy *uint       `gorm:"null"`
	Creator    User        `gorm:"foreignKey:CreatedBy;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`  //relation to user table
	Modifier   User        `gorm:"foreignKey:ModifiedBy;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"` //relation to user table                                                              // 1 = Active & 0 = Inactive
}
