package main

import (
	"log"

	"chat_agent/internal/ai"
	"chat_agent/internal/api"
	"chat_agent/internal/config"
	"chat_agent/internal/repository"
	"chat_agent/internal/service"

	"github.com/gin-gonic/gin"
)

func Main() {
	// 加载配置
	cfg := config.LoadConfig()

	// 设置Gin模式
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 初始化数据库连接
	db, err := config.InitDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 数据库迁移
	if err := config.MigrateDatabase(db); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// 初始化种子数据
	if err := config.SeedData(db); err != nil {
		log.Printf("Warning: Failed to seed data: %v", err)
	}

	// 初始化仓库
	starRepo := repository.NewStarRepository(db)
	chatRepo := repository.NewChatRepository(db)
	messageRepo := repository.NewMessageRepository(db)

	// 初始化AI组件
	promptBuilder := ai.NewPromptTemplate()
	llmClient := ai.NewOpenAIClient(cfg.OpenAIAPIKey, cfg.LLMBaseURL, cfg.LLMModel)
	memoryManager := ai.NewInMemoryManager() // 创建内存记忆管理器

	// 初始化服务
	starService := service.NewStarService(starRepo)
	chatService := service.NewChatService(chatRepo, messageRepo, starRepo, llmClient, memoryManager, promptBuilder)

	// 初始化API处理器
	chatHandler := api.NewChatHandler(chatService)
	starHandler := api.NewStarHandler(starService)

	// 设置路由
	router := api.SetupRouter(chatHandler, starHandler)

	// 启动服务器
	serverAddr := cfg.ServerHost + ":" + cfg.ServerPort
	log.Printf("Server starting on %s", serverAddr)
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// main 是命令行入口点
func main() {
	Main()
}
