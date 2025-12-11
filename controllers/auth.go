package controllers

import (
	"blog-api/database"
	"blog-api/middleware"
	"blog-api/models"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthController struct {
	DB *gorm.DB
}

func NewAuthController() *AuthController {
	return &AuthController{
		DB: database.GetDB(),
	}
}

// RegisterInput 注册输入结构体
type RegisterInput struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginInput 登录输入结构体
type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Register 用户注册
func (ac *AuthController) Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Printf("DEBUG: 注册输入 - 用户名: %s, 邮箱: %s, 密码长度: %d\n",
		input.Username, input.Email, len(input.Password))

	// 创建用户对象
	user := models.User{
		Username: input.Username,
		Email:    input.Email,
		Password: input.Password,
	}

	// 加密密码
	if err := user.HashPassword(); err != nil {
		fmt.Printf("DEBUG: 密码加密失败: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not hash password: " + err.Error()})
		return
	}

	// 创建用户
	if err := ac.DB.Create(&user).Error; err != nil {
		fmt.Printf("DEBUG: 创建用户失败: %v\n", err)
		if strings.Contains(err.Error(), "Duplicate entry") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Username or email already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

// Login 用户登录
func (ac *AuthController) Login(c *gin.Context) {
	var input LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Printf("DEBUG: 登录尝试 - 用户名: %s\n", input.Username)

	var user models.User
	if err := ac.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		fmt.Printf("DEBUG: 用户查找失败: %v\n", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	fmt.Printf("DEBUG: 找到用户 - ID: %d, 密码哈希长度: %d\n", user.ID, len(user.Password))

	if err := user.CheckPassword(input.Password); err != nil {
		fmt.Printf("DEBUG: 密码验证失败: %v\n", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	fmt.Printf("DEBUG: 密码验证成功\n")

	token, err := middleware.GenerateToken(user.ID)
	if err != nil {
		fmt.Printf("DEBUG: Token 生成失败: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}
