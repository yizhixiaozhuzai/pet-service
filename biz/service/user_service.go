package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
	"pet-service/biz/model"
	"pet-service/biz/repository"
	"pet-service/pkg/jwt"
	"pet-service/pkg/logger"
	"pet-service/pkg/redis"
)

// UserService 用户服务接口
type UserService interface {
	CreateUser(ctx context.Context, req *model.CreateUserRequest) (*model.User, error)
	UpdateUser(ctx context.Context, id uint, req *model.UpdateUserRequest) (*model.User, error)
	DeleteUser(ctx context.Context, id uint) error
	GetUser(ctx context.Context, id uint) (*model.User, error)
	GetUserList(ctx context.Context, req *model.ListUserRequest) ([]*model.User, int64, error)
	Login(ctx context.Context, req *model.LoginRequest, jwtManager *jwt.JWTManager) (*model.LoginResponse, error)
}

// userService 用户服务实现
type userService struct {
	userRepo repository.UserRepository
}

// NewUserService 创建用户服务
func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

// CreateUser 创建用户
func (s *userService) CreateUser(ctx context.Context, req *model.CreateUserRequest) (*model.User, error) {
	// 检查用户名是否已存在
	if _, err := s.userRepo.GetByUsername(ctx, req.Username); err == nil {
		logger.Warn(ctx, "用户名已存在", logger.String("username", req.Username))
		return nil, errors.New("用户名已存在")
	}

	// 检查邮箱是否已存在
	if _, err := s.userRepo.GetByEmail(ctx, req.Email); err == nil {
		logger.Warn(ctx, "邮箱已存在", logger.String("email", req.Email))
		return nil, errors.New("邮箱已存在")
	}

	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error(ctx, "密码加密失败", logger.ErrorField(err))
		return nil, errors.New("密码加密失败")
	}

	user := &model.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Email:    req.Email,
		Phone:    req.Phone,
		Nickname: req.Nickname,
		Status:   1,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// 清除缓存
	_ = redis.Del(ctx, "users:all")

	logger.Info(ctx, "用户创建成功", logger.String("username", req.Username))
	return user, nil
}

// UpdateUser 更新用户
func (s *userService) UpdateUser(ctx context.Context, id uint, req *model.UpdateUserRequest) (*model.User, error) {
	// 获取用户
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 更新字段
	if req.Email != "" {
		// 检查邮箱是否已被其他用户使用
		if existUser, err := s.userRepo.GetByEmail(ctx, req.Email); err == nil && existUser.ID != id {
			return nil, errors.New("邮箱已被其他用户使用")
		}
		user.Email = req.Email
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}
	if req.Status != nil {
		user.Status = *req.Status
	}

	if err := s.userRepo.Update(ctx, id, user); err != nil {
		return nil, err
	}

	// 清除缓存
	_ = redis.Del(ctx, "users:all", fmt.Sprintf("user:%d", id))

	logger.Info(ctx, "用户更新成功", logger.Int("id", int(id)))
	return user, nil
}

// DeleteUser 删除用户
func (s *userService) DeleteUser(ctx context.Context, id uint) error {
	if err := s.userRepo.Delete(ctx, id); err != nil {
		return err
	}

	// 清除缓存
	_ = redis.Del(ctx, "users:all", fmt.Sprintf("user:%d", id))

	logger.Info(ctx, "用户删除成功", logger.Int("id", int(id)))
	return nil
}

// GetUser 获取用户详情
func (s *userService) GetUser(ctx context.Context, id uint) (*model.User, error) {
	// 先从缓存获取
	cacheKey := fmt.Sprintf("user:%d", id)
	if cached, err := redis.Get(ctx, cacheKey); err == nil && cached != "" {
		logger.Debug(ctx, "从缓存获取用户", logger.Int("id", int(id)))
		// TODO: 反序列化缓存数据
		// user := &model.User{}
		// json.Unmarshal([]byte(cached), user)
		// return user, nil
	}

	// 从数据库获取
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 设置缓存
	// TODO: 序列化并缓存
	// if data, err := json.Marshal(user); err == nil {
	// 	redis.Set(ctx, cacheKey, data, 5*time.Minute)
	// }

	return user, nil
}

// GetUserList 获取用户列表
func (s *userService) GetUserList(ctx context.Context, req *model.ListUserRequest) ([]*model.User, int64, error) {
	offset := (req.Page - 1) * req.PageSize

	users, total, err := s.userRepo.List(ctx, offset, req.PageSize, req.Keyword, req.Status)
	if err != nil {
		return nil, 0, err
	}

	logger.Info(ctx, "获取用户列表成功",
		logger.Int("page", req.Page),
		logger.Int("page_size", req.PageSize),
		logger.Int64("total", total),
	)

	return users, total, nil
}

// Login 用户登录
func (s *userService) Login(ctx context.Context, req *model.LoginRequest, jwtManager *jwt.JWTManager) (*model.LoginResponse, error) {
	// 获取用户
	user, err := s.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		logger.Warn(ctx, "登录失败,用户不存在", logger.String("username", req.Username))
		return nil, errors.New("用户名或密码错误")
	}

	// 检查用户状态
	if user.Status != 1 {
		logger.Warn(ctx, "登录失败,用户已被禁用", logger.String("username", req.Username))
		return nil, errors.New("用户已被禁用")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		logger.Warn(ctx, "登录失败,密码错误", logger.String("username", req.Username))
		return nil, errors.New("用户名或密码错误")
	}

	// 生成token
	token, expiresIn, err := jwtManager.GenerateToken(user.ID, user.Username)
	if err != nil {
		logger.Error(ctx, "生成token失败", logger.ErrorField(err))
		return nil, errors.New("生成token失败")
	}

	// 设置缓存
	cacheKey := fmt.Sprintf("token:%s", token)
	_ = redis.Set(ctx, cacheKey, user.ID, 24*time.Hour)

	logger.Info(ctx, "用户登录成功", logger.String("username", req.Username))

	response := &model.LoginResponse{
		Token:     token,
		TokenType: "Bearer",
		ExpiresIn: expiresIn,
		User: model.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			Phone:     user.Phone,
			Nickname:  user.Nickname,
			Avatar:    user.Avatar,
			Status:    user.Status,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	}

	return response, nil
}
