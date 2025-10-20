package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response API响应结构
type Response struct {
	Code    int         `json:"code"`    // 响应码，200表示成功，非200表示错误
	Message string      `json:"message"` // 响应消息
	Data    interface{} `json:"data"`    // 响应数据
}

// Pagination 分页信息
type Pagination struct {
	Total       int64 `json:"total"`       // 总数据量
	Page        int   `json:"page"`        // 当前页码
	PageSize    int   `json:"page_size"`   // 每页大小
	TotalPages  int   `json:"total_pages"` // 总页数
}

// PaginatedResponse 分页响应
type PaginatedResponse struct {
	Code       int         `json:"code"`       // 响应码
	Message    string      `json:"message"`    // 响应消息
	Data       interface{} `json:"data"`       // 分页数据列表
	Pagination Pagination  `json:"pagination"` // 分页信息
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "success",
		Data:    data,
	})
}

// SuccessWithMessage 带自定义消息的成功响应
func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: message,
		Data:    data,
	})
}

// SuccessPagination 分页成功响应
func SuccessPagination(c *gin.Context, data interface{}, total int64, page, pageSize int) {
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, PaginatedResponse{
		Code:    200,
		Message: "success",
		Data:    data,
		Pagination: Pagination{
			Total:      total,
			Page:       page,
			PageSize:   pageSize,
			TotalPages: totalPages,
		},
	})
}

// Fail 失败响应
func Fail(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}

// FailWithError 带错误信息的失败响应
func FailWithError(c *gin.Context, code int, err error) {
	message := "操作失败"
	if err != nil {
		message = err.Error()
	}

	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}

// BadRequest 400错误响应
func BadRequest(c *gin.Context, message string) {
	Fail(c, 400, message)
}

// Unauthorized 401错误响应
func Unauthorized(c *gin.Context) {
	Fail(c, 401, "未授权访问")
}

// Forbidden 403错误响应
func Forbidden(c *gin.Context) {
	Fail(c, 403, "权限不足")
}

// NotFound 404错误响应
func NotFound(c *gin.Context, message string) {
	if message == "" {
		message = "资源不存在"
	}
	Fail(c, 404, message)
}

// ServerError 500错误响应
func ServerError(c *gin.Context, err error) {
	message := "服务器内部错误"
	if err != nil {
		message = err.Error()
	}
	Fail(c, 500, message)
}

// ParamError 参数错误响应
func ParamError(c *gin.Context, err error) {
	message := "参数错误"
	if err != nil {
		message = "参数错误: " + err.Error()
	}
	BadRequest(c, message)
}