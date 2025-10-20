package main

import (
	"context"
	"fmt"

	"chat_agent/internal/ai"
	"chat_agent/internal/models"
)

// MockStarRepository 模拟的明星仓库，用于测试
 type MockStarRepository struct{}

// GetByID 实现仓库接口
func (r *MockStarRepository) GetByID(ctx context.Context, id uint) (*models.Star, error) {
	return nil, nil
}

// Create 实现仓库接口
func (r *MockStarRepository) Create(ctx context.Context, star *models.Star) error {
	return nil
}

// GetAllActive 实现仓库接口
func (r *MockStarRepository) GetAllActive(ctx context.Context) ([]models.Star, error) {
	return nil, nil
}

// GetList 实现仓库接口
func (r *MockStarRepository) GetList(ctx context.Context, page, pageSize int) ([]models.Star, int64, error) {
	return nil, 0, nil
}

// Update 实现仓库接口
func (r *MockStarRepository) Update(ctx context.Context, star *models.Star) error {
	return nil
}

// Delete 实现仓库接口
func (r *MockStarRepository) Delete(ctx context.Context, id uint) error {
	return nil
}

// Search 实现仓库接口
func (r *MockStarRepository) Search(ctx context.Context, keyword string) ([]models.Star, error) {
	return nil, nil
}

func main() {
	// 创建上下文
	ctx := context.Background()

	// 创建模拟仓库
	mockRepo := &MockStarRepository{}

	// 创建一些测试明星数据
	stars := []models.Star{
		{
			Name:         "周杰伦",
			EnglishName:  "Jay Chou",
			Introduction: "华语流行歌手、音乐人",
		},
		{
			Name:         "杨幂",
			EnglishName:  "Yang Mi",
			Introduction: "中国女演员、制片人",
		},
	}

	// 打印明星列表
	fmt.Printf("测试明星数量: %d\n", len(stars))
	for _, star := range stars {
		fmt.Printf("- %s\n", star.Name)
	}

	// 初始化简单爬虫
	crawler := ai.NewSimpleSocialCrawler(mockRepo)

	// 为每个明星测试爬虫功能
	for i := range stars {
		star := &stars[i] // 获取指针
		fmt.Printf("\n正在处理明星: %s\n", star.Name)
		
		// 测试抓取明星内容
		content, err := crawler.CrawlStarContent(ctx, star)
		if err != nil {
			fmt.Printf("抓取明星内容失败: %v\n", err)
			continue
		}
		fmt.Printf("抓取到的内容数量: %d\n", len(content))
		if len(content) > 0 {
			fmt.Printf("第一条内容: %s\n", content[0])
		}
		
		// 测试更新风格特征
		err = crawler.UpdateStarStyleFeatures(star, content)
		if err != nil {
			fmt.Printf("更新风格特征失败: %v\n", err)
		} else {
			fmt.Println("风格特征更新成功")
		}
	}

	fmt.Println("\n测试完成！")
}