package models

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Username  string `gorm:"size:50;uniqueIndex" json:"username"`
	Email     string `gorm:"size:100;uniqueIndex" json:"email"`
	Password  string `gorm:"size:200" json:"-"` // 密码不返回给前端
	Nickname  string `gorm:"size:50" json:"nickname"`
	Avatar    string `gorm:"size:500" json:"avatar"`
	IsActive  bool   `gorm:"default:true" json:"is_active"`

	// 关联关系
	Chats []Chat `gorm:"foreignKey:UserID" json:"chats,omitempty"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// UserResponse 用户响应数据
type UserResponse struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Nickname  string    `json:"nickname"`
	Avatar    string    `json:"avatar"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

// ToUserResponse 转换为响应数据
func (u *User) ToUserResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		Nickname:  u.Nickname,
		Avatar:    u.Avatar,
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt,
	}
}