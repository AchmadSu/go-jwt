package validator

import (
	"net/http"

	"example.com/m/dto"
	"example.com/m/errs"
	"example.com/m/repositories"
	"example.com/m/utils"
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

func (v *UserValidatorService) ValidateUserLogin(input *dto.LoginUserInput) (dto.PublicUser, error) {
	user, result := v.userRepo.FindByEmail(input.Email)
	if result.Error != nil {
		return dto.PublicUser{}, result.Error
	}

	if result.RowsAffected == 0 {
		return dto.PublicUser{}, errs.New("Email has not registered yet. Please Sign Up!", http.StatusNotFound)
	}

	isValidPass, _ := CheckPassword(user.Password, input.Password)
	if !isValidPass {
		return dto.PublicUser{}, errs.New("Password is not correct!", http.StatusUnauthorized)
	}

	return utils.ToPublicUser(user), nil
}

func CheckPassword(password string, input string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(password), []byte(input))

	if err != nil {
		return false, err
	}

	return true, nil
}
