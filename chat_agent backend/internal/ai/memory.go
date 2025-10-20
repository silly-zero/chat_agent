package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// MemoryType 记忆类型
type MemoryType string

const (
	// ShortTermMemory 短期记忆（当前对话上下文）
	ShortTermMemory MemoryType = "short_term"
	// LongTermMemory 长期记忆（关键信息）
	LongTermMemory MemoryType = "long_term"
)

// MemoryItem 记忆项结构
type MemoryItem struct {
	ID        string    `json:"id"`
	ChatID    uint      `json:"chat_id"`
	Type      MemoryType `json:"type"`
	Content   string    `json:"content"`
	Weight    float64   `json:"weight"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// MemoryManager 记忆管理器接口
type MemoryManager interface {
	// 短期记忆管理
	AddShortTermMemory(ctx context.Context, chatID uint, content string) error
	GetShortTermMemory(ctx context.Context, chatID uint, limit int) ([]string, error)
	ClearShortTermMemory(ctx context.Context, chatID uint) error

	// 长期记忆管理
	AddLongTermMemory(ctx context.Context, chatID uint, content string, weight float64) error
	GetLongTermMemory(ctx context.Context, chatID uint, limit int) ([]string, error)
	UpdateMemoryWeight(ctx context.Context, memoryID string, weight float64) error

	// 记忆检索
	SearchMemory(ctx context.Context, chatID uint, query string, limit int) ([]string, error)
}

// InMemoryManager 内存实现的记忆管理器（简单版本）
type InMemoryManager struct {
	memories map[uint][]MemoryItem
}

// NewInMemoryManager 创建新的内存记忆管理器
func NewInMemoryManager() *InMemoryManager {
	return &InMemoryManager{
		memories: make(map[uint][]MemoryItem),
	}
}

// AddShortTermMemory 添加短期记忆
func (m *InMemoryManager) AddShortTermMemory(ctx context.Context, chatID uint, content string) error {
	item := MemoryItem{
		ID:        fmt.Sprintf("%d_%d", chatID, time.Now().UnixNano()),
		ChatID:    chatID,
		Type:      ShortTermMemory,
		Content:   content,
		Weight:    1.0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if _, exists := m.memories[chatID]; !exists {
		m.memories[chatID] = []MemoryItem{}
	}

	// 添加到开头（最新的在前）
	m.memories[chatID] = append([]MemoryItem{item}, m.memories[chatID]...)

	// 限制短期记忆数量（保持最近的10条）
	if len(m.memories[chatID]) > 50 {
		// 筛选出所有短期记忆
		var shortTermMemories []MemoryItem
		var otherMemories []MemoryItem
		for _, mem := range m.memories[chatID] {
			if mem.Type == ShortTermMemory {
				shortTermMemories = append(shortTermMemories, mem)
			} else {
				otherMemories = append(otherMemories, mem)
			}
		}

		// 只保留最近的10条短期记忆
		if len(shortTermMemories) > 10 {
			shortTermMemories = shortTermMemories[:10]
		}

		// 合并回所有记忆
		m.memories[chatID] = append(shortTermMemories, otherMemories...)
	}

	return nil
}

// GetShortTermMemory 获取短期记忆
func (m *InMemoryManager) GetShortTermMemory(ctx context.Context, chatID uint, limit int) ([]string, error) {
	if limit <= 0 {
		limit = 10 // 默认返回10条
	}

	memories, exists := m.memories[chatID]
	if !exists {
		return []string{}, nil
	}

	var result []string
	count := 0

	// 倒序遍历，获取最近的记忆
	for i := len(memories) - 1; i >= 0; i-- {
		if memories[i].Type == ShortTermMemory {
			result = append(result, memories[i].Content)
			count++
			if count >= limit {
				break
			}
		}
	}

	return result, nil
}

// ClearShortTermMemory 清除短期记忆
func (m *InMemoryManager) ClearShortTermMemory(ctx context.Context, chatID uint) error {
	memories, exists := m.memories[chatID]
	if !exists {
		return nil
	}

	// 保留非短期记忆
	var otherMemories []MemoryItem
	for _, mem := range memories {
		if mem.Type != ShortTermMemory {
			otherMemories = append(otherMemories, mem)
		}
	}

	m.memories[chatID] = otherMemories
	return nil
}

// AddLongTermMemory 添加长期记忆
func (m *InMemoryManager) AddLongTermMemory(ctx context.Context, chatID uint, content string, weight float64) error {
	// 检查是否已存在相似内容的记忆
	if m.isDuplicateMemory(chatID, content) {
		return nil // 避免重复添加
	}

	item := MemoryItem{
		ID:        fmt.Sprintf("%d_%d", chatID, time.Now().UnixNano()),
		ChatID:    chatID,
		Type:      LongTermMemory,
		Content:   content,
		Weight:    weight,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if _, exists := m.memories[chatID]; !exists {
		m.memories[chatID] = []MemoryItem{}
	}

	// 添加长期记忆（按权重排序）
	m.memories[chatID] = append(m.memories[chatID], item)

	// 限制长期记忆数量（最多50条）
	m.limitAndSortMemories(chatID)

	return nil
}

// GetLongTermMemory 获取长期记忆
func (m *InMemoryManager) GetLongTermMemory(ctx context.Context, chatID uint, limit int) ([]string, error) {
	if limit <= 0 {
		limit = 20 // 默认返回20条
	}

	memories, exists := m.memories[chatID]
	if !exists {
		return []string{}, nil
	}

	var longTermMemories []MemoryItem
	for _, mem := range memories {
		if mem.Type == LongTermMemory {
			longTermMemories = append(longTermMemories, mem)
		}
	}

	// 按权重排序（降序）
	m.sortMemoriesByWeight(longTermMemories)

	var result []string
	for i, mem := range longTermMemories {
		result = append(result, mem.Content)
		if i+1 >= limit {
			break
		}
	}

	return result, nil
}

// UpdateMemoryWeight 更新记忆权重
func (m *InMemoryManager) UpdateMemoryWeight(ctx context.Context, memoryID string, weight float64) error {
	// 在所有聊天中查找并更新记忆权重
	for chatID, memories := range m.memories {
		for i, mem := range memories {
			if mem.ID == memoryID {
				memories[i].Weight = weight
				memories[i].UpdatedAt = time.Now()
				m.memories[chatID] = memories

				// 重新排序和限制数量
				m.limitAndSortMemories(chatID)
				return nil
			}
		}
	}

	return fmt.Errorf("memory not found: %s", memoryID)
}

// SearchMemory 搜索记忆
func (m *InMemoryManager) SearchMemory(ctx context.Context, chatID uint, query string, limit int) ([]string, error) {
	if limit <= 0 {
		limit = 5 // 默认返回5条
	}

	memories, exists := m.memories[chatID]
	if !exists {
		return []string{}, nil
	}

	var matchingMemories []MemoryItem
	for _, mem := range memories {
		// 简单的字符串匹配，实际应用中可以使用更复杂的算法
		if contains(mem.Content, query) {
			matchingMemories = append(matchingMemories, mem)
		}
	}

	// 按权重和时间排序
	m.sortMemoriesByWeight(matchingMemories)

	var result []string
	for i, mem := range matchingMemories {
		result = append(result, mem.Content)
		if i+1 >= limit {
			break
		}
	}

	return result, nil
}

// isDuplicateMemory 检查是否存在相似内容的记忆
func (m *InMemoryManager) isDuplicateMemory(chatID uint, content string) bool {
	memories, exists := m.memories[chatID]
	if !exists {
		return false
	}

	for _, mem := range memories {
		if mem.Type == LongTermMemory && mem.Content == content {
			return true
		}
	}

	return false
}

// limitAndSortMemories 限制记忆数量并排序
func (m *InMemoryManager) limitAndSortMemories(chatID uint) {
	memories := m.memories[chatID]
	var longTermMemories []MemoryItem

	// 提取长期记忆
	for _, mem := range memories {
		if mem.Type == LongTermMemory {
			longTermMemories = append(longTermMemories, mem)
		}
	}

	// 按权重排序（降序）
	m.sortMemoriesByWeight(longTermMemories)

	// 限制长期记忆数量
	if len(longTermMemories) > 50 {
		longTermMemories = longTermMemories[:50]
	}

	// 重新构建记忆列表：短期记忆 + 排序后的长期记忆
	var newMemories []MemoryItem
	for _, mem := range memories {
		if mem.Type == ShortTermMemory {
			newMemories = append(newMemories, mem)
		}
	}
	newMemories = append(newMemories, longTermMemories...)

	m.memories[chatID] = newMemories
}

// sortMemoriesByWeight 按权重排序记忆（降序）
func (m *InMemoryManager) sortMemoriesByWeight(memories []MemoryItem) {
	// 简单的冒泡排序，实际应用中可以使用更高效的排序算法
	for i := 0; i < len(memories)-1; i++ {
		for j := 0; j < len(memories)-i-1; j++ {
			if memories[j].Weight < memories[j+1].Weight {
				memories[j], memories[j+1] = memories[j+1], memories[j]
			}
		}
	}
}

// contains 简单的字符串包含检查
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// SerializeMemory 序列化记忆数据
func SerializeMemory(memories []MemoryItem) (string, error) {
	data, err := json.Marshal(memories)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// DeserializeMemory 反序列化记忆数据
func DeserializeMemory(data string) ([]MemoryItem, error) {
	var memories []MemoryItem
	err := json.Unmarshal([]byte(data), &memories)
	if err != nil {
		return nil, err
	}
	return memories, nil
}