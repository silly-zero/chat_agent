package ai

import (
	"context"
	"fmt"
	"io"
	ark "github.com/sashabaranov/go-openai"
)

// LLMClient 大语言模型客户端接口
type LLMClient interface {
	GenerateResponse(ctx context.Context, messages []map[string]string, model string) (string, error)
	GenerateStreamResponse(ctx context.Context, messages []map[string]string, model string, callback func(string) error) error
}

// OpenAIClient 大语言模型客户端实现
type OpenAIClient struct {
	client  *ark.Client
	model   string
	baseURL string // 存储baseURL用于调试
}
// NewOpenAIClient 创建一个新的大语言模型客户端
func NewOpenAIClient(apiKey, baseURL string, model string) *OpenAIClient {
	config := ark.DefaultConfig(apiKey)
	config.BaseURL = baseURL
	client := ark.NewClientWithConfig(config)
	// 如果没有提供模型，使用默认值
	if model == "" {
		model = "doubao-1.5-pro-32k-250115"
	}
	return &OpenAIClient{
		client:  client,
		model:   model,
		baseURL: baseURL, // 保存baseURL用于调试
	}
}

// 转换消息格式为API所需格式
func convertMessages(messages []map[string]string) []ark.ChatCompletionMessage {
	var converted []ark.ChatCompletionMessage
	for _, msg := range messages {
		role := msg["role"]
		// 确保使用正确的角色字符串，直接使用标准角色名
		converted = append(converted, ark.ChatCompletionMessage{
			Role:    role, // 直接使用原始角色名称，豆包API应该支持标准角色名
			Content: msg["content"],
		})
	}
	return converted
}

// GenerateResponse 生成非流式响应
func (c *OpenAIClient) GenerateResponse(ctx context.Context, messages []map[string]string, model string) (string, error) {
	// 使用配置的模型或传入的模型
	useModel := c.model
	if model != "" {
		useModel = model
	}
	
	// 重要调试信息 - 确保日志中能看到这个
	fmt.Printf("[CRITICAL] 最终使用的模型名称: %s\n", useModel)
	fmt.Printf("[CRITICAL] 客户端配置的模型: %s\n", c.model)
	fmt.Printf("[CRITICAL] 传入的模型参数: %s\n", model)
	
	// 如果模型名称不是预期的豆包模型，记录警告
	if useModel != "doubao-1.5-pro-32k-250115" {
		fmt.Printf("[WARNING] 模型名称不是豆包模型: %s\n", useModel)
	}
	
	fmt.Printf("[DEBUG] 消息数量: %d\n", len(messages))
	
	// 转换消息格式
	convertedMessages := convertMessages(messages)
	
	// 创建聊天完成请求
	req := ark.ChatCompletionRequest{
		Model:    useModel,
		Messages: convertedMessages,
	}
	
	// 非常详细的调试信息
	fmt.Printf("[CRITICAL] 准备发送API请求\n")
	fmt.Printf("[CRITICAL] API基础URL: %s\n", c.baseURL)
	fmt.Printf("[CRITICAL] 使用的模型: %s\n", req.Model)
	fmt.Printf("[CRITICAL] 请求消息数量: %d\n", len(req.Messages))
	
	// 尝试手动构建完整的API URL进行调试
	if c.baseURL != "" {
		fullURL := c.baseURL + "/chat/completions"
		fmt.Printf("[CRITICAL] 完整API URL猜测: %s\n", fullURL)
	}
	
	fmt.Printf("[DEBUG] 发送请求: 模型=%s, 消息数量=%d\n", req.Model, len(req.Messages))
	resp, err := c.client.CreateChatCompletion(ctx, req)
	
	if err != nil {
		fmt.Printf("[ERROR] API调用失败: %v\n", err)
		// 返回友好的错误信息
		return "", fmt.Errorf("ChatCompletion error: %v", err)
	}
	
	fmt.Printf("[DEBUG] API调用成功，返回选择数量: %d\n", len(resp.Choices))
	
	// 返回响应内容
	if len(resp.Choices) > 0 && resp.Choices[0].Message.Content != "" {
		fmt.Printf("[DEBUG] 收到有效响应，内容长度: %d\n", len(resp.Choices[0].Message.Content))
		return resp.Choices[0].Message.Content, nil
	}
	
	fmt.Printf("[ERROR] 未收到响应内容\n")
	return "", fmt.Errorf("no response content received")
}

// GenerateStreamResponse 生成流式响应
func (c *OpenAIClient) GenerateStreamResponse(ctx context.Context, messages []map[string]string, model string, callback func(string) error) error {
	// 使用配置的模型或传入的模型
	useModel := c.model
	if model != "" {
		useModel = model
	}
	
	// 转换消息格式
	convertedMessages := convertMessages(messages)
	
	// 创建流式聊天完成请求
	stream, err := c.client.CreateChatCompletionStream(
		ctx,
		ark.ChatCompletionRequest{
			Model:    useModel,
			Messages: convertedMessages,
		},
	)
	
	if err != nil {
		return fmt.Errorf("stream chat error: %w", err)
	}
	defer stream.Close()
	
	// 处理流式响应
	for {
		recv, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return fmt.Errorf("Stream chat error: %w", err)
		}
		
		// 提取内容并调用回调函数
		if len(recv.Choices) > 0 && recv.Choices[0].Delta.Content != "" {
			if err := callback(recv.Choices[0].Delta.Content); err != nil {
				return fmt.Errorf("callback error: %w", err)
			}
		}
	}
}
