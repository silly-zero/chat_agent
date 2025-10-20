package ai

import (
	"fmt"
	"strings"

	"chat_agent/internal/models"
)

// PromptTemplate 提示词模板生成器
type PromptTemplate struct{}

// NewPromptTemplate 创建新的提示词模板生成器
func NewPromptTemplate() *PromptTemplate {
	return &PromptTemplate{}
}

// BuildSystemPrompt 构建系统提示词
func (p *PromptTemplate) BuildSystemPrompt(star *models.Star) string {
	return fmt.Sprintf(`你现在需要扮演%s（%s），请严格按照以下要求进行对话：

## 人物背景
- 姓名：%s
- 英文名：%s
- 性别：%s
- 出生日期：%s
- 国籍：%s
- 职业：%s
- 个人简介：%s

## 语言风格要求
%s

## 对话规则
1. 完全沉浸在%s的角色中，使用符合其身份和风格的语言
2. 回答要自然、真实，避免过于机械或模板化
3. 保持对话的连贯性，记得之前的对话内容和用户提到的重要信息
4. 不要泄露你是AI助手的身份，始终保持角色一致性
5. 回答要简洁明了，符合日常对话习惯
6. 可以使用适当的表情符号和口头禅，但要符合明星的风格特点
7. 如果遇到敏感或不合适的问题，可以委婉拒绝回答
8. 对于音乐相关的问题，可以适当引用你（角色）的作品或风格
9. 保持积极向上的态度，展现明星的魅力和个性

请以%s的身份开始对话。`,
		star.Name,
		star.EnglishName,
		star.Name,
		star.EnglishName,
		star.Gender,
		star.BirthDate,
		star.Nationality,
		star.Occupation,
		star.Introduction,
		star.StyleFeatures,
		star.Name,
		star.Name,
	)
}

// BuildUserPrompt 构建用户提示词（包含历史对话上下文）
func (p *PromptTemplate) BuildUserPrompt(messages []models.Message, currentMessage string) string {
	var historyBuilder strings.Builder

	// 只保留最近的对话历史，避免token过长
	maxHistoryLength := 10
	historyStart := 0
	if len(messages) > maxHistoryLength {
		historyStart = len(messages) - maxHistoryLength
	}

	// 添加历史对话上下文
	for i := historyStart; i < len(messages); i++ {
		msg := messages[i]
		sender := "用户"
		if msg.SenderType == "star" {
			sender = "你"
		}
		historyBuilder.WriteString(fmt.Sprintf("%s: %s\n", sender, msg.Content))
	}

	// 添加当前用户消息
	historyBuilder.WriteString(fmt.Sprintf("用户: %s", currentMessage))

	return historyBuilder.String()
}

// BuildMemoryPrompt 构建记忆增强提示词
func (p *PromptTemplate) BuildMemoryPrompt(star *models.Star, memories []string) string {
	if len(memories) == 0 {
		return ""
	}

	// 按重要性排序（这里简单按时间顺序，最新的在前面）
	memoriesText := strings.Join(memories, "\n- ")
	return fmt.Sprintf(`作为%s，请记住以下重要信息，并在对话中自然地引用和参考：
- %s

这些信息来自之前的对话，请将它们融入到你的回应中，保持自然流畅。`, star.Name, memoriesText)
}

// ExtractKeyInfo 从对话中提取关键信息（用于长期记忆）
func (p *PromptTemplate) ExtractKeyInfo(conversation string) []string {
	// 增强的关键信息提取逻辑
	keywords := []string{
		"名字", "称呼", "爱好", "喜欢", "讨厌", "生日", "年龄",
		"工作", "学习", "家乡", "城市", "朋友", "家人", "宠物",
		"梦想", "目标", "愿望", "希望", "期待", "经历", "故事",
		"习惯", "特长", "技能", "喜欢的歌", "喜欢的电影", "喜欢的食物",
	}

	var keyInfos []string
	lines := strings.Split(conversation, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "用户:") || strings.HasPrefix(line, "你:") {
			continue
		}

		// 查找用户信息的模式
		for _, keyword := range keywords {
			// 处理冒号和逗号的情况
			if strings.Contains(line, keyword) {
				// 提取包含关键词的整个句子作为信息
				if len(line) < 100 { // 避免过长的信息
					// 检查是否已经包含类似信息
					isDuplicate := false
					for _, existingInfo := range keyInfos {
						if strings.Contains(existingInfo, line) || strings.Contains(line, existingInfo) {
							isDuplicate = true
							break
						}
					}
					if !isDuplicate {
						keyInfos = append(keyInfos, line)
					}
				}
			}
		}
	}

	// 限制返回的关键信息数量
	maxMemories := 15
	if len(keyInfos) > maxMemories {
		keyInfos = keyInfos[:maxMemories]
	}

	return keyInfos
}

// BuildChatCompletionMessages 构建完整的聊天完成请求消息
func (p *PromptTemplate) BuildChatCompletionMessages(
	star *models.Star,
	messages []models.Message,
	currentMessage string,
	memories []string,
) []map[string]string {
	var completionMessages []map[string]string

	// 添加系统提示词
	systemPrompt := p.BuildSystemPrompt(star)
	completionMessages = append(completionMessages, map[string]string{
		"role":    "system",
		"content": systemPrompt,
	})

	// 添加记忆增强提示词
	memoryPrompt := p.BuildMemoryPrompt(star, memories)
	if memoryPrompt != "" {
		completionMessages = append(completionMessages, map[string]string{
			"role":    "system",
			"content": memoryPrompt,
		})
	}

	// 添加对话历史和当前消息
	historyPrompt := p.BuildUserPrompt(messages, currentMessage)
	completionMessages = append(completionMessages, map[string]string{
		"role":    "user",
		"content": historyPrompt,
	})

	return completionMessages
}