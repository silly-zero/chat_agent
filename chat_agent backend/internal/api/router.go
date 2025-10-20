package api

import (
	"chat_agent/internal/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SetupRouter 配置路由
func SetupRouter(
	chatHandler *ChatHandler,
	starHandler *StarHandler,
) *gin.Engine {
	// 创建Gin引擎
	r := gin.Default()

	// 添加中间件
	r.Use(middleware.CORSMiddleware())

	// 健康检查
	r.GET("/health", HealthCheck)

	// API路由组 - 移除用户相关路由，只保留聊天和明星相关路由
	api := r.Group("/api/v1")
	{
		// 注册简化的路由
		chatHandler.RegisterRoutes(api)
		starHandler.RegisterRoutes(api)
	}

	// 静态文件服务
	r.Static("/static", "./static")

	// 提供主页访问
	r.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})

	return r
}

// HealthCheck 健康检查
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "服务运行正常",
	})
}