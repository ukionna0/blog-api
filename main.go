package main

import (
	"blog-api/database"
	"blog-api/models"
	"blog-api/routes"
	"log"
)

func main() {
	// 连接数据库
	database.Connect()

	// 自动迁移数据库表
	db := database.GetDB()
	err := db.AutoMigrate(&models.User{}, &models.Article{}, &models.Comment{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// 设置路由
	r := routes.SetupRouter()

	// 启动服务器
	log.Println("Server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
