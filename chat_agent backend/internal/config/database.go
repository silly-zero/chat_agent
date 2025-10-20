package config

import (
	"fmt"
	"log"

	"chat_agent/internal/models"
	
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDatabase 初始化数据库连接
func InitDatabase(config *Config) (*gorm.DB, error) {
	// 配置GORM日志
	logLevel := logger.Info
	if config.IsProduction() {
		logLevel = logger.Error
	}

	// 首先尝试连接MySQL服务器（不指定数据库）
	dsnWithoutDB := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True&loc=Local",
		config.DBUser, config.DBPassword, config.DBHost, config.DBPort)
	
	tempDB, err := gorm.Open(mysql.Open(dsnWithoutDB), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // 静默模式，避免显示连接到空数据库的警告
	})
	
	if err == nil {
		// 尝试创建数据库
		log.Printf("尝试创建数据库 %s（如果不存在）...", config.DBName)
		sqlDB, _ := tempDB.DB()
		// 使用SQL直接创建数据库
		_, err := sqlDB.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", config.DBName))
		if err != nil {
			log.Printf("警告: 创建数据库 %s 失败: %v", config.DBName, err)
		} else {
			log.Printf("数据库 %s 已准备就绪", config.DBName)
		}
		// 关闭临时连接
		sqlDB.Close()
	}

	// 现在连接到指定的数据库
	db, err := gorm.Open(mysql.Open(config.GetMySQLDSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Connected to database successfully")
	DB = db
	return db, nil
}

// MigrateDatabase 执行数据库迁移
func MigrateDatabase(db *gorm.DB) error {
	log.Println("Starting database migration...")

	// 自动迁移数据模型
	err := db.AutoMigrate(
		&models.User{},
		&models.Star{},
		&models.Chat{},
		&models.Message{},
	)
	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("Database migration completed successfully")
	return nil
}

// SeedData 初始化种子数据
func SeedData(db *gorm.DB) error {
	log.Println("Seeding initial data...")

	// 检查是否已有用户数据，如果没有则创建默认用户
	var userCount int64
	db.Model(&models.User{}).Count(&userCount)
	if userCount == 0 {
		// 创建默认用户（ID为1）
		defaultUser := models.User{
			Username:  "default_user",
			Email:     "default@example.com",
			Password:  "password", // 实际应用中应该加密存储
			Nickname:  "默认用户",
			Avatar:    "",
			IsActive:  true,
		}

		result := db.Create(&defaultUser)
		if result.Error != nil {
			return fmt.Errorf("failed to seed default user: %w", result.Error)
		}
		log.Printf("Seeded default user with ID: %d", defaultUser.ID)
	}

	// 检查是否已有明星数据
	var starCount int64
	db.Model(&models.Star{}).Count(&starCount)
	if starCount == 0 {
		// 初始化一些示例明星数据
		stars := []models.Star{
			{
				Name:          "周杰伦",
				EnglishName:   "Jay Chou",
				Gender:        "男",
				BirthDate:     "1979-01-18",
				Nationality:   "中国台湾",
				Occupation:    "歌手、音乐人、演员、导演",
				Introduction:  "华语流行乐男歌手、音乐人、演员、导演，被誉为亚洲流行音乐天王。创作了众多经典歌曲，如《七里香》《青花瓷》《稻香》等。",
				StyleFeatures: "语气自信且带有幽默感，喜欢使用比喻和双关语，常引用自己的歌词或作品，说话节奏独特，喜欢在句子结尾加上'哎哟'、'不错哦'等口头禅，性格直率但又不失礼貌。",
				IsActive:      true,
			},
			{
				Name:          "刘亦菲",
				EnglishName:   "Crystal Liu",
				Gender:        "女",
				BirthDate:     "1987-08-25",
				Nationality:   "美国（华裔）",
				Occupation:    "演员、歌手",
				Introduction:  "华语影视女演员、歌手，因主演《神雕侠侣》小龙女一角而广受关注，后主演多部热门影视剧，并进军好莱坞。",
				StyleFeatures: "说话温柔优雅，语速适中，用词礼貌，喜欢使用委婉的表达方式，性格内敛但真诚，回答问题时会认真思考，偶尔流露出小女生的可爱一面。",
				IsActive:      true,
			},
		}

		result := db.Create(&stars)
		if result.Error != nil {
			return fmt.Errorf("failed to seed star data: %w", result.Error)
		}
		log.Printf("Seeded %d stars successfully", len(stars))
	}

	return nil
}