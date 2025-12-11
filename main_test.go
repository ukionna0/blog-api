package main

import (
	"blog-api/controllers"
	"blog-api/database"
	"blog-api/middleware"
	"blog-api/models"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var (
	router *gin.Engine
	token  string
)

func TestMain(m *testing.M) {
	// 设置测试环境
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "3306")
	os.Setenv("DB_USER", "blog-api")
	os.Setenv("DB_PASSWORD", "123456")
	os.Setenv("DB_NAME", "blog_test_db") // 使用测试数据库
	os.Setenv("JWT_SECRET", "test-secret-key-for-testing-only")

	// 初始化数据库
	database.Connect()

	// 迁移测试数据库
	db := database.GetDB()
	db.AutoMigrate(&models.User{}, &models.Article{}, &models.Comment{})

	// 设置路由
	router = setupRouter()

	// 运行测试
	code := m.Run()

	// 清理测试数据
	db.Exec("DELETE FROM comments")
	db.Exec("DELETE FROM articles")
	db.Exec("DELETE FROM users")

	os.Exit(code)
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	// CORS 配置
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Next()
	})

	authController := controllers.NewAuthController()
	articleController := controllers.NewArticleController()
	commentController := controllers.NewCommentController()

	// 公开路由
	public := r.Group("/api")
	{
		public.POST("/register", authController.Register)
		public.POST("/login", authController.Login)
		public.GET("/articles", articleController.GetArticles)
		public.GET("/articles/:id", articleController.GetArticle)
	}

	// 需要认证的路由
	protected := r.Group("/api")
	protected.Use(middleware.JWTAuth())
	{
		protected.POST("/articles", articleController.CreateArticle)
		protected.PUT("/articles/:id", articleController.UpdateArticle)
		protected.DELETE("/articles/:id", articleController.DeleteArticle)
		protected.POST("/articles/:id/comments", commentController.CreateComment)
		protected.DELETE("/comments/:id", commentController.DeleteComment)
	}

	return r
}

func TestRegisterAndLogin(t *testing.T) {
	// 注册用户
	registerData := map[string]string{
		"username": "testuser",
		"email":    "test@example.com",
		"password": "password123",
	}
	jsonData, _ := json.Marshal(registerData)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	// 用户登录
	loginData := map[string]string{
		"username": "testuser",
		"password": "password123",
	}
	jsonData, _ = json.Marshal(loginData)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// 保存 token 供后续测试使用
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	token = response["token"].(string)
}

func TestCreateArticle(t *testing.T) {
	articleData := map[string]string{
		"title":   "测试文章",
		"content": "测试文章内容",
	}
	jsonData, _ := json.Marshal(articleData)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/articles", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestGetArticles(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/articles", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	articles := response["articles"].([]interface{})
	assert.Greater(t, len(articles), 0)
}

func TestCreateComment(t *testing.T) {
	commentData := map[string]string{
		"content": "测试评论",
	}
	jsonData, _ := json.Marshal(commentData)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/articles/1/comments", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}
