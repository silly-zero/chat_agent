package models

import (
	"time"

	"gorm.io/gorm"
)

// 消息发送者类型常量
const (
	SenderTypeUser   = "user"
	SenderTypeStar   = "star"
	SenderTypeSystem = "system"
)

// 消息状态常量
const (
	MessageStatusSending    = "sending"
	MessageStatusSent       = "sent"
	MessageStatusDelivered  = "delivered"
	MessageStatusRead       = "read"
	MessageStatusFailed     = "failed"
)

// 消息类型常量
const (
	MessageTypeText  = "text"
	MessageTypeImage = "image"
	MessageTypeVoice = "voice"
)

// Message 消息模型
type Message struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	ChatID    uint   `gorm:"not null;index" json:"chat_id"`
	SenderID  uint   `gorm:"index" json:"sender_id"` // 0表示系统/明星，非0表示用户ID
	SenderType string `gorm:"size:20;not null" json:"sender_type"` // "user", "star", "system"
	Content   string `gorm:"type:text;not null" json:"content"`
	MessageType string `gorm:"size:20;default:'text'" json:"message_type"` // "text", "image", "voice"
	Status    string `gorm:"size:20;default:'sent'" json:"status"` // "sending", "sent", "delivered", "read", "failed"

	// 关联关系
	Chat Chat `gorm:"foreignKey:ChatID" json:"chat,omitempty"`
}

// TableName 指定表名
func (Message) TableName() string {
	return "messages"
}

// MessageResponse 消息响应数据
type MessageResponse struct {
	ID          uint      `json:"id"`
	ChatID      uint      `json:"chat_id"`
	SenderID    uint      `json:"sender_id"`
	SenderType  string    `json:"sender_type"`
	Content     string    `json:"content"`
	MessageType string    `json:"message_type"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ToMessageResponse 转换为消息响应数据
func (m *Message) ToMessageResponse() MessageResponse {
	return MessageResponse{
		ID:          m.ID,
		ChatID:      m.ChatID,
		SenderID:    m.SenderID,
		SenderType:  m.SenderType,
		Content:     m.Content,
		MessageType: m.MessageType,
		Status:      m.Status,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

// SendMessageRequest 发送消息请求
type SendMessageRequest struct {
	ChatID      uint   `json:"chat_id" binding:"required"`
	Content     string `json:"content" binding:"required"`
	MessageType string `json:"message_type" binding:"omitempty,oneof=text image voice"`
	Model       string `json:"model" binding:"omitempty"`
}

// MessageListQuery 消息列表查询参数
type MessageListQuery struct {
	Page     int `form:"page,default=1" binding:"min=1"`
	PageSize int `form:"page_size,default=50" binding:"min=1,max=200"`
	BeforeID uint `form:"before_id" binding:"omitempty"` // 获取该ID之前的消息
	AfterID  uint `form:"after_id" binding:"omitempty"`  // 获取该ID之后的消息
}