package controllers

import (
	"net/http"
	"os"
	"time"

	"example.com/m/initializers"
	"example.com/m/models"
	"example.com/m/repositories"
	"example.com/m/services"
	"example.com/m/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type PublicUser struct {
	ID        uint      `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

var userService = services.NewUserService(repositories.NewUserRepository())

func SignUp(c *gin.Context) {
	var body struct {
		Email    string
		Password string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})

		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash password",
		})

		return
	}

	user := models.User{Email: body.Email, Password: string(hash)}
	result := initializers.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create user",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User has registered successfully",
	})
}

func Login(c *gin.Context) {
	var body struct {
		Email    string
		Password string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})

		return
	}

	var user models.User
	initializers.DB.First(&user, "email = ?", body.Email)

	if user.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found. Try another email",
		})

		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})

		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_API_KEY")))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create token",
		})

		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successfully!",
		"token":   tokenString,
	})
}

func GetUsers(c *gin.Context) {
	resp := utils.NewResponse()
	id := c.Query("id")
	email := c.Query("email")

	if email != "" || id != "" {
		user, result := userService.GetUser(id, email)

		if result.Error != nil {
			resp.SetStatus(http.StatusInternalServerError).
				SetMessage("Failed to get user data").
				SetError("Unknown error occurred").
				Send(c)
			return
		}

		if result.RowsAffected == 0 {
			resp.SetStatus(http.StatusNotFound).
				SetMessage("User not found").
				Send(c)
			return
		}

		publicUser := ToPublicUser(user)
		resp.SetMessage("Get user successfully").
			SetPayload(publicUser).
			Send(c)
		return
	}

	users, pg, err := userService.GetAllUsers(c)
	if err != nil {
		resp.SetStatus(http.StatusInternalServerError).
			SetMessage("Failed to fetch user data").
			SetError("Unknown error occurred").
			Send(c)
		return
	}

	publicUsers := make([]PublicUser, 0, len(users))
	for _, u := range users {
		publicUsers = append(publicUsers, ToPublicUser(u))
	}

	resp.SetMessage("Get users successfully").
		SetPayload(publicUsers).
		SetMeta(gin.H{
			"total":        pg.Total,
			"current_page": pg.Page,
			"total_pages":  pg.TotalPages,
		}).
		Send(c)
}

func ToPublicUser(u models.User) PublicUser {
	return PublicUser{
		ID:        u.ID,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func Validate(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "I'm logged in",
	})
}
