package validator

import (
	"net/http"

	"example.com/m/errs"
	"example.com/m/models"
	"example.com/m/repositories"
	"golang.org/x/crypto/bcrypt"
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

func (v *UserValidatorService) ValidateUserLogin(input models.LoginUserInput) (models.User, error) {
	user, result := v.userRepo.FindByEmail(input.Email)
	if result.Error != nil {
		return models.User{}, result.Error
	}

	if result.RowsAffected == 0 {
		return models.User{}, errs.New("Email has not registered yet. Please Sign Up!", http.StatusNotFound)
	}

	isValidPass, _ := CheckPassword(user.Password, input.Password)
	if !isValidPass {
		return models.User{}, errs.New("Password is not correct!", http.StatusUnauthorized)
	}

	return user, nil
}

func CheckPassword(password string, input string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(password), []byte(input))
	if err != nil {
		return false, err
	}

	return true, nil
}
