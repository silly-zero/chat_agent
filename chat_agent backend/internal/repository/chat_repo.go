package repository

import (
	"context"
	"time"

	"chat_agent/internal/models"

	"gorm.io/gorm"
)

// ChatRepository 聊天仓库接口
type ChatRepository interface {
	// 创建聊天会话
	Create(ctx context.Context, chat *models.Chat) error

	// 根据ID获取聊天会话
	GetByID(ctx context.Context, id uint) (*models.Chat, error)

	// 获取用户的聊天会话列表
	GetUserChats(ctx context.Context, userID uint, page, pageSize int) ([]models.Chat, int64, error)

	// 获取用户与特定明星的聊天会话
	GetUserStarChat(ctx context.Context, userID, starID uint) (*models.Chat, error)

	// 更新聊天会话
	Update(ctx context.Context, chat *models.Chat) error

	// 删除聊天会话
	Delete(ctx context.Context, id uint) error

	// 更新会话最后活动时间和最后消息
	UpdateLastActive(ctx context.Context, chatID uint, lastMessage string) error

	// 更新消息计数
	IncrementMessageCount(ctx context.Context, chatID uint) error
}

// ChatRepositoryImpl 聊天仓库实现
type ChatRepositoryImpl struct {
	db *gorm.DB
}

// NewChatRepository 创建新的聊天仓库
func NewChatRepository(db *gorm.DB) ChatRepository {
	return &ChatRepositoryImpl{db: db}
}

// Create 创建聊天会话
func (r *ChatRepositoryImpl) Create(ctx context.Context, chat *models.Chat) error {
	return r.db.WithContext(ctx).Create(chat).Error
}

// GetByID 根据ID获取聊天会话
func (r *ChatRepositoryImpl) GetByID(ctx context.Context, id uint) (*models.Chat, error) {
	var chat models.Chat
	err := r.db.WithContext(ctx).First(&chat, id).Error
	if err != nil {
		return nil, err
	}
	return &chat, nil
}

// GetUserChats 获取用户的聊天会话列表
func (r *ChatRepositoryImpl) GetUserChats(ctx context.Context, userID uint, page, pageSize int) ([]models.Chat, int64, error) {
	var chats []models.Chat
	var total int64

	// 计算偏移量
	offset := (page - 1) * pageSize

	// 获取总数
	if err := r.db.WithContext(ctx).Model(&models.Chat{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据，按最后活动时间倒序排列，并预加载明星信息
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Preload("Star").
		Order("last_active DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&chats).Error

	if err != nil {
		return nil, 0, err
	}

	return chats, total, nil
}

// GetUserStarChat 获取用户与特定明星的聊天会话
func (r *ChatRepositoryImpl) GetUserStarChat(ctx context.Context, userID, starID uint) (*models.Chat, error) {
	var chat models.Chat
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND star_id = ?", userID, starID).
		Preload("Star").
		First(&chat).Error

	if err != nil {
		return nil, err
	}

	return &chat, nil
}

// Update 更新聊天会话
func (r *ChatRepositoryImpl) Update(ctx context.Context, chat *models.Chat) error {
	return r.db.WithContext(ctx).Save(chat).Error
}

// Delete 删除聊天会话
func (r *ChatRepositoryImpl) Delete(ctx context.Context, id uint) error {
	// 使用事务确保数据一致性
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 先删除相关消息
		if err := tx.Where("chat_id = ?", id).Delete(&models.Message{}).Error; err != nil {
			return err
		}

		// 再删除聊天会话
		if err := tx.Delete(&models.Chat{}, id).Error; err != nil {
			return err
		}

		return nil
	})
}

// UpdateLastActive 更新会话最后活动时间和最后消息
func (r *ChatRepositoryImpl) UpdateLastActive(ctx context.Context, chatID uint, lastMessage string) error {
	return r.db.WithContext(ctx).Model(&models.Chat{}).Where("id = ?", chatID).Updates(map[string]interface{}{
		"last_message": lastMessage,
		"last_active":  time.Now(),
	}).Error
}

// IncrementMessageCount 更新消息计数
func (r *ChatRepositoryImpl) IncrementMessageCount(ctx context.Context, chatID uint) error {
	return r.db.WithContext(ctx).Model(&models.Chat{}).Where("id = ?", chatID).UpdateColumn("message_count", gorm.Expr("message_count + ?", 1)).Error
}

// MessageRepository 消息仓库接口
type MessageRepository interface {
	// 创建消息
	Create(ctx context.Context, message *models.Message) error

	// 批量创建消息
	CreateBatch(ctx context.Context, messages []models.Message) error

	// 根据ID获取消息
	GetByID(ctx context.Context, id uint) (*models.Message, error)

	// 获取聊天会话的消息列表
	GetChatMessages(ctx context.Context, chatID uint, query models.MessageListQuery) ([]models.Message, int64, error)

	// 更新消息状态
	UpdateStatus(ctx context.Context, messageID uint, status string) error

	// 删除消息
	Delete(ctx context.Context, id uint) error

	// 获取会话的最后几条消息
	GetLastMessages(ctx context.Context, chatID uint, limit int) ([]models.Message, error)
}

// MessageRepositoryImpl 消息仓库实现
type MessageRepositoryImpl struct {
	db *gorm.DB
}

// NewMessageRepository 创建新的消息仓库
func NewMessageRepository(db *gorm.DB) MessageRepository {
	return &MessageRepositoryImpl{db: db}
}

// Create 创建消息
func (r *MessageRepositoryImpl) Create(ctx context.Context, message *models.Message) error {
	return r.db.WithContext(ctx).Create(message).Error
}

// CreateBatch 批量创建消息
func (r *MessageRepositoryImpl) CreateBatch(ctx context.Context, messages []models.Message) error {
	return r.db.WithContext(ctx).CreateInBatches(messages, 100).Error
}

// GetByID 根据ID获取消息
func (r *MessageRepositoryImpl) GetByID(ctx context.Context, id uint) (*models.Message, error) {
	var message models.Message
	err := r.db.WithContext(ctx).First(&message, id).Error
	if err != nil {
		return nil, err
	}
	return &message, nil
}

// GetChatMessages 获取聊天会话的消息列表
func (r *MessageRepositoryImpl) GetChatMessages(ctx context.Context, chatID uint, query models.MessageListQuery) ([]models.Message, int64, error) {
	var messages []models.Message
	var total int64

	// 计算偏移量
	offset := (query.Page - 1) * query.PageSize

	db := r.db.WithContext(ctx).Model(&models.Message{}).Where("chat_id = ?", chatID)

	// 处理分页查询条件
	if query.BeforeID > 0 {
		db = db.Where("id < ?", query.BeforeID)
	}
	if query.AfterID > 0 {
		db = db.Where("id > ?", query.AfterID)
	}

	// 获取总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据，按创建时间倒序排列
	err := db.Order("created_at DESC").Offset(offset).Limit(query.PageSize).Find(&messages).Error
	if err != nil {
		return nil, 0, err
	}

	return messages, total, nil
}

// UpdateStatus 更新消息状态
func (r *MessageRepositoryImpl) UpdateStatus(ctx context.Context, messageID uint, status string) error {
	return r.db.WithContext(ctx).Model(&models.Message{}).Where("id = ?", messageID).Update("status", status).Error
}

// Delete 删除消息
func (r *MessageRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Message{}, id).Error
}

// GetLastMessages 获取会话的最后几条消息
func (r *MessageRepositoryImpl) GetLastMessages(ctx context.Context, chatID uint, limit int) ([]models.Message, error) {
	var messages []models.Message
	err := r.db.WithContext(ctx).
		Where("chat_id = ?", chatID).
		Order("created_at DESC").
		Limit(limit).
		Find(&messages).Error

	if err != nil {
		return nil, err
	}

	return messages, nil
}