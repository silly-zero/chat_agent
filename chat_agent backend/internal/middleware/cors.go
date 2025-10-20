package middleware

import (
	"github.com/gin-gonic/gin"
)

// CORSMiddleware 处理跨域请求的中间件
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 设置允许的源
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

		// 设置允许的HTTP方法
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")

		// 设置允许的请求头
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Accept, Cache-Control, X-Requested-With")

		// 设置是否允许发送Cookie
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		// 处理预检请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		// 继续处理请求
		c.Next()
	}
}