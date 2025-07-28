package utils

import (
	"example.com/m/dto"
	"example.com/m/models"
)

func ToPublicUser(u models.User) dto.PublicUser {
	return dto.PublicUser{
		ID:        u.ID,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func IsEmptyUser(u dto.PublicUser) bool {
	return u.ID == 0 && u.Email == ""
}
