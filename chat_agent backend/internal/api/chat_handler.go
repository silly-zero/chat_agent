package api

import (
	"fmt"
	"net/http"
	"strconv"

	"chat_agent/internal/models"
	"chat_agent/internal/service"

	"github.com/gin-gonic/gin"
)

// ChatHandler 聊天API处理器
type ChatHandler struct {
	chatService service.ChatService
}

// NewChatHandler 创建新的聊天API处理器
func NewChatHandler(chatService service.ChatService) *ChatHandler {
	return &ChatHandler{
		chatService: chatService,
	}
}

// CreateChat 创建聊天会话
func (h *ChatHandler) CreateChat(c *gin.Context) {
	// 使用固定用户ID 1，简化为无需登录的聊天
	userID := uint(1)

	var req struct {
		StarID uint `json:"star_id" binding:"required"`
	}

	// 绑定请求参数
	if err := c.ShouldBindJSON(&req); err != nil {
		ParamError(c, err)
		return
	}

	// 调用服务层创建聊天会话
	chat, err := h.chatService.CreateChat(c.Request.Context(), userID, req.StarID)
	if err != nil {
		BadRequest(c, err.Error())
		return
	}

	// 返回成功响应
	SuccessWithMessage(c, "创建会话成功", chat)
}

// GetUserChats 获取用户的聊天会话列表
func (h *ChatHandler) GetUserChats(c *gin.Context) {
	// 使用固定用户ID 1，简化为无需登录的聊天
	userID := uint(1)

	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	// 参数校验
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// 调用服务层获取聊天会话列表
	chats, total, err := h.chatService.GetUserChats(c.Request.Context(), userID, page, pageSize)
	if err != nil {
		ServerError(c, err)
		return
	}

	// 返回成功响应
	SuccessPagination(c, chats, total, page, pageSize)
}

// GetChatByID 获取聊天会话详情
func (h *ChatHandler) GetChatByID(c *gin.Context) {
	// 使用固定用户ID 1，简化为无需登录的聊天
	userID := uint(1)

	// 获取聊天ID
	chatID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		ParamError(c, err)
		return
	}

	// 调用服务层获取聊天会话详情
	chat, err := h.chatService.GetChatByID(c.Request.Context(), userID, uint(chatID))
	if err != nil {
		NotFound(c, err.Error())
		return
	}

	// 返回成功响应
	Success(c, chat)
}

// GetOrCreateChatWithStar 获取或创建与特定明星的聊天会话
func (h *ChatHandler) GetOrCreateChatWithStar(c *gin.Context) {
	// 使用固定用户ID 1，简化为无需登录的聊天
	userID := uint(1)

	// 获取明星ID
	starID, err := strconv.ParseUint(c.Query("star_id"), 10, 32)
	if err != nil {
		ParamError(c, err)
		return
	}

	// 调用服务层获取或创建聊天会话
	chat, err := h.chatService.GetOrCreateChatWithStar(c.Request.Context(), userID, uint(starID))
	if err != nil {
		BadRequest(c, err.Error())
		return
	}

	// 返回成功响应
	Success(c, chat)
}

// UpdateChat 更新聊天会话信息
func (h *ChatHandler) UpdateChat(c *gin.Context) {
	// 使用固定用户ID 1，简化为无需登录的聊天
	userID := uint(1)

	// 获取聊天ID
	chatID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		ParamError(c, err)
		return
	}

	var req struct {
		Title string `json:"title" binding:"required"`
	}

	// 绑定请求参数
	if err := c.ShouldBindJSON(&req); err != nil {
		ParamError(c, err)
		return
	}

	// 构建更新字段
	updates := map[string]interface{}{
		"title": req.Title,
	}

	// 调用服务层更新聊天会话信息
	chat, err := h.chatService.UpdateChat(c.Request.Context(), userID, uint(chatID), updates)
	if err != nil {
		BadRequest(c, err.Error())
		return
	}

	// 返回成功响应
	SuccessWithMessage(c, "更新成功", chat)
}

// DeleteChat 删除聊天会话
func (h *ChatHandler) DeleteChat(c *gin.Context) {
	// 使用固定用户ID 1，简化为无需登录的聊天
	userID := uint(1)

	// 获取聊天ID
	chatID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		ParamError(c, err)
		return
	}

	// 调用服务层删除聊天会话
	err = h.chatService.DeleteChat(c.Request.Context(), userID, uint(chatID))
	if err != nil {
		BadRequest(c, err.Error())
		return
	}

	// 返回成功响应
	SuccessWithMessage(c, "删除成功", nil)
}

// SendMessage 发送消息
func (h *ChatHandler) SendMessage(c *gin.Context) {
	// 使用固定用户ID 1，简化为无需登录的聊天
	userID := uint(1)
	
	// 添加详细日志
	fmt.Println("收到发送消息请求")

	var req models.SendMessageRequest

	// 绑定请求参数
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Printf("请求参数绑定失败: %v\n", err)
		ParamError(c, err)
		return
	}
	
	// 记录请求参数
	fmt.Printf("请求参数: ChatID=%d, Content=%s, Model=%s\n", req.ChatID, req.Content, req.Model)

	// 保留原始模型参数，如果未指定则使用默认豆包模型
	if req.Model == "" {
		req.Model = "doubao-1.5-pro-32k-250115"
	}
	fmt.Printf("使用模型: %s\n", req.Model)

	// 调用服务层发送消息
	fmt.Println("调用服务层发送消息...")
	responseMessage, err := h.chatService.SendMessage(c.Request.Context(), userID, &req)
	if err != nil {
		fmt.Printf("服务层处理失败: %v\n", err)
		ServerError(c, err)
		return
	}
	
	fmt.Printf("消息处理成功，响应ID: %d\n", responseMessage.ID)

	// 返回成功响应
	Success(c, responseMessage)
}

// SendMessageStream 流式发送消息
func (h *ChatHandler) SendMessageStream(c *gin.Context) {
	// 使用固定用户ID 1，简化为无需登录的聊天
	userID := uint(1)

	var req models.SendMessageRequest

	// 绑定请求参数
	if err := c.ShouldBindJSON(&req); err != nil {
		ParamError(c, err)
		return
	}

	// 保留原始模型参数，如果未指定则使用默认豆包模型
	if req.Model == "" {
		req.Model = "doubao-1.5-pro-32k-250115"
	}
	fmt.Printf("使用模型: %s\n", req.Model)
	

	// 调用服务层流式发送消息
	streamChan, errChan, err := h.chatService.SendMessageStream(c.Request.Context(), userID, &req)
	if err != nil {
		// 设置响应头
		c.Header("Content-Type", "text/event-stream")
		c.Header("Cache-Control", "no-cache")
		c.Header("Connection", "keep-alive")
		
		// 返回默认回复
		defaultResponse := "你好！很高兴能和你聊天。虽然我的AI功能暂时无法使用，但我依然可以陪伴你。有什么想聊的吗？"
		c.String(http.StatusOK, "data: "+defaultResponse+"\n\n")
		c.String(http.StatusOK, "data: [DONE]\n\n")
		c.Writer.Flush()
		return
	}

	// 设置响应头
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	// 处理流式响应
	for {
		select {
		case chunk, ok := <-streamChan:
			if !ok {
				// 流式响应结束
				c.String(http.StatusOK, "data: [DONE]\n\n")
				return
			}

			// 发送数据块
			c.String(http.StatusOK, "data: "+chunk+"\n\n")
			c.Writer.Flush()

		case err, ok := <-errChan:
			if !ok {
				return
			}

			// 发生错误
			c.String(http.StatusOK, "error: "+err.Error()+"\n\n")
			return

		case <-c.Request.Context().Done():
			// 客户端断开连接
			return
		}
	}
}

// GetChatMessages 获取聊天消息列表
func (h *ChatHandler) GetChatMessages(c *gin.Context) {
	// 使用固定用户ID 1，简化为无需登录的聊天
	userID := uint(1)

	// 获取聊天ID
	chatID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		ParamError(c, err)
		return
	}

	// 获取查询参数
	query := models.MessageListQuery{
		Page:     1,
		PageSize: 50,
	}

	// 绑定查询参数
	if page, err := strconv.Atoi(c.Query("page")); err == nil && page > 0 {
		query.Page = page
	}

	if pageSize, err := strconv.Atoi(c.Query("page_size")); err == nil && pageSize > 0 && pageSize <= 100 {
		query.PageSize = pageSize
	}

	if beforeID, err := strconv.ParseUint(c.Query("before_id"), 10, 32); err == nil {
		query.BeforeID = uint(beforeID)
	}

	if afterID, err := strconv.ParseUint(c.Query("after_id"), 10, 32); err == nil {
		query.AfterID = uint(afterID)
	}

	// 调用服务层获取聊天消息列表
	messages, total, err := h.chatService.GetChatMessages(c.Request.Context(), userID, uint(chatID), query)
	if err != nil {
		ServerError(c, err)
		return
	}

	// 返回成功响应
	SuccessPagination(c, messages, total, query.Page, query.PageSize)
}

// DeleteMessage 删除消息
func (h *ChatHandler) DeleteMessage(c *gin.Context) {
	// 使用固定用户ID 1，简化为无需登录的聊天
	userID := uint(1)

	// 获取消息ID
	messageID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		ParamError(c, err)
		return
	}

	// 调用服务层删除消息
	err = h.chatService.DeleteMessage(c.Request.Context(), userID, uint(messageID))
	if err != nil {
		BadRequest(c, err.Error())
		return
	}

	// 返回成功响应
	SuccessWithMessage(c, "删除成功", nil)
}

// RegisterRoutes 注册聊天相关路由
func (h *ChatHandler) RegisterRoutes(router *gin.RouterGroup) {
	chats := router.Group("/chats")
	// 移除认证中间件，允许直接访问
	{
		// 聊天会话相关路由
		chats.GET("", h.GetUserChats)
		chats.POST("", h.CreateChat)
		chats.GET("/star", h.GetOrCreateChatWithStar)
		chats.GET("/:id", h.GetChatByID)
		chats.PUT("/:id", h.UpdateChat)
		chats.DELETE("/:id", h.DeleteChat)

		// 消息相关路由
		chats.GET("/:id/messages", h.GetChatMessages)
		chats.POST("/messages", h.SendMessage)
		chats.POST("/messages/stream", h.SendMessageStream)
		chats.DELETE("/messages/:id", h.DeleteMessage)
	}
}
