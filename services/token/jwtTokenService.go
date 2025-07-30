package token

import (
	"net/http"
	"os"
	"time"

	"example.com/m/config"
	"example.com/m/errs"
	"example.com/m/models"
	"example.com/m/repositories"
	"example.com/m/utils"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v4"
)

type JwtTokenService interface {
	CreateToken(sub int) (string, int, error)
	ValidateToken(tokenString string) (models.User, error)
	BlacklistToken(tokenString string, exp int64) error
}

type jwtTokenService struct {
	userRepo repositories.UserRepository
}

func NewJwtTokenService(userRepo repositories.UserRepository) *jwtTokenService {
	return &jwtTokenService{userRepo: userRepo}
}

func (j *jwtTokenService) CreateToken(sub int) (string, int, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": sub,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(), // token expire within a month
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_API_KEY")))
	if err != nil {
		return "", 0, err
	}
	exp := 3600 * 24 * 30
	return tokenString, exp, err
}

func (j *jwtTokenService) ValidateToken(tokenString string) (models.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			alg, _ := token.Header["alg"].(string)
			errorMessage := "Unexpected signing method: " + alg
			return nil, errs.New(errorMessage, http.StatusUnauthorized)
		}
		return []byte(os.Getenv("SECRET_API_KEY")), nil
	})
	if err != nil {
		return models.User{}, errs.New(utils.GetSafeErrorMessage(err, "Invalid Token Format"), http.StatusUnauthorized)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return models.User{}, errs.New("Invalid token claims", http.StatusUnauthorized)
	}

	expUnix, ok := claims["exp"].(float64)
	if !ok || float64(time.Now().Unix()) > expUnix {
		return models.User{}, errs.New("Token expired. Please re-login.", http.StatusUnauthorized)
	}

	val, err := config.RedisClient.Get(config.Ctx, tokenString).Result()
	if err == nil && val == "blacklisted" {
		return models.User{}, errs.New("Token is blacklisted. Please re-login.", http.StatusUnauthorized)
	}
	if err != nil && err != redis.Nil {
		return models.User{}, errs.New("Token validation failed. Redis error.", http.StatusInternalServerError)
	}

	subFloat, ok := claims["sub"].(float64)
	if !ok {
		return models.User{}, errs.New("Invalid token subject", http.StatusInternalServerError)
	}
	sub := int(subFloat)

	user, result := j.userRepo.FindByID(sub)
	if result.Error != nil {
		return models.User{}, errs.New(utils.GetSafeErrorMessage(result.Error, "Failed fetching user"), http.StatusInternalServerError)
	}

	return user, nil
}

func (j *jwtTokenService) BlacklistToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			alg, _ := token.Header["alg"].(string)
			errorMessage := "Unexpected signing method: " + alg
			return nil, errs.New(errorMessage, http.StatusUnauthorized)
		}
		return []byte(os.Getenv("SECRET_API_KEY")), nil
	})

	if err != nil {
		return errs.New(utils.GetSafeErrorMessage(err, "Invalid Token Format"), http.StatusUnauthorized)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return errs.New("Invalid token claims", http.StatusUnauthorized)
	}

	expUnix, ok := claims["exp"].(float64)
	if !ok {
		return errs.New("Failed get expired Token", http.StatusInternalServerError)
	}
	ttl := time.Until(time.Unix(int64(expUnix), 0))
	err = config.RedisClient.Set(config.Ctx, tokenString, "blacklisted", ttl).Err()
	if err != nil {
		return errs.New(utils.GetSafeErrorMessage(err, "Unknown Redis set blaclisted error"), http.StatusInternalServerError)
	}

	return nil
}
