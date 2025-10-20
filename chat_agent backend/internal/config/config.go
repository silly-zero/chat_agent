package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config 应用程序配置
type Config struct {
	// 服务器配置
	ServerHost string
	ServerPort string

	// 数据库配置
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	// Redis配置
	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       int

	// 大语言模型配置
	OpenAIAPIKey string
	LLMModel     string
	LLMBaseURL   string
	LLMTimeout   int

	// 应用配置
	Environment string
}

// AppConfig 全局配置实例
var AppConfig *Config

// LoadConfig 从环境变量加载配置
func LoadConfig() *Config {
	// 加载.env文件
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	redisDB, _ := strconv.Atoi(getEnv("REDIS_DB", "0"))

	config := &Config{
		// 服务器配置
		ServerHost: getEnv("SERVER_HOST", "0.0.0.0"),
		ServerPort: getEnv("SERVER_PORT", "8000"),

		// 数据库配置
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", "password"),
		DBName:     getEnv("DB_NAME", "chat_agent"),

		// Redis配置
		RedisHost:     getEnv("REDIS_HOST", "localhost"),
		RedisPort:     getEnv("REDIS_PORT", "6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       redisDB,

		// 大语言模型配置
		OpenAIAPIKey: getEnv("LLM_API_KEY", ""),
		LLMModel:     getEnv("LLM_MODEL", "gpt-3.5-turbo"),
		LLMBaseURL:   getEnv("LLM_BASE_URL", "https://api.openai.com/v1"),
		LLMTimeout:   30,

		// 应用配置
		Environment: getEnv("GO_ENV", "development"),
	}

	AppConfig = config
	return config
}

// GetMySQLDSN 获取MySQL连接字符串
func (c *Config) GetMySQLDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
}

// GetRedisAddr 获取Redis地址
func (c *Config) GetRedisAddr() string {
	return fmt.Sprintf("%s:%s", c.RedisHost, c.RedisPort)
}

// IsProduction 判断是否为生产环境
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
