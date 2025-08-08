package utils

import (
	"example.com/m/dto"
	"example.com/m/models"
)

func ToPublicUser(u models.User) dto.PublicUser {
	return dto.PublicUser{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		IsActive:  dto.UserStatus(u.IsActive),
		Status:    u.IsActive.String(),
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func IsEmptyUser(u dto.PublicUser) bool {
	return u.ID == 0 && u.Email == ""
}
