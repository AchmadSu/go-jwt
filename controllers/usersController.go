package controllers

import (
	"net/http"
	"os"
	"time"

	"example.com/m/initializers"
	"example.com/m/models"
	"example.com/m/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

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
	id := c.Query("id")
	email := c.Query("email")

	type PublicUser struct {
		ID        uint      `json:"id"`
		Email     string    `json:"email"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	if email != "" || id != "" {
		user, result := FindUserByIDOrEmail(id, email)
		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User not found",
			})
			return
		}
		publicUser := PublicUser{
			ID:        user.ID,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}
		c.JSON(http.StatusOK, gin.H{"user": publicUser})
	} else {
		pg, err := utils.Paginate(c, &models.User{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Pagination failed"})
			return
		}
		var users []models.User
		// initializers.DB.Find(&users)
		if err := initializers.DB.Limit(pg.Limit).Offset(pg.Offset).Find(&users).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
			return
		}

		publicUsers := make([]PublicUser, 0, len(users))
		for _, u := range users {
			publicUsers = append(publicUsers, PublicUser{
				ID:        u.ID,
				Email:     u.Email,
				CreatedAt: u.CreatedAt,
				UpdatedAt: u.UpdatedAt,
			})
		}

		// c.JSON(http.StatusOK, gin.H{"users": publicUsers})
		c.JSON(http.StatusOK, gin.H{
			"users":        publicUsers,
			"total":        pg.Total,
			"current_page": pg.Page,
			"total_pages":  pg.TotalPages,
		})
	}

}

func FindUserByIDOrEmail(id, email string) (models.User, *gorm.DB) {
	var user models.User
	var result *gorm.DB

	if id != "" {
		result = initializers.DB.First(&user, "id = ?", id)
	} else {
		result = initializers.DB.First(&user, "email = ?", email)
	}

	return user, result
}

func Validate(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "I'm logged in",
	})
}
