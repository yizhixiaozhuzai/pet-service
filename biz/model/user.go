package model

import (
	"time"
)

// User 用户模型
type User struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Username  string    `json:"username" gorm:"type:varchar(50);uniqueIndex;not null;comment:用户名"`
	Password  string    `json:"-" gorm:"type:varchar(255);not null;comment:密码"`
	Email     string    `json:"email" gorm:"type:varchar(100);uniqueIndex;comment:邮箱"`
	Phone     string    `json:"phone" gorm:"type:varchar(20);uniqueIndex;comment:手机号"`
	Nickname  string    `json:"nickname" gorm:"type:varchar(50);comment:昵称"`
	Avatar    string    `json:"avatar" gorm:"type:varchar(255);comment:头像"`
	Status    int       `json:"status" gorm:"type:tinyint;default:1;comment:状态:0禁用,1正常"`
	IsDeleted int       `json:"is_deleted" gorm:"type:tinyint;default:0;comment:是否删除:0否,1是"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6"`
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone" binding:"omitempty,len=11"`
	Nickname string `json:"nickname" binding:"omitempty,max=50"`
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	Email    string `json:"email" binding:"omitempty,email"`
	Phone    string `json:"phone" binding:"omitempty,len=11"`
	Nickname string `json:"nickname" binding:"omitempty,max=50"`
	Avatar   string `json:"avatar" binding:"omitempty,max=255"`
	Status   *int   `json:"status" binding:"omitempty,oneof=0 1"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UserResponse 用户响应
type UserResponse struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Nickname  string    `json:"nickname"`
	Avatar    string    `json:"avatar"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token     string       `json:"token"`
	TokenType string       `json:"token_type"`
	ExpiresIn int64        `json:"expires_in"`
	User      UserResponse `json:"user"`
}

// ListUserRequest 用户列表请求
type ListUserRequest struct {
	Page     int    `form:"page,default=1" binding:"min=1"`
	PageSize int    `form:"page_size,default=10" binding:"min=1,max=100"`
	Keyword  string `form:"keyword"`
	Status   *int   `form:"status" binding:"omitempty,oneof=0 1"`
}
