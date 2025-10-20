package ai

import (
	"context"
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"sync"
	"time"

	"chat_agent/internal/models"
)

// 社交媒体内容结构体
type SocialContent struct {
	Text      string
	Timestamp time.Time
	Likes     int
	Source    string
}

// RealSocialCrawlerEnhanceStarBasicInfo 增强明星基本资料实现
func RealSocialCrawlerEnhanceStarBasicInfo(ctx context.Context, star *models.Star) error {
	// 这里可以实现从各种来源获取明星基本资料的逻辑
	// 例如从百科、官方网站等获取详细信息

	// 示例：从模拟的百科API获取信息
	enhancedInfo, err := fetchEncyclopediaInfo(star.Name)
	if err != nil {
		fmt.Printf("获取百科信息失败: %v\n", err)
		return err
	}

	// 更新明星信息
	if enhancedInfo.BirthDate != "" {
		star.BirthDate = enhancedInfo.BirthDate
	}
	if enhancedInfo.Nationality != "" {
		star.Nationality = enhancedInfo.Nationality
	}
	if enhancedInfo.Occupation != "" {
		star.Occupation = enhancedInfo.Occupation
	}
	if enhancedInfo.Biography != "" && len(enhancedInfo.Biography) > len(star.Introduction) {
		star.Introduction = enhancedInfo.Biography
	}

	return nil
}

// 百科信息结构体
type EncyclopediaInfo struct {
	BirthDate   string
	Nationality string
	Occupation  string
	Biography   string
}

// fetchEncyclopediaInfo 模拟从百科获取明星信息
func fetchEncyclopediaInfo(starName string) (*EncyclopediaInfo, error) {
	// 在实际项目中，这里应该调用真实的百科API或抓取百科网页
	// 这里提供模拟实现

	// 模拟网络延迟
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)

	// 模拟不同明星的百科信息
	switch strings.TrimSpace(starName) {
	case "周杰伦":
		return &EncyclopediaInfo{
			BirthDate:   "1979-01-18",
			Nationality: "中国台湾",
			Occupation:  "歌手、词曲作家、音乐制作人、演员、导演",
			Biography:   "周杰伦（Jay Chou），1979年1月18日出生于台湾省新北市，华语流行乐男歌手、音乐人、演员、导演、编剧、监制、商人。2000年发行首张个人专辑《Jay》。2001年发行专辑《范特西》。2002年举行The One世界巡回演唱会。2003年登上美国《时代周刊》封面。2004年获得世界音乐大奖中国区最畅销艺人奖。2005年凭借动作片《头文字D》获得台湾电影金马奖、香港电影金像奖最佳新人奖。2006年起连续三年获得世界音乐大奖中国区最畅销艺人奖。2007年自编自导的文艺片《不能说的秘密》获得台湾电影金马奖年度台湾杰出电影奖。",
		}, nil
	case "杨幂":
		return &EncyclopediaInfo{
			BirthDate:   "1986-09-12",
			Nationality: "中国内地",
			Occupation:  "演员、歌手、制片人",
			Biography:   "杨幂，1986年9月12日出生于北京市，中国内地影视女演员、流行乐歌手、影视制片人。2005年，杨幂进入北京电影学院表演系本科班就读。2006年，因出演金庸武侠剧《神雕侠侣》而崭露头角。2008年，凭借古装剧《王昭君》获得第24届中国电视金鹰奖观众喜爱的电视剧女演员奖提名。2009年，杨幂在\"80后新生代娱乐大明星\"评选活动中与其她三位女演员共同被评为\"四小花旦\"。2011年，因主演穿越剧《宫锁心玉》赢得广泛关注。",
		}, nil
	default:
		// 返回通用信息
		return &EncyclopediaInfo{
			Biography: "暂无详细资料",
		}, nil
	}
}

// RealSocialCrawlerCrawlContent 抓取明星的社交媒体内容实现
func RealSocialCrawlerCrawlContent(ctx context.Context, star *models.Star) ([]string, error) {
	// 随机选择一个User-Agent
	var mutex sync.Mutex
	userAgents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Firefox/89.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
	}
	mutex.Lock()
	userAgent := userAgents[rand.Intn(len(userAgents))]
	mutex.Unlock()

	var allContents []string
	var wg sync.WaitGroup
	var mu sync.Mutex
	var fetchErr error

	// 模拟从多个社交媒体平台抓取内容
	platforms := []string{"weibo", "douyin", "xiaohongshu"}

	for _, platform := range platforms {
		wg.Add(1)
		go func(p string) {
			defer wg.Done()
			
			// 模拟网络延迟
			time.Sleep(time.Duration(rand.Intn(2000)) * time.Millisecond)
			
			contents, err := fetchFromPlatform(ctx, star.Name, p, userAgent)
			if err != nil {
				fmt.Printf("从平台 %s 抓取内容失败: %v\n", p, err)
				fetchErr = err
				return
			}
			
			mu.Lock()
			allContents = append(allContents, contents...)
			mu.Unlock()
		}(platform)
	}

	wg.Wait()

	// 如果抓取失败，返回一些基础的模拟数据
	if len(allContents) == 0 && fetchErr != nil {
		return []string{
			"感谢大家的支持！",
			"最近在忙着新作品的创作，敬请期待！",
			"很高兴认识大家！",
		}, nil
	}

	return allContents, nil
}

// fetchFromPlatform 从特定平台抓取内容（模拟环境下直接返回模拟数据）
func fetchFromPlatform(ctx context.Context, starName, platform, userAgent string) ([]string, error) {
	// 在模拟环境中，直接返回模拟数据而不发送实际请求
	// 添加随机延迟模拟网络延迟
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	
	// 直接返回模拟数据
	return generateMockContentForPlatform(starName, platform), nil
}

// generateMockContentForPlatform 为不同平台生成模拟内容
func generateMockContentForPlatform(starName, platform string) []string {
	// 为不同平台生成不同风格的内容
	platformContents := make(map[string]map[string][]string)

	// 微博内容
	platformContents["weibo"] = map[string][]string{
		"周杰伦": {
			"新歌《最伟大的作品》MV明天发布，敬请期待！#周杰伦新歌 #新歌预告",
			"感谢粉丝们的支持，演唱会门票已售罄！下一站，北京！❤️",
			"今天在录音室和朋友们一起创作，灵感源源不断！#音乐创作 #工作室日常",
			"分享一张旧照片，怀念那段美好的时光～",
			"刚刚完成了新专辑的最后一首歌，感觉很棒！",
		},
		"杨幂": {
			"新剧开机啦！这是一个非常有挑战性的角色，期待与大家分享！🎬",
			"今天的拍摄很顺利，感谢所有工作人员的努力！#剧组日常 #拍戏日常",
			"分享今日穿搭，你们喜欢吗？💕 #时尚穿搭 #日常穿搭",
			"感谢粉丝们的生日礼物，很感动！爱你们！❤️",
			"参加品牌活动，遇见了很多老朋友，开心！",
		},
	}

	// 抖音内容
	platformContents["douyin"] = map[string][]string{
		"周杰伦": {
			"分享一段新歌的创作片段，猜猜这首歌叫什么名字？",
			"今天教大家弹钢琴，这个和弦你们学会了吗？",
			"和粉丝们互动的瞬间，谢谢你们的支持！",
			"分享今日的美食，自己做的哦！#美食分享",
			"工作室的日常，看看我们在忙什么？",
		},
		"杨幂": {
			"新剧的造型曝光，你们喜欢吗？#新剧造型",
			"今天的舞蹈练习，给大家跳一段！#舞蹈挑战",
			"和剧组小伙伴的搞笑日常，笑到停不下来！",
			"护肤小技巧分享，你们有什么好方法吗？#护肤分享",
			"健身打卡，坚持就是胜利！#健身日常",
		},
	}

	// 小红书内容
	platformContents["xiaohongshu"] = map[string][]string{
		"周杰伦": {
			"分享我的创作心得，如何保持灵感源源不断？#创作分享 #音乐制作",
			"推荐几本我最近在看的书，很有启发！#读书笔记",
			"我的私人歌单分享，这些歌陪伴了我很久～",
			"工作室的角落，我的创作空间～",
			"旅行随笔，记录美丽的风景！#旅行日记",
		},
		"杨幂": {
			"夏日穿搭指南，清爽又时尚！#夏日穿搭 #时尚搭配",
			"分享我的护肤routine，保持肌肤水润的秘诀！#护肤分享",
			"剧组的美食日常，每天都有好吃的！#美食分享",
			"推荐几个我常用的美妆产品，效果超赞！#美妆推荐",
			"生活小窍门，让你的生活更便捷！#生活技巧",
		},
	}

	// 如果有该明星在该平台的模拟内容，则返回
	if platformData, exists := platformContents[platform]; exists {
		if contents, exists := platformData[starName]; exists {
			return contents
		}
	}

	// 否则返回通用的模拟内容
	return []string{
		fmt.Sprintf("在%s上分享日常，感谢大家的支持！", platform),
		fmt.Sprintf("今天的%s更新，希望大家喜欢！", platform),
		fmt.Sprintf("和%s的朋友们一起度过了愉快的一天！", platform),
	}
}

// RealSocialCrawlerUpdateStyleFeatures 更新明星的语言风格特征实现
func RealSocialCrawlerUpdateStyleFeatures(star *models.Star, contents []string) error {
	if len(contents) == 0 {
		return nil
	}

	// 更高级的语言风格分析
	analysis := analyzeLanguageStyle(contents)

	// 结合原有的风格特征
	if star.StyleFeatures != "" {
		star.StyleFeatures = star.StyleFeatures + "\n\n" + analysis
	} else {
		star.StyleFeatures = analysis
	}

	return nil
}

// analyzeLanguageStyle 分析语言风格
func analyzeLanguageStyle(contents []string) string {
	// 初始化分析计数器
	totalLength := 0
	hashtagCount := 0
	exclamationCount := 0
	emojiCount := 0
	questionCount := 0
	averageWords := 0
	totalWords := 0

	// 正则表达式匹配表情符号
	emojiRegex := regexp.MustCompile(`[❤️💕😊🎉🎬✨]`)

	// 分析内容
	for _, content := range contents {
		totalLength += len(content)
		hashtagCount += strings.Count(content, "#")
		exclamationCount += strings.Count(content, "!")
		questionCount += strings.Count(content, "?")
		totalWords += len(strings.Fields(content))

		// 统计表情符号数量
		emojiMatches := emojiRegex.FindAllString(content, -1)
		emojiCount += len(emojiMatches)
	}

	// 计算平均值
	avgLength := totalLength / len(contents)
	if len(contents) > 0 {
		averageWords = totalWords / len(contents)
	}

	// 构建分析报告
	var analysis strings.Builder
	analysis.WriteString("语言风格分析报告：\n")

	// 内容长度分析
	if avgLength < 50 {
		analysis.WriteString("- 语言简洁，内容简短\n")
	} else if avgLength < 100 {
		analysis.WriteString("- 内容适中，表达清晰\n")
	} else {
		analysis.WriteString("- 内容详细，喜欢分享细节\n")
	}

	// 词汇丰富度分析
	if averageWords < 10 {
		analysis.WriteString("- 用词简练，言简意赅\n")
	} else if averageWords < 20 {
		analysis.WriteString("- 表达完整，信息量适中\n")
	} else {
		analysis.WriteString("- 喜欢详细描述，表达充分\n")
	}

	// 情感表达分析
	if emojiCount > len(contents)/2 {
		analysis.WriteString("- 频繁使用表情符号，情感表达丰富\n")
	}
	if exclamationCount > len(contents)/3 {
		analysis.WriteString("- 常用感叹号，语气活泼热情\n")
	}
	if questionCount > len(contents)/5 {
		analysis.WriteString("- 喜欢提问互动，注重粉丝交流\n")
	}

	// 话题标签分析
	if hashtagCount > len(contents) {
		analysis.WriteString("- 常用话题标签，注重内容分类\n")
	}

	return analysis.String()
}