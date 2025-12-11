package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	JWTSecret  string
	ServerPort string
}

func LoadConfig() *Config {
	// 加载 .env 文件
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	config := &Config{
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		JWTSecret:  os.Getenv("JWT_SECRET"),
		ServerPort: getEnv("SERVER_PORT", "8080"), // 只有端口可以有合理的默认值
	}

	// 验证必需的环境变量
	if err := config.validate(); err != nil {
		log.Fatal("Configuration error: ", err)
	}

	return config
}

func (c *Config) validate() error {
	required := map[string]string{
		"DB_HOST":     c.DBHost,
		"DB_PORT":     c.DBPort,
		"DB_USER":     c.DBUser,
		"DB_PASSWORD": c.DBPassword,
		"DB_NAME":     c.DBName,
		"JWT_SECRET":  c.JWTSecret,
	}

	for key, value := range required {
		if value == "" {
			return fmt.Errorf("%s environment variable is required", key)
		}
	}

	// 检查JWT密钥长度
	if len(c.JWTSecret) < 32 {
		log.Println("Warning: JWT_SECRET is too short, should be at least 32 characters")
	}

	return nil
}

// 只有对于真正可有默认值的配置项使用这个函数
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
