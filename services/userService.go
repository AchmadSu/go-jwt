package services

import (
	"net/http"
	"strconv"

	"example.com/m/dto"
	"example.com/m/errs"
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
	Logout(tokenString string) error
}

type userService struct {
	repo      repositories.UserRepository
	validator validator.UserValidatorService
	token     token.JwtTokenService
}

func NewUserService(
	repo repositories.UserRepository,
	validator validator.UserValidatorService,
	token token.JwtTokenService,
) UserService {
	return &userService{
		repo:      repo,
		validator: validator,
		token:     token,
	}
}

func (s *userService) GetUser(id, email string) (dto.PublicUser, error) {
	var publicUser dto.PublicUser
	var errResult error

	if id != "" {
		parsedID, err := strconv.Atoi(id)
		if err != nil {
			return dto.PublicUser{}, errs.New("User ID is not a number!", http.StatusBadRequest)
		}
		user, result := s.repo.FindByID(parsedID)
		publicUser = utils.ToPublicUser(user)
		errResult = result.Error
	} else {
		user, result := s.repo.FindByEmail(email)
		publicUser = utils.ToPublicUser(user)
		errResult = result.Error
	}

	if utils.IsEmptyUser(publicUser) {
		return dto.PublicUser{}, errs.New("User not found", http.StatusNotFound)
	}

	return publicUser, errResult
}

func (s *userService) GetAllUsers(request *dto.PaginationRequest) (*dto.PaginationResponse[dto.PublicUser], error) {
	pg, err := s.repo.FindAll(request)
	if err != nil {
		return nil, err
	}

	return pg, err
}

func (s *userService) Register(input *dto.CreateUserInput) (dto.PublicUser, error) {
	isValid, err := s.validator.ValidateUserRegister(input.Email)
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
	user, err := s.validator.ValidateUserLogin(input)
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

func (s *userService) Logout(tokenString string) error {
	return s.token.BlacklistToken(tokenString)
}
