package service

import (
	"context"
	"errors"
	"fmt"

	"chat_agent/internal/ai"
	"chat_agent/internal/models"
	"chat_agent/internal/repository"

	"gorm.io/gorm"
)

// StarService 明星服务接口
type StarService interface {
	// 获取明星详情
	GetStarByID(ctx context.Context, starID uint) (*models.StarResponse, error)

	// 获取所有活跃的明星
	GetAllActiveStars(ctx context.Context) ([]models.StarResponse, error)

	// 分页获取明星列表
	GetStarList(ctx context.Context, page, pageSize int) ([]models.StarResponse, int64, error)

	// 搜索明星
	SearchStars(ctx context.Context, keyword string) ([]models.StarResponse, error)

	// 创建明星（管理员功能）
	CreateStar(ctx context.Context, req *models.CreateStarRequest) (*models.StarResponse, error)

	// 更新明星信息（管理员功能）
	UpdateStar(ctx context.Context, starID uint, req *models.UpdateStarRequest) (*models.StarResponse, error)

	// 删除明星（管理员功能）
	DeleteStar(ctx context.Context, starID uint) error

	// ToggleStarActive 切换明星活跃状态（管理员功能）
ToggleStarActive(ctx context.Context, starID uint, isActive bool) error

	// EnhanceStarProfile 增强明星资料
	EnhanceStarProfile(ctx context.Context, starID uint) error
}

// StarServiceImpl 明星服务实现
type StarServiceImpl struct {
	starRepo repository.StarRepository
}

// NewStarService 创建新的明星服务
func NewStarService(starRepo repository.StarRepository) StarService {
	return &StarServiceImpl{
		starRepo: starRepo,
	}
}

// GetStarByID 获取明星详情
func (s *StarServiceImpl) GetStarByID(ctx context.Context, starID uint) (*models.StarResponse, error) {
	star, err := s.starRepo.GetByID(ctx, starID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("明星不存在")
		}
		return nil, err
	}

	// 检查明星是否活跃
	if !star.IsActive {
		return nil, errors.New("明星当前不可用")
	}
	response := star.ToStarResponse()
	return &response, nil
}

// GetAllActiveStars 获取所有活跃的明星
func (s *StarServiceImpl) GetAllActiveStars(ctx context.Context) ([]models.StarResponse, error) {
	stars, err := s.starRepo.GetAllActive(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]models.StarResponse, len(stars))
	for i, star := range stars {
		responses[i] = star.ToStarResponse()
	}

	return responses, nil
}

// GetStarList 分页获取明星列表
func (s *StarServiceImpl) GetStarList(ctx context.Context, page, pageSize int) ([]models.StarResponse, int64, error) {
	// 参数校验
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	stars, total, err := s.starRepo.GetList(ctx, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]models.StarResponse, len(stars))
	for i, star := range stars {
		responses[i] = star.ToStarResponse()
	}

	return responses, total, nil
}

// SearchStars 搜索明星
func (s *StarServiceImpl) SearchStars(ctx context.Context, keyword string) ([]models.StarResponse, error) {
	// 参数校验
	if keyword == "" {
		return nil, errors.New("搜索关键词不能为空")
	}

	stars, err := s.starRepo.Search(ctx, keyword)
	if err != nil {
		return nil, err
	}

	responses := make([]models.StarResponse, len(stars))
	for i, star := range stars {
		responses[i] = star.ToStarResponse()
	}

	return responses, nil
}

// CreateStar 创建明星（管理员功能）
func (s *StarServiceImpl) CreateStar(ctx context.Context, req *models.CreateStarRequest) (*models.StarResponse, error) {
	// 创建明星对象
	star := &models.Star{
		Name:          req.Name,
		EnglishName:   req.EnglishName,
		Gender:        req.Gender,
		BirthDate:     req.BirthDate,
		Nationality:   req.Nationality,
		Occupation:    req.Occupation,
		Avatar:        req.Avatar,
		CoverImage:    req.CoverImage,
		Introduction:  req.Introduction,
		StyleFeatures: req.StyleFeatures,
		IsActive:      true, // 默认激活
	}

	// 保存到数据库
	if err := s.starRepo.Create(ctx, star); err != nil {
		return nil, err
	}

	response := star.ToStarResponse()
	return &response, nil
}

// UpdateStar 更新明星信息（管理员功能）
func (s *StarServiceImpl) UpdateStar(ctx context.Context, starID uint, req *models.UpdateStarRequest) (*models.StarResponse, error) {
	// 获取现有明星
	star, err := s.starRepo.GetByID(ctx, starID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("明星不存在")
		}
		return nil, err
	}

	// 更新字段
	if req.Name != "" {
		star.Name = req.Name
	}
	if req.EnglishName != "" {
		star.EnglishName = req.EnglishName
	}
	if req.Nationality != "" {
		star.Nationality = req.Nationality
	}
	if req.Occupation != "" {
		star.Occupation = req.Occupation
	}
	if req.BirthDate != "" {
		star.BirthDate = req.BirthDate
	}
	if req.Introduction != "" {
		star.Introduction = req.Introduction
	}
	if req.StyleFeatures != "" {
		star.StyleFeatures = req.StyleFeatures
	}
	if req.Gender != "" {
		star.Gender = req.Gender
	}
	if req.Avatar != "" {
		star.Avatar = req.Avatar
	}
	if req.CoverImage != "" {
		star.CoverImage = req.CoverImage
	}
	if req.IsActive != nil {
		star.IsActive = *req.IsActive
	}

	// 保存更新
	if err := s.starRepo.Update(ctx, star); err != nil {
		return nil, err
	}

	response := star.ToStarResponse()
	return &response, nil
}

// DeleteStar 删除明星（管理员功能）
func (s *StarServiceImpl) DeleteStar(ctx context.Context, starID uint) error {
	// 检查明星是否存在
	_, err := s.starRepo.GetByID(ctx, starID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("明星不存在")
		}
		return err
	}

	// 删除明星
	return s.starRepo.Delete(ctx, starID)
}

// ToggleStarActive 切换明星活跃状态（管理员功能）
func (s *StarServiceImpl) ToggleStarActive(ctx context.Context, starID uint, isActive bool) error {
	// 获取明星
	star, err := s.starRepo.GetByID(ctx, starID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("明星不存在")
		}
		return err
	}

	// 更新活跃状态
	star.IsActive = isActive

	// 保存更新
	return s.starRepo.Update(ctx, star)
}

// EnhanceStarProfile 增强明星资料
func (s *StarServiceImpl) EnhanceStarProfile(ctx context.Context, starID uint) error {
	// 调用AI模块的爬虫功能增强明星资料
	_, err := ai.EnhanceStarProfile(ctx, s.starRepo, starID)
	if err != nil {
		return fmt.Errorf("增强明星资料失败: %w", err)
	}
	return nil
}
