package api

import (
	"context"
	"fmt"
	"strconv"

	"chat_agent/internal/models"
	"chat_agent/internal/service"

	"github.com/gin-gonic/gin"
)
// StarHandler 明星API处理器
type StarHandler struct {
	starService service.StarService
}

// NewStarHandler 创建新的明星API处理器
func NewStarHandler(starService service.StarService) *StarHandler {
	return &StarHandler{
		starService: starService,
	}
}

// GetStarByID 获取明星详情
func (h *StarHandler) GetStarByID(c *gin.Context) {
	// 获取明星ID
	starID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		ParamError(c, err)
		return
	}

	// 调用服务层获取明星详情
	star, err := h.starService.GetStarByID(c.Request.Context(), uint(starID))
	if err != nil {
		NotFound(c, err.Error())
		return
	}

	// 返回成功响应
	Success(c, star)
}

// GetAllActiveStars 获取所有活跃的明星
func (h *StarHandler) GetAllActiveStars(c *gin.Context) {
	// 调用服务层获取所有活跃的明星
	stars, err := h.starService.GetAllActiveStars(c.Request.Context())
	if err != nil {
		ServerError(c, err)
		return
	}

	// 返回成功响应
	Success(c, stars)
}

// GetStarList 分页获取明星列表
func (h *StarHandler) GetStarList(c *gin.Context) {
	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	// 调用服务层分页获取明星列表
	stars, total, err := h.starService.GetStarList(c.Request.Context(), page, pageSize)
	if err != nil {
		ServerError(c, err)
		return
	}

	// 返回成功响应
	SuccessPagination(c, stars, total, page, pageSize)
}

// SearchStars 搜索明星
func (h *StarHandler) SearchStars(c *gin.Context) {
	// 获取搜索关键词
	keyword := c.Query("keyword")
	if keyword == "" {
		BadRequest(c, "搜索关键词不能为空")
		return
	}

	// 调用服务层搜索明星
	stars, err := h.starService.SearchStars(c.Request.Context(), keyword)
	if err != nil {
		ServerError(c, err)
		return
	}

	// 返回成功响应
	Success(c, stars)
}

// CreateStar 创建明星（管理员功能）
func (h *StarHandler) CreateStar(c *gin.Context) {
	var req models.CreateStarRequest

	// 绑定请求参数
	if err := c.ShouldBindJSON(&req); err != nil {
		ParamError(c, err)
		return
	}

	// 调用服务层创建明星
	star, err := h.starService.CreateStar(c.Request.Context(), &req)
	if err != nil {
		ServerError(c, err)
		return
	}

	// 返回成功响应
	SuccessWithMessage(c, "创建明星成功", star)
}

// UpdateStar 更新明星信息（管理员功能）
func (h *StarHandler) UpdateStar(c *gin.Context) {
	// 获取明星ID
	starID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		ParamError(c, err)
		return
	}

	var req models.UpdateStarRequest

	// 绑定请求参数
	if err := c.ShouldBindJSON(&req); err != nil {
		ParamError(c, err)
		return
	}

	// 调用服务层更新明星信息
	star, err := h.starService.UpdateStar(c.Request.Context(), uint(starID), &req)
	if err != nil {
		ServerError(c, err)
		return
	}

	// 返回成功响应
	SuccessWithMessage(c, "更新成功", star)
}

// DeleteStar 删除明星（管理员功能）
func (h *StarHandler) DeleteStar(c *gin.Context) {
	// 获取明星ID
	starID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		ParamError(c, err)
		return
	}

	// 调用服务层删除明星
	err = h.starService.DeleteStar(c.Request.Context(), uint(starID))
	if err != nil {
		ServerError(c, err)
		return
	}

	// 返回成功响应
	SuccessWithMessage(c, "删除成功", nil)
}

// ToggleStarActive 切换明星活跃状态（管理员功能）
func (h *StarHandler) ToggleStarActive(c *gin.Context) {
	// 获取明星ID
	starID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		ParamError(c, err)
		return
	}

	var req struct {
		IsActive bool `json:"is_active" binding:"required"`
	}

	// 绑定请求参数
	if err := c.ShouldBindJSON(&req); err != nil {
		ParamError(c, err)
		return
	}

	// 调用服务层切换明星活跃状态
	err = h.starService.ToggleStarActive(c.Request.Context(), uint(starID), req.IsActive)
	if err != nil {
		ServerError(c, err)
		return
	}

	// 返回成功响应
	status := "激活"
	if !req.IsActive {
		status = "禁用"
	}
	SuccessWithMessage(c, status+"成功", nil)
}

// EnhanceStarProfile 手动触发爬虫增强明星资料
func (h *StarHandler) EnhanceStarProfile(c *gin.Context) {
	// 从路径参数中获取明星ID
	starID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		BadRequest(c, "无效的明星ID")
		return
	}

	// 异步调用爬虫增强明星资料
	go func() {
		// 创建新的context用于异步操作
		ctx := context.Background()
		
		// 调用service层的EnhanceStarProfile方法
		err := h.starService.EnhanceStarProfile(ctx, uint(starID))
		if err != nil {
			// 记录错误日志
			// 在实际项目中，应该有更完善的错误处理机制
			// 比如发送通知、写入日志文件等
			fmt.Printf("增强明星资料失败: %v\n", err)
		}
	}()

	// 立即返回响应，不等待爬虫完成
	Success(c, gin.H{"message": "已开始爬取明星资料，稍后完成"})
}

// RegisterRoutes 注册明星相关路由
func (h *StarHandler) RegisterRoutes(router *gin.RouterGroup) {
	stars := router.Group("/stars")
	{
		// 所有路由都直接访问，不使用认证中间件
		stars.GET("", h.GetAllActiveStars)
		stars.GET("/list", h.GetStarList)
		stars.GET("/search", h.SearchStars)
		stars.GET("/:id", h.GetStarByID)

		// 管理员功能也直接访问（演示版本）
		stars.POST("", h.CreateStar)
		stars.PUT("/:id", h.UpdateStar)
		stars.DELETE("/:id", h.DeleteStar)
		stars.PUT("/:id/active", h.ToggleStarActive)
		// 增强明星资料（爬虫）
		stars.POST("/:id/enhance", h.EnhanceStarProfile)
	}
}
