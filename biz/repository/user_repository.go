package repository

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"pet-service/biz/model"
	"pet-service/pkg/logger"
)

// UserRepository 用户仓储接口
type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	Update(ctx context.Context, id uint, user *model.User) error
	Delete(ctx context.Context, id uint) error
	GetByID(ctx context.Context, id uint) (*model.User, error)
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	List(ctx context.Context, offset, limit int, keyword string, status *int) ([]*model.User, int64, error)
}

// userRepository 用户仓储实现
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓储
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// Create 创建用户
func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	err := r.db.WithContext(ctx).Create(user).Error
	if err != nil {
		logger.Error(ctx, "创建用户失败", logger.String("username", user.Username), logger.ErrorField(err))
		return err
	}
	logger.Info(ctx, "创建用户成功", logger.Int("id", int(user.ID)))
	return nil
}

// Update 更新用户
func (r *userRepository) Update(ctx context.Context, id uint, user *model.User) error {
	result := r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).Updates(user)
	if result.Error != nil {
		logger.Error(ctx, "更新用户失败", logger.Int("id", int(id)), logger.ErrorField(result.Error))
		return result.Error
	}
	if result.RowsAffected == 0 {
		logger.Warn(ctx, "更新用户失败,用户不存在", logger.Int("id", int(id)))
		return errors.New("用户不存在")
	}
	logger.Info(ctx, "更新用户成功", logger.Int("id", int(id)))
	return nil
}

// Delete 删除用户(软删除)
func (r *userRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).Update("is_deleted", 1)
	if result.Error != nil {
		logger.Error(ctx, "删除用户失败", logger.Int("id", int(id)), logger.ErrorField(result.Error))
		return result.Error
	}
	if result.RowsAffected == 0 {
		logger.Warn(ctx, "删除用户失败,用户不存在", logger.Int("id", int(id)))
		return errors.New("用户不存在")
	}
	logger.Info(ctx, "删除用户成功", logger.Int("id", int(id)))
	return nil
}

// GetByID 根据ID获取用户
func (r *userRepository) GetByID(ctx context.Context, id uint) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("id = ? AND is_deleted = 0", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Warn(ctx, "获取用户失败,用户不存在", logger.Int("id", int(id)))
			return nil, errors.New("用户不存在")
		}
		logger.Error(ctx, "获取用户失败", logger.Int("id", int(id)), logger.ErrorField(err))
		return nil, err
	}
	return &user, nil
}

// GetByUsername 根据用户名获取用户
func (r *userRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("username = ? AND is_deleted = 0", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		logger.Error(ctx, "根据用户名获取用户失败", logger.String("username", username), logger.ErrorField(err))
		return nil, err
	}
	return &user, nil
}

// GetByEmail 根据邮箱获取用户
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("email = ? AND is_deleted = 0", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		logger.Error(ctx, "根据邮箱获取用户失败", logger.String("email", email), logger.ErrorField(err))
		return nil, err
	}
	return &user, nil
}

// List 获取用户列表
func (r *userRepository) List(ctx context.Context, offset, limit int, keyword string, status *int) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	query := r.db.WithContext(ctx).Model(&model.User{}).Where("is_deleted = 0")

	// 关键词搜索
	if keyword != "" {
		query = query.Where("username LIKE ? OR nickname LIKE ? OR email LIKE ? OR phone LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 状态过滤
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		logger.Error(ctx, "获取用户总数失败", logger.ErrorField(err))
		return nil, 0, err
	}

	// 获取列表
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&users).Error; err != nil {
		logger.Error(ctx, "获取用户列表失败", logger.ErrorField(err))
		return nil, 0, err
	}

	logger.Debug(ctx, "获取用户列表成功", logger.Int64("total", total), logger.Int("count", len(users)))
	return users, total, nil
}
