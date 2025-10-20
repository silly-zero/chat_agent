package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"chat_agent/internal/ai"
	"chat_agent/internal/models"
	"chat_agent/internal/repository"

	"gorm.io/gorm"
)

// ChatService 聊天服务接口
type ChatService interface {
	// 创建新的聊天会话
	CreateChat(ctx context.Context, userID, starID uint) (*models.ChatResponse, error)

	// 获取用户的聊天会话列表
	GetUserChats(ctx context.Context, userID uint, page, pageSize int) ([]models.ChatResponse, int64, error)

	// 获取聊天会话详情
	GetChatByID(ctx context.Context, userID, chatID uint) (*models.ChatResponse, error)

	// 获取或创建与特定明星的聊天会话
	GetOrCreateChatWithStar(ctx context.Context, userID, starID uint) (*models.ChatResponse, error)

	// 更新聊天会话信息
	UpdateChat(ctx context.Context, userID, chatID uint, updates map[string]interface{}) (*models.ChatResponse, error)

	// 删除聊天会话
	DeleteChat(ctx context.Context, userID, chatID uint) error

	// 发送消息
	SendMessage(ctx context.Context, userID uint, req *models.SendMessageRequest) (*models.MessageResponse, error)

	// 流式发送消息
	SendMessageStream(ctx context.Context, userID uint, req *models.SendMessageRequest) (<-chan string, <-chan error, error)

	// 获取聊天消息列表
	GetChatMessages(ctx context.Context, userID, chatID uint, query models.MessageListQuery) ([]models.MessageResponse, int64, error)

	// 删除消息
	DeleteMessage(ctx context.Context, userID, messageID uint) error
}

// ChatServiceImpl 聊天服务实现
type ChatServiceImpl struct {
	chatRepo     repository.ChatRepository
	messageRepo  repository.MessageRepository
	starRepo     repository.StarRepository
	llmClient    ai.LLMClient
	memoryManager ai.MemoryManager
	promptBuilder *ai.PromptTemplate
}

// NewChatService 创建新的聊天服务
func NewChatService(
	chatRepo repository.ChatRepository,
	messageRepo repository.MessageRepository,
	starRepo repository.StarRepository,
	llmClient ai.LLMClient,
	memoryManager ai.MemoryManager,
	promptBuilder *ai.PromptTemplate,
) ChatService {
	return &ChatServiceImpl{
		chatRepo:     chatRepo,
		messageRepo:  messageRepo,
		starRepo:     starRepo,
		llmClient:    llmClient,
		memoryManager: memoryManager,
		promptBuilder: promptBuilder,
	}
}

// CreateChat 创建新的聊天会话
func (s *ChatServiceImpl) CreateChat(ctx context.Context, userID, starID uint) (*models.ChatResponse, error) {
	// 检查明星是否存在且活跃
	star, err := s.starRepo.GetByID(ctx, starID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("明星不存在")
		}
		return nil, err
	}

	if !star.IsActive {
		return nil, errors.New("明星不可用")
	}

	// 检查是否已存在与该明星的聊天会话
	existingChat, err := s.chatRepo.GetUserStarChat(ctx, userID, starID)
	if err == nil && existingChat != nil {
		response := existingChat.ToChatResponse(true)
		return &response, nil
	}

	// 创建新的聊天会话
	chat := &models.Chat{
		UserID:      userID,
		StarID:      starID,
		Title:       fmt.Sprintf("与%s的聊天", star.Name),
		MessageCount: 0,
		LastActive:  time.Now(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// 保存到数据库
	if err := s.chatRepo.Create(ctx, chat); err != nil {
		return nil, err
	}

	// 预加载明星信息
	chat.Star = *star
	// 确保返回的响应包含正确的ID和所有必要字段
	response := chat.ToChatResponse(true)
	return &response, nil
}

// GetUserChats 获取用户的聊天会话列表
func (s *ChatServiceImpl) GetUserChats(ctx context.Context, userID uint, page, pageSize int) ([]models.ChatResponse, int64, error) {
	chats, total, err := s.chatRepo.GetUserChats(ctx, userID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]models.ChatResponse, len(chats))
	for i, chat := range chats {
		responses[i] = chat.ToChatResponse(false)
	}

	return responses, total, nil
}

// GetChatByID 获取聊天会话详情
func (s *ChatServiceImpl) GetChatByID(ctx context.Context, userID, chatID uint) (*models.ChatResponse, error) {
	chat, err := s.chatRepo.GetByID(ctx, chatID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("聊天会话不存在")
		}
		return nil, err
	}

	// 验证是否是用户自己的聊天会话
	if chat.UserID != userID {
		return nil, errors.New("无权访问此聊天会话")
	}

	response := chat.ToChatResponse(false)
	return &response, nil
}

// GetOrCreateChatWithStar 获取或创建与特定明星的聊天会话
func (s *ChatServiceImpl) GetOrCreateChatWithStar(ctx context.Context, userID, starID uint) (*models.ChatResponse, error) {
	// 尝试获取已存在的会话
	chat, err := s.chatRepo.GetUserStarChat(ctx, userID, starID)
	if err == nil && chat != nil {
		response := chat.ToChatResponse(true)
		return &response, nil
	}

	// 如果不存在，创建新会话
	return s.CreateChat(ctx, userID, starID)
}

// UpdateChat 更新聊天会话信息
func (s *ChatServiceImpl) UpdateChat(ctx context.Context, userID, chatID uint, updates map[string]interface{}) (*models.ChatResponse, error) {
	// 获取聊天会话
	chat, err := s.chatRepo.GetByID(ctx, chatID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("聊天会话不存在")
		}
		return nil, err
	}

	// 验证是否是用户自己的聊天会话
	if chat.UserID != userID {
		return nil, errors.New("无权修改此聊天会话")
	}

	// 更新字段
	if title, ok := updates["title"].(string); ok && title != "" {
		chat.Title = title
	}

	chat.UpdatedAt = time.Now()

	// 保存更新
	if err := s.chatRepo.Update(ctx, chat); err != nil {
		return nil, err
	}
	response := chat.ToChatResponse(false)
	return &response, nil
}

// DeleteChat 删除聊天会话
func (s *ChatServiceImpl) DeleteChat(ctx context.Context, userID, chatID uint) error {
	// 获取聊天会话
	chat, err := s.chatRepo.GetByID(ctx, chatID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("聊天会话不存在")
		}
		return err
	}

	// 验证是否是用户自己的聊天会话
	if chat.UserID != userID {
		return errors.New("无权删除此聊天会话")
	}

	// 删除会话（包括相关消息）
	return s.chatRepo.Delete(ctx, chatID)
}

// SendMessage 发送消息
func (s *ChatServiceImpl) SendMessage(ctx context.Context, userID uint, req *models.SendMessageRequest) (*models.MessageResponse, error) {
	// 获取聊天会话
	chat, err := s.chatRepo.GetByID(ctx, req.ChatID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("聊天会话不存在")
		}
		return nil, err
	}

	// 验证是否是用户自己的聊天会话
	if chat.UserID != userID {
		return nil, errors.New("无权在该聊天会话中发送消息")
	}

	// 获取明星信息
	star, err := s.starRepo.GetByID(ctx, chat.StarID)
	if err != nil {
		// 如果明星信息不存在，使用默认值
		star = &models.Star{
			ID:           1,
			Name:         "默认明星",
			Introduction: "这是一个默认明星",
			IsActive:     true,
		}
	}

	// 创建用户消息
	userMessage := &models.Message{
		ChatID:     req.ChatID,
		SenderID:   userID,
		SenderType: models.SenderTypeUser,
		Content:    req.Content,
		Status:     models.MessageStatusSent,
		CreatedAt:  time.Now(),
	}

	// 保存用户消息
	if err := s.messageRepo.Create(ctx, userMessage); err != nil {
		return nil, err
	}

	// 更新聊天会话信息
	if err := s.chatRepo.UpdateLastActive(ctx, req.ChatID, req.Content); err != nil {
		return nil, err
	}

	if err := s.chatRepo.IncrementMessageCount(ctx, req.ChatID); err != nil {
		return nil, err
	}

	// 尝试使用爬虫增强明星资料（异步执行，不阻塞主流程）
	go func() {
		enhancedStar, err := ai.EnhanceStarProfile(context.Background(), s.starRepo, star.ID)
		if err == nil {
			// 如果成功增强了明星资料，使用增强后的资料
			star = enhancedStar
		}
	}()

	// 获取最近的聊天记录作为上下文
	recentMessages, err := s.messageRepo.GetLastMessages(ctx, req.ChatID, 10) // 增加历史记录数量
	if err != nil {
		return nil, err
	}

	// 获取长期记忆
	longTermMemories, err := s.memoryManager.GetLongTermMemory(ctx, req.ChatID, 10) // 获取最近10条长期记忆
	if err != nil {
		// 记忆获取失败不影响主流程
		longTermMemories = []string{}
	}

	// 构建提示词
	messages := s.promptBuilder.BuildChatCompletionMessages(star, recentMessages, req.Content, longTermMemories)

	// 添加到记忆
	s.memoryManager.AddShortTermMemory(ctx, req.ChatID, req.Content)

	// 尝试调用LLM获取回复
	response, err := s.llmClient.GenerateResponse(ctx, messages, req.Model)
	if err != nil {
		return nil, err
	}

	// 创建AI回复消息
	aiMessage := &models.Message{
		ChatID:     req.ChatID,
		SenderID:   star.ID,
		SenderType: models.SenderTypeStar,
		Content:    response,
		Status:     models.MessageStatusSent,
		CreatedAt:  time.Now(),
	}

	// 保存AI回复消息
	if err := s.messageRepo.Create(ctx, aiMessage); err != nil {
		return nil, err
	}

	// 再次更新聊天会话信息
	if err := s.chatRepo.UpdateLastActive(ctx, req.ChatID, response); err != nil {
		return nil, err
	}

	if err := s.chatRepo.IncrementMessageCount(ctx, req.ChatID); err != nil {
		return nil, err
	}

	// 添加AI回复到记忆
	s.memoryManager.AddShortTermMemory(ctx, req.ChatID, response)
	
	// 提取对话中的关键信息，更新长期记忆
	s.memoryManager.AddLongTermMemory(ctx, req.ChatID, response, 1.0) // weight=1.0表示重要性一般

	// 返回AI回复消息
	aiMessageResponse := aiMessage.ToMessageResponse()
	return &aiMessageResponse, nil
}

// SendMessageStream 流式发送消息
func (s *ChatServiceImpl) SendMessageStream(ctx context.Context, userID uint, req *models.SendMessageRequest) (<-chan string, <-chan error, error) {
	// 获取聊天会话
	chat, err := s.chatRepo.GetByID(ctx, req.ChatID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, errors.New("聊天会话不存在")
		}
		return nil, nil, err
	}

	// 验证是否是用户自己的聊天会话
	if chat.UserID != userID {
		return nil, nil, errors.New("无权在该聊天会话中发送消息")
	}

	// 获取明星信息
	star, err := s.starRepo.GetByID(ctx, chat.StarID)
	if err != nil {
		return nil, nil, err
	}

	// 创建用户消息
	userMessage := &models.Message{
		ChatID:     req.ChatID,
		SenderID:   userID,
		SenderType: models.SenderTypeUser,
		Content:    req.Content,
		Status:     models.MessageStatusSent,
		CreatedAt:  time.Now(),
	}

	// 保存用户消息
	if err := s.messageRepo.Create(ctx, userMessage); err != nil {
		return nil, nil, err
	}

	// 更新聊天会话信息
	if err := s.chatRepo.UpdateLastActive(ctx, req.ChatID, req.Content); err != nil {
		return nil, nil, err
	}

	if err := s.chatRepo.IncrementMessageCount(ctx, req.ChatID); err != nil {
		return nil, nil, err
	}

	// 尝试使用爬虫增强明星资料（异步执行，不阻塞主流程）
	go func() {
		enhancedStar, err := ai.EnhanceStarProfile(context.Background(), s.starRepo, star.ID)
		if err == nil {
			// 如果成功增强了明星资料，使用增强后的资料
			star = enhancedStar
		}
	}()

	// 获取最近的聊天记录作为上下文
	recentMessages, err := s.messageRepo.GetLastMessages(ctx, req.ChatID, 10) // 增加历史记录数量
	if err != nil {
		return nil, nil, err
	}

	// 获取长期记忆
	longTermMemories, err := s.memoryManager.GetLongTermMemory(ctx, req.ChatID, 10) // 获取最近10条长期记忆
	if err != nil {
		// 记忆获取失败不影响主流程
		longTermMemories = []string{}
	}

	// 构建提示词
	messages := s.promptBuilder.BuildChatCompletionMessages(star, recentMessages, req.Content, longTermMemories)

	// 添加到记忆
	s.memoryManager.AddShortTermMemory(ctx, req.ChatID, req.Content)

	// 创建响应通道
	streamChan := make(chan string)
	errChan := make(chan error)

	// 调用LLM的流式生成功能
	go func() {
		defer close(streamChan)
		defer close(errChan)

		err := s.llmClient.GenerateStreamResponse(ctx, messages, req.Model, func(chunk string) error {
			streamChan <- chunk
			return nil
		})

		if err != nil {
			// 如果调用失败，发送默认回复
			defaultResponse := "你好！很高兴能和你聊天。虽然我的AI功能暂时无法使用，但我依然可以陪伴你。有什么想聊的吗？"
			streamChan <- defaultResponse
			// 不发送错误，而是使用默认回复
		}
	}()

	// 创建最终的响应通道
	responseChan := make(chan string)
	finalErrChan := make(chan error)

	// 收集完整回复内容
	var fullResponse string

	// 处理流式响应
	go func() {
		defer close(responseChan)
		defer close(finalErrChan)

		for {
			select {
			case chunk, ok := <-streamChan:
				if !ok {
					// 流式响应结束，保存AI回复
					aiMessage := &models.Message{
						ChatID:     req.ChatID,
						SenderID:   star.ID,
						SenderType: models.SenderTypeStar,
						Content:    fullResponse,
						Status:     models.MessageStatusSent,
						CreatedAt:  time.Now(),
					}

					// 保存AI回复消息
					if err := s.messageRepo.Create(context.Background(), aiMessage); err != nil {
						finalErrChan <- err
						return
					}

					// 更新聊天会话信息
					if err := s.chatRepo.UpdateLastActive(context.Background(), req.ChatID, fullResponse); err != nil {
						finalErrChan <- err
						return
					}

					if err := s.chatRepo.IncrementMessageCount(context.Background(), req.ChatID); err != nil {
						finalErrChan <- err
						return
					}

					// 添加AI回复到记忆
			s.memoryManager.AddShortTermMemory(context.Background(), req.ChatID, fullResponse)
			
			// 提取对话中的关键信息，更新长期记忆
			s.memoryManager.AddLongTermMemory(context.Background(), req.ChatID, fullResponse, 1.0) // weight=1.0表示重要性一般

					return
				}

				// 发送数据块到响应通道
				responseChan <- chunk
				fullResponse += chunk

			case err, ok := <-errChan:
				if !ok {
					return
				}
				finalErrChan <- err
				return

			case <-ctx.Done():
				finalErrChan <- ctx.Err()
				return
			}
		}
	}()

	return streamChan, errChan, nil
}

// GetChatMessages 获取聊天消息列表
func (s *ChatServiceImpl) GetChatMessages(ctx context.Context, userID, chatID uint, query models.MessageListQuery) ([]models.MessageResponse, int64, error) {
	// 验证聊天会话权限
	chat, err := s.chatRepo.GetByID(ctx, chatID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, errors.New("聊天会话不存在")
		}
		return nil, 0, err
	}

	if chat.UserID != userID {
		return nil, 0, errors.New("无权访问此聊天会话的消息")
	}

	// 获取消息列表
	messages, total, err := s.messageRepo.GetChatMessages(ctx, chatID, query)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]models.MessageResponse, len(messages))
	for i, message := range messages {
		responses[i] = message.ToMessageResponse()
	}

	return responses, total, nil
}

// DeleteMessage 删除消息
func (s *ChatServiceImpl) DeleteMessage(ctx context.Context, userID, messageID uint) error {
	// 获取消息
	message, err := s.messageRepo.GetByID(ctx, messageID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("消息不存在")
		}
		return err
	}

	// 获取聊天会话
	chat, err := s.chatRepo.GetByID(ctx, message.ChatID)
	if err != nil {
		return err
	}

	// 验证权限（只有用户自己发送的消息可以删除，且必须是用户自己的聊天会话）
	if chat.UserID != userID || message.SenderType != models.SenderTypeUser || message.SenderID != userID {
		return errors.New("无权删除此消息")
	}

	// 删除消息
	return s.messageRepo.Delete(ctx, messageID)
}