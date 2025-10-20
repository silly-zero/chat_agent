package repository

import (
	"context"

	"chat_agent/internal/models"

	"gorm.io/gorm"
)

// StarRepository 明星仓库接口
type StarRepository interface {
	// 创建明星
	Create(ctx context.Context, star *models.Star) error

	// 根据ID获取明星
	GetByID(ctx context.Context, id uint) (*models.Star, error)

	// 获取所有活跃的明星
	GetAllActive(ctx context.Context) ([]models.Star, error)

	// 分页获取明星列表
	GetList(ctx context.Context, page, pageSize int) ([]models.Star, int64, error)

	// 更新明星信息
	Update(ctx context.Context, star *models.Star) error

	// 删除明星
	Delete(ctx context.Context, id uint) error

	// 搜索明星
	Search(ctx context.Context, keyword string) ([]models.Star, error)
}

// StarRepositoryImpl 明星仓库实现
type StarRepositoryImpl struct {
	db *gorm.DB
}

// NewStarRepository 创建新的明星仓库
func NewStarRepository(db *gorm.DB) StarRepository {
	return &StarRepositoryImpl{db: db}
}

// Create 创建明星
func (r *StarRepositoryImpl) Create(ctx context.Context, star *models.Star) error {
	return r.db.WithContext(ctx).Create(star).Error
}

// GetByID 根据ID获取明星
func (r *StarRepositoryImpl) GetByID(ctx context.Context, id uint) (*models.Star, error) {
	var star models.Star
	err := r.db.WithContext(ctx).First(&star, id).Error
	if err != nil {
		return nil, err
	}
	return &star, nil
}

// GetAllActive 获取所有活跃的明星
func (r *StarRepositoryImpl) GetAllActive(ctx context.Context) ([]models.Star, error) {
	var stars []models.Star
	err := r.db.WithContext(ctx).Where("is_active = ?", true).Find(&stars).Error
	if err != nil {
		return nil, err
	}
	return stars, nil
}

// GetList 分页获取明星列表
func (r *StarRepositoryImpl) GetList(ctx context.Context, page, pageSize int) ([]models.Star, int64, error) {
	var stars []models.Star
	var total int64

	// 计算偏移量
	offset := (page - 1) * pageSize

	// 获取总数
	if err := r.db.WithContext(ctx).Model(&models.Star{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err := r.db.WithContext(ctx).Offset(offset).Limit(pageSize).Find(&stars).Error
	if err != nil {
		return nil, 0, err
	}

	return stars, total, nil
}

// Update 更新明星信息
func (r *StarRepositoryImpl) Update(ctx context.Context, star *models.Star) error {
	return r.db.WithContext(ctx).Save(star).Error
}

// Delete 删除明星
func (r *StarRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Star{}, id).Error
}

// Search 搜索明星
func (r *StarRepositoryImpl) Search(ctx context.Context, keyword string) ([]models.Star, error) {
	var stars []models.Star
	query := r.db.WithContext(ctx)

	// 根据关键词搜索姓名、英文名、国籍等字段
	searchPattern := "%" + keyword + "%"
	query = query.Where(
		"name LIKE ? OR english_name LIKE ? OR nationality LIKE ? OR occupation LIKE ? OR introduction LIKE ?",
		searchPattern, searchPattern, searchPattern, searchPattern, searchPattern,
	)

	// 只返回活跃的明星
	query = query.Where("is_active = ?", true)

	err := query.Find(&stars).Error
	if err != nil {
		return nil, err
	}

	return stars, nil
}