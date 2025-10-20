package models

import (
	"time"

	"gorm.io/gorm"
)

// Chat 聊天会话模型
type Chat struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	UserID      uint   `gorm:"not null;index" json:"user_id"`
	StarID      uint   `gorm:"not null;index" json:"star_id"`
	Title       string `gorm:"size:200" json:"title"` // 会话标题，可根据第一条消息生成
	LastMessage string `gorm:"size:500" json:"last_message"`
	LastActive  time.Time `json:"last_active"`
	MessageCount int      `gorm:"default:0" json:"message_count"`

	// 关联关系
	Star     Star      `gorm:"foreignKey:StarID" json:"star,omitempty"`
	Messages []Message `gorm:"foreignKey:ChatID" json:"messages,omitempty"`
}

// TableName 指定表名
func (Chat) TableName() string {
	return "chats"
}

// ChatResponse 聊天会话响应数据
type ChatResponse struct {
	ID           uint         `json:"id"`
	UserID       uint         `json:"user_id"`
	StarID       uint         `json:"star_id"`
	Title        string       `json:"title"`
	LastMessage  string       `json:"last_message"`
	LastActive   time.Time    `json:"last_active"`
	MessageCount int          `json:"message_count"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
	Star         StarResponse `json:"star,omitempty"`
}

// ToChatResponse 转换为聊天会话响应数据
func (c *Chat) ToChatResponse(withStar bool) ChatResponse {
	response := ChatResponse{
		ID:           c.ID,
		UserID:       c.UserID,
		StarID:       c.StarID,
		Title:        c.Title,
		LastMessage:  c.LastMessage,
		LastActive:   c.LastActive,
		MessageCount: c.MessageCount,
		CreatedAt:    c.CreatedAt,
		UpdatedAt:    c.UpdatedAt,
	}

	if withStar && c.Star.ID != 0 {
		response.Star = c.Star.ToStarResponse()
	}

	return response
}

// CreateChatRequest 创建聊天会话请求
type CreateChatRequest struct {
	StarID uint `json:"star_id" binding:"required"`
	Title  string `json:"title"`
}

// UpdateChatRequest 更新聊天会话请求
type UpdateChatRequest struct {
	Title string `json:"title"`
}