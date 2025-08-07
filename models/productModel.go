package models

import (
	"example.com/m/config"
	"gorm.io/gorm"
)

type ProductStatus int

func (Product) TableName() string {
	return "products"
}

func (ps ProductStatus) String() string {
	switch ps {
	case config.Active:
		return "Active"
	case config.Draft:
		return "Draft"
	default:
		return "Inactive"
	}
}

type Product struct {
	gorm.Model
	Code       string        `gorm:"unique;size:20;not null"`
	Name       string        `gorm:"index;size:50;"`
	Desc       string        `gorm:"type:text"`
	IsActive   ProductStatus `gorm:"type:int;default:1"` // 1 = Active & 0 = Inactive
	CreatedBy  *uint         `gorm:"null"`
	ModifiedBy *uint         `gorm:"null"`
	Creator    User          `gorm:"foreignKey:CreatedBy;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`  //relation to user table
	Modifier   User          `gorm:"foreignKey:ModifiedBy;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"` //relation to user table
}
