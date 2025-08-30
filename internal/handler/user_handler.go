package handler

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"strings"
	"time"

	"github.com/Incipe-win/ai-tshirt-shop/internal/model"
	"github.com/Incipe-win/ai-tshirt-shop/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6,max=100"`
	Email    string `json:"email" binding:"required,email,max=100"`
}

type RegisterResponse struct {
	ID           uint   `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	Message      string `json:"message"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	ID           uint   `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	Message      string `json:"message"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RefreshTokenResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	Message      string `json:"message"`
}

// generateRefreshToken generates a secure random refresh token
func generateRefreshToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// generateTokens generates both access and refresh tokens
func generateTokens(userID uint, username string) (string, string, error) {
	jwtSecret := viper.GetString("jwt.secret")
	if jwtSecret == "" {
		return "", "", jwt.ErrInvalidKey
	}

	// Generate access token (1 hour expiry)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":   userID,
		"username": username,
		"exp":      time.Now().Add(time.Hour).Unix(),
		"iat":      time.Now().Unix(),
	})

	accessToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", "", err
	}

	// Generate refresh token
	refreshToken, err := generateRefreshToken()
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// Register godoc
// @Summary 用户注册
// @Description 创建新用户账户
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "注册请求参数"
// @Success 201 {object} RegisterResponse "注册成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 409 {object} map[string]interface{} "用户名或邮箱已存在"
// @Failure 500 {object} map[string]interface{} "内部服务器错误"
// @Router /auth/register [post]
func Register(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	if len(req.Username) < 3 || len(req.Username) > 50 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "用户名长度必须在3到50个字符之间",
		})
		return
	}

	if len(req.Password) < 6 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "密码长度至少需要6个字符",
		})
		return
	}

	db := repository.GetDB()

	var existingUser model.User
	result := db.Where("username = ? OR email = ?", req.Username, req.Email).First(&existingUser)
	if result.Error == nil {
		if existingUser.Username == req.Username {
			c.JSON(http.StatusConflict, gin.H{
				"error": "用户名已存在",
			})
			return
		}
		if existingUser.Email == req.Email {
			c.JSON(http.StatusConflict, gin.H{
				"error": "邮箱已被注册",
			})
			return
		}
	} else if result.Error != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "数据库错误",
		})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "密码加密失败",
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
			"error": "创建用户失败",
		})
		return
	}

	// Generate JWT and refresh tokens for auto-login after registration
	accessToken, refreshToken, err := generateTokens(user.ID, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "生成token失败",
		})
		return
	}

	// Save refresh token to database
	user.RefreshToken = refreshToken
	if err := db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "保存refresh token失败",
		})
		return
	}

	response := RegisterResponse{
		ID:           user.ID,
		Username:     user.Username,
		Email:        user.Email,
		Token:        accessToken,
		RefreshToken: refreshToken,
		Message:      "User registered successfully",
	}

	c.JSON(http.StatusCreated, response)
}

// Login godoc
// @Summary 用户登录
// @Description 用户登录获取访问令牌
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "登录请求参数"
// @Success 200 {object} LoginResponse "登录成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 401 {object} map[string]interface{} "用户名或密码错误"
// @Failure 500 {object} map[string]interface{} "内部服务器错误"
// @Router /auth/login [post]
func Login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
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
				"error": "用户名不存在",
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
			"error": "密码错误",
		})
		return
	}

	// Generate JWT and refresh tokens
	accessToken, refreshToken, err := generateTokens(user.ID, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "生成token失败",
		})
		return
	}

	// Save refresh token to database
	user.RefreshToken = refreshToken
	if err := db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "保存refresh token失败",
		})
		return
	}

	response := LoginResponse{
		ID:           user.ID,
		Username:     user.Username,
		Email:        user.Email,
		Token:        accessToken,
		RefreshToken: refreshToken,
		Message:      "Login successful",
	}

	c.JSON(http.StatusOK, response)
}

// RefreshToken godoc
// @Summary 刷新令牌
// @Description 使用刷新令牌获取新的访问令牌
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RefreshTokenRequest true "刷新令牌请求参数"
// @Success 200 {object} RefreshTokenResponse "令牌刷新成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 401 {object} map[string]interface{} "非法刷新令牌"
// @Failure 500 {object} map[string]interface{} "内部服务器错误"
// @Router /auth/refresh [post]
func RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	db := repository.GetDB()

	var user model.User
	result := db.Where("refresh_token = ?", req.RefreshToken).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid refresh token",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database error",
		})
		return
	}

	// Generate new tokens
	accessToken, refreshToken, err := generateTokens(user.ID, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "生成token失败",
		})
		return
	}

	// Update refresh token in database
	user.RefreshToken = refreshToken
	if err := db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "保存refresh token失败",
		})
		return
	}

	response := RefreshTokenResponse{
		Token:        accessToken,
		RefreshToken: refreshToken,
		Message:      "Tokens refreshed successfully",
	}

	c.JSON(http.StatusOK, response)
}
