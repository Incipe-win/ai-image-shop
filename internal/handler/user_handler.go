package handler

import (
	"net/http"
	"strings"
	"time"

	"github.com/Incipe-win/ai-tshirt-shop/internal/model"
	"github.com/Incipe-win/ai-tshirt-shop/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6,max=100"`
	Email    string `json:"email" binding:"required,email,max=100"`
}

type RegisterResponse struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Token     string `json:"token"`
	Message   string `json:"message"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Token    string `json:"token"`
	Message  string `json:"message"`
}

func Register(c *gin.Context) {
	var req RegisterRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	if len(req.Username) < 3 || len(req.Username) > 50 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Username must be between 3 and 50 characters",
		})
		return
	}

	if len(req.Password) < 6 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Password must be at least 6 characters",
		})
		return
	}

	db := repository.GetDB()

	var existingUser model.User
	result := db.Where("username = ? OR email = ?", req.Username, req.Email).First(&existingUser)
	if result.Error == nil {
		if existingUser.Username == req.Username {
			c.JSON(http.StatusConflict, gin.H{
				"error": "Username already exists",
			})
			return
		}
		if existingUser.Email == req.Email {
			c.JSON(http.StatusConflict, gin.H{
				"error": "Email already exists",
			})
			return
		}
	} else if result.Error != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database error",
		})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	user := model.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Email:    req.Email,
	}

	result = db.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create user",
		})
		return
	}

	// Generate JWT token for auto-login after registration
	jwtSecret := viper.GetString("jwt.secret")
	if jwtSecret == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "JWT secret not configured",
		})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": user.ID,
		"username": user.Username,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate token",
		})
		return
	}

	response := RegisterResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Token:    tokenString,
		Message:  "User registered successfully",
	}

	c.JSON(http.StatusCreated, response)
}

func Login(c *gin.Context) {
	var req LoginRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	req.Username = strings.TrimSpace(req.Username)

	db := repository.GetDB()

	var user model.User
	result := db.Where("username = ?", req.Username).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid username or password",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database error",
		})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid username or password",
		})
		return
	}

	jwtSecret := viper.GetString("jwt.secret")
	if jwtSecret == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "JWT secret not configured",
		})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": user.ID,
		"username": user.Username,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate token",
		})
		return
	}

	response := LoginResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Token:    tokenString,
		Message:  "Login successful",
	}

	c.JSON(http.StatusOK, response)
}