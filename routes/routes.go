package routes

import (
	"blog-api/controllers"
	"blog-api/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// CORS 配置
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

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
		// 文章路由
		protected.POST("/articles", articleController.CreateArticle)
		protected.PUT("/articles/:id", articleController.UpdateArticle)
		protected.DELETE("/articles/:id", articleController.DeleteArticle)

		// 评论路由
		protected.POST("/articles/:id/comments", commentController.CreateComment)
		protected.DELETE("/comments/:id", commentController.DeleteComment)
	}

	return r
}
