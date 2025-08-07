package models

import (
	"example.com/m/config"
	"gorm.io/gorm"
)

type UserStatus int

func (us UserStatus) String() string {
	switch us {
	case config.Active:
		return "Active"
	case config.Draft:
		return "Draft"
	default:
		return "Inactive"
	}
}

type User struct {
	gorm.Model
	Name     string
	Email    string     `gorm:"unique"`
	IsActive UserStatus `gorm:"type:int;default:1"`
	Password string
}
