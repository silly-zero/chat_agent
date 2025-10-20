package models

import (
	"time"

	"gorm.io/gorm"
)

// Star 明星模型
type Star struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Name          string `gorm:"size:100;not null" json:"name"`
	EnglishName   string `gorm:"size:100" json:"english_name"`
	Gender        string `gorm:"size:10" json:"gender"`
	BirthDate     string `gorm:"size:20" json:"birth_date"`
	Nationality   string `gorm:"size:50" json:"nationality"`
	Occupation    string `gorm:"size:100" json:"occupation"`
	Avatar        string `gorm:"size:500" json:"avatar"`
	CoverImage    string `gorm:"size:500" json:"cover_image"`
	Introduction  string `gorm:"type:text" json:"introduction"`
	StyleFeatures string `gorm:"type:text" json:"style_features"` // 语言风格特征描述
	IsActive      bool   `gorm:"default:true" json:"is_active"`

	// 关联关系
	Chats []Chat `gorm:"foreignKey:StarID" json:"-"`
}

// TableName 指定表名
func (Star) TableName() string {
	return "stars"
}

// StarResponse 明星响应数据
type StarResponse struct {
	ID            uint      `json:"id"`
	Name          string    `json:"name"`
	EnglishName   string    `json:"english_name"`
	Gender        string    `json:"gender"`
	BirthDate     string    `json:"birth_date"`
	Nationality   string    `json:"nationality"`
	Occupation    string    `json:"occupation"`
	Avatar        string    `json:"avatar"`
	CoverImage    string    `json:"cover_image"`
	Introduction  string    `json:"introduction"`
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
}

// ToStarResponse 转换为明星响应数据
func (s *Star) ToStarResponse() StarResponse {
	return StarResponse{
		ID:            s.ID,
		Name:          s.Name,
		EnglishName:   s.EnglishName,
		Gender:        s.Gender,
		BirthDate:     s.BirthDate,
		Nationality:   s.Nationality,
		Occupation:    s.Occupation,
		Avatar:        s.Avatar,
		CoverImage:    s.CoverImage,
		Introduction:  s.Introduction,
		IsActive:      s.IsActive,
		CreatedAt:     s.CreatedAt,
	}
}

// CreateStarRequest 创建明星请求
type CreateStarRequest struct {
	Name          string `json:"name" binding:"required"`
	EnglishName   string `json:"english_name"`
	Gender        string `json:"gender"`
	BirthDate     string `json:"birth_date"`
	Nationality   string `json:"nationality"`
	Occupation    string `json:"occupation"`
	Avatar        string `json:"avatar"`
	CoverImage    string `json:"cover_image"`
	Introduction  string `json:"introduction"`
	StyleFeatures string `json:"style_features" binding:"required"`
}

// UpdateStarRequest 更新明星请求
type UpdateStarRequest struct {
	Name          string `json:"name"`
	EnglishName   string `json:"english_name"`
	Gender        string `json:"gender"`
	BirthDate     string `json:"birth_date"`
	Nationality   string `json:"nationality"`
	Occupation    string `json:"occupation"`
	Avatar        string `json:"avatar"`
	CoverImage    string `json:"cover_image"`
	Introduction  string `json:"introduction"`
	StyleFeatures string `json:"style_features"`
	IsActive      *bool  `json:"is_active"`
}