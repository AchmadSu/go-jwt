package validator

import (
	"net/http"

	"example.com/m/errs"
	"example.com/m/repositories"
)

type UserValidatorService struct {
	userRepo repositories.UserRepository
}

func NewUserValidatorService(repo repositories.UserRepository) *UserValidatorService {
	return &UserValidatorService{userRepo: repo}
}

func (v *UserValidatorService) ValidateUserRegister(email string) (bool, error) {
	_, result := v.userRepo.FindByEmail(email)
	if result.Error != nil {
		return false, result.Error
	}

	if result.RowsAffected > 0 {
		return false, errs.New("Email is already exists. Please try another email!", http.StatusNotAcceptable)
	}

	return true, nil
}
