package utils

import (
	"example.com/m/models"
)

func ToPublicUser(u models.User) models.PublicUser {
	return models.PublicUser{
		ID:        u.ID,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
