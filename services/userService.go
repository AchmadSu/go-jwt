package services

import (
	"strconv"

	"example.com/m/dto"
	"example.com/m/repositories"
	"example.com/m/services/token"
	"example.com/m/services/validator"
	"example.com/m/utils"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	GetUser(id, email string) (dto.PublicUser, error)
	GetAllUsers(paginate *dto.PaginationRequest) (*dto.PaginationResponse[dto.PublicUser], error)
	Register(input *dto.CreateUserInput) (dto.PublicUser, error)
	Login(input *dto.LoginUserInput) dto.LoginResult
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo}
}

func (s *userService) GetUser(id, email string) (dto.PublicUser, error) {
	if id != "" {
		parsedID, err := strconv.Atoi(id)
		if err != nil {
			return dto.PublicUser{}, err
		}
		user, result := s.repo.FindByID(parsedID)
		return utils.ToPublicUser(user), result.Error
	}
	user, result := s.repo.FindByEmail(email)
	return utils.ToPublicUser(user), result.Error
}

func (s *userService) GetAllUsers(paginate *dto.PaginationRequest) (*dto.PaginationResponse[dto.PublicUser], error) {
	pg, err := s.repo.FindAll(paginate)
	if err != nil {
		return nil, err
	}

	return pg, err
}

func (s *userService) Register(input *dto.CreateUserInput) (dto.PublicUser, error) {
	validator := validator.NewUserValidatorService(s.repo)
	isValid, err := validator.ValidateUserRegister(input.Email)
	if !isValid {
		return dto.PublicUser{}, err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), 10)
	if err != nil {
		return dto.PublicUser{}, err
	}
	input.Password = string(hash)
	return s.repo.Create(input)
}

func (s *userService) Login(input *dto.LoginUserInput) dto.LoginResult {
	validator := validator.NewUserValidatorService(s.repo)
	user, err := validator.ValidateUserLogin(input)
	if err != nil {
		return dto.LoginResult{Err: err}
	}
	token := token.NewJwtTokenService(s.repo)
	tokenString, exp, err := token.CreateToken(int(user.ID))
	if err != nil {
		return dto.LoginResult{Err: err}
	}
	return dto.LoginResult{User: user, Token: tokenString, Exp: exp, Err: nil}
}
