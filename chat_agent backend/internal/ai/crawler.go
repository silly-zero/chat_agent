package ai

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"chat_agent/internal/models"
	"chat_agent/internal/repository"
)

// SocialMediaCrawler 社交媒体爬虫接口
type SocialMediaCrawler interface {
	// 抓取明星的社交媒体内容
	CrawlStarContent(ctx context.Context, star *models.Star) ([]string, error)
	// 更新明星的语言风格特征
	UpdateStarStyleFeatures(star *models.Star, contents []string) error
}

// SimpleSocialCrawler 简单的社交媒体爬虫实现
type SimpleSocialCrawler struct {
	starRepo repository.StarRepository
	client   *http.Client
}

// NewSimpleSocialCrawler 创建新的简单社交媒体爬虫
func NewSimpleSocialCrawler(starRepo repository.StarRepository) *SimpleSocialCrawler {
	return &SimpleSocialCrawler{
		starRepo: starRepo,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// RealSocialCrawler 真实的社交媒体爬虫实现
type RealSocialCrawler struct {
	starRepo    repository.StarRepository
	client      *http.Client
	userAgents  []string
	mutex       sync.Mutex
}

// NewRealSocialCrawler 创建新的真实社交媒体爬虫
func NewRealSocialCrawler(starRepo repository.StarRepository) *RealSocialCrawler {
	return &RealSocialCrawler{
		starRepo: starRepo,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		userAgents: []string{
			"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
			"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Firefox/89.0 Safari/537.36",
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
			"Mozilla/5.0 (iPhone; CPU iPhone OS 14_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.1.1 Mobile/15E148 Safari/604.1",
		},
	}
}

// EnhanceStarBasicInfo 增强明星基本资料
func (c *RealSocialCrawler) EnhanceStarBasicInfo(ctx context.Context, star *models.Star) error {
	return RealSocialCrawlerEnhanceStarBasicInfo(ctx, star)
}

// CrawlStarContent 抓取明星的社交媒体内容
func (c *RealSocialCrawler) CrawlStarContent(ctx context.Context, star *models.Star) ([]string, error) {
	return RealSocialCrawlerCrawlContent(ctx, star)
}

// UpdateStarStyleFeatures 根据抓取的内容更新明星的语言风格特征
func (c *RealSocialCrawler) UpdateStarStyleFeatures(star *models.Star, contents []string) error {
	return RealSocialCrawlerUpdateStyleFeatures(star, contents)
}

// MockSocialContent 模拟社交媒体内容（由于实际爬虫需要处理各种反爬措施，这里提供模拟数据）
func (c *SimpleSocialCrawler) CrawlStarContent(ctx context.Context, star *models.Star) ([]string, error) {
	// 在实际项目中，这里应该实现真实的爬虫逻辑
	// 由于环境限制，这里返回模拟数据

	// 根据明星名称返回模拟的社交媒体内容
	mockContents := make(map[string][]string)

	// 周杰伦的模拟社交媒体内容
	mockContents["周杰伦"] = []string{
		"今天在录音室创作新歌，灵感爆发，哎哟不错哦！#创作日常 #音乐",
		"想念大家了，下次演唱会见！等不及要和你们一起唱《七里香》了～",
		"收到很多朋友的生日祝福，谢谢你们一路以来的支持！爱你们！❤️",
		"最近在学新乐器，挑战自我的感觉真好！",
		"美食分享：今天做了一道台湾小吃，味道超赞！有空给大家分享食谱～",
		"电影拍摄中，这个角色很有挑战性，期待和大家见面！",
		"深夜创作，灵感总是在凌晨涌现，你们也是这样吗？",
		"刚刚看了粉丝的留言，真的很感动，谢谢你们的支持！",
	}

	// 更多明星的模拟内容可以在这里添加

	// 如果有该明星的模拟内容，则返回
	if contents, exists := mockContents[star.Name]; exists {
		return contents, nil
	}

	// 否则返回通用的模拟内容
	return []string{
		"今天的天气真好，心情也跟着好起来了！",
		"工作顺利完成，给自己点个赞！",
		"感谢粉丝们的支持，爱你们！",
		"分享今日的小确幸～",
		"期待接下来的工作计划！",
	}, nil
}

// UpdateStarStyleFeatures 根据抓取的内容更新明星的语言风格特征
func (c *SimpleSocialCrawler) UpdateStarStyleFeatures(star *models.Star, contents []string) error {
	// 分析社交媒体内容，提取语言风格特征
	styleAnalysis := ""

	// 计算内容特征
	contentLength := 0
	hasEmoji := false
	hasHashtag := false
	hasExclamation := false
	avgWordsPerSentence := 0

	if len(contents) > 0 {
		// 简单的特征分析
		totalWords := 0
		for _, content := range contents {
			contentLength += len(content)
			totalWords += len(strings.Fields(content))

			if strings.Contains(content, "#") {
				hasHashtag = true
			}
			if strings.Contains(content, "!") {
				hasExclamation = true
			}
			// 简单的表情符号检测
			if strings.Contains(content, "❤️") || strings.Contains(content, "😊") || strings.Contains(content, "🎉") {
				hasEmoji = true
			}
		}

		avgWordsPerSentence = totalWords / len(contents)

		// 构建风格分析
		styleAnalysis = "根据社交媒体分析：\n"
		if hasEmoji {
			styleAnalysis += "- 喜欢使用表情符号表达情感\n"
		}
		if hasHashtag {
			styleAnalysis += "- 常使用话题标签\n"
		}
		if hasExclamation {
			styleAnalysis += "- 语气活泼，常用感叹号\n"
		}
		if avgWordsPerSentence < 10 {
			styleAnalysis += "- 发言简洁明了\n"
		} else if avgWordsPerSentence > 20 {
			styleAnalysis += "- 喜欢详细表达想法\n"
		}
	}

	// 结合原有的风格特征
	if star.StyleFeatures != "" {
		star.StyleFeatures = star.StyleFeatures + "\n\n" + styleAnalysis
	} else {
		star.StyleFeatures = styleAnalysis
	}

	return nil
}

// EnhanceStarProfile 增强明星资料，通过爬虫获取更多信息
func EnhanceStarProfile(ctx context.Context, starRepo repository.StarRepository, starID uint) (*models.Star, error) {
	// 使用增强版的明星资料增强函数，使用真实爬虫
	return NewEnhancedStarProfile(ctx, starRepo, starID)
}

// NewEnhancedStarProfile 增强版的明星资料增强函数，使用真实爬虫
func NewEnhancedStarProfile(ctx context.Context, starRepo repository.StarRepository, starID uint) (*models.Star, error) {
	// 获取明星信息
	star, err := starRepo.GetByID(ctx, starID)
	if err != nil {
		return nil, fmt.Errorf("获取明星信息失败: %w", err)
	}

	// 创建真实爬虫
	crawler := NewRealSocialCrawler(starRepo)

	// 增强基本资料
	err = crawler.EnhanceStarBasicInfo(ctx, star)
	if err != nil {
		// 基本资料增强失败不影响后续流程
		fmt.Printf("增强明星基本资料失败: %v\n", err)
	}

	// 抓取社交媒体内容
	contents, err := crawler.CrawlStarContent(ctx, star)
	if err != nil {
		fmt.Printf("抓取明星社交媒体内容失败: %v\n", err)
		// 回退到使用简单爬虫
		fallbackCrawler := NewSimpleSocialCrawler(starRepo)
		contents, _ = fallbackCrawler.CrawlStarContent(ctx, star)
	}

	// 更新语言风格特征
	err = crawler.UpdateStarStyleFeatures(star, contents)
	if err != nil {
		fmt.Printf("更新明星语言风格特征失败: %v\n", err)
	} else {
		// 保存更新后的明星信息
		err = starRepo.Update(ctx, star)
		if err != nil {
			fmt.Printf("保存明星信息失败: %v\n", err)
		}
	}

	return star, nil
}
