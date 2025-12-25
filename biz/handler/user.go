package handler

import (
	"context"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"pet-service/biz/model"
	"pet-service/biz/service"
	"pet-service/pkg/jwt"
	"pet-service/pkg/logger"
	"pet-service/pkg/middleware"
)

// UserHandler 用户处理器
type UserHandler struct {
	userService service.UserService
	jwtManager  *jwt.JWTManager
}

// NewUserHandler 创建用户处理器
func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
		jwtManager:  middleware.GetJWTManager(),
	}
}

// CreateUser 创建用户
// @Summary 创建用户
// @Description 创建新用户
// @Tags 用户
// @Accept json
// @Produce json
// @Param request body model.CreateUserRequest true "创建用户请求"
// @Success 200 {object} utils.H
// @Router /api/v1/users [post]
func (h *UserHandler) CreateUser(ctx context.Context, c *app.RequestContext) {
	var req model.CreateUserRequest
	if err := c.BindAndValidate(&req); err != nil {
		logger.Error(ctx, "创建用户参数错误", logger.ErrorField(err))
		c.JSON(consts.StatusBadRequest, utils.H{
			"code":    400,
			"message": "参数错误",
			"error":   err.Error(),
		})
		return
	}

	user, err := h.userService.CreateUser(ctx, &req)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, utils.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(consts.StatusOK, utils.H{
		"code":    0,
		"message": "创建成功",
		"data": utils.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

// UpdateUser 更新用户
// @Summary 更新用户
// @Description 更新用户信息
// @Tags 用户
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Param request body model.UpdateUserRequest true "更新用户请求"
// @Success 200 {object} utils.H
// @Router /api/v1/users/{id} [put]
func (h *UserHandler) UpdateUser(ctx context.Context, c *app.RequestContext) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(consts.StatusBadRequest, utils.H{
			"code":    400,
			"message": "用户ID不能为空",
		})
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{
			"code":    400,
			"message": "用户ID格式错误",
		})
		return
	}

	var req model.UpdateUserRequest
	if err := c.BindAndValidate(&req); err != nil {
		logger.Error(ctx, "更新用户参数错误", logger.ErrorField(err))
		c.JSON(consts.StatusBadRequest, utils.H{
			"code":    400,
			"message": "参数错误",
			"error":   err.Error(),
		})
		return
	}

	userID := uint(id)
	user, err := h.userService.UpdateUser(ctx, userID, &req)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, utils.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(consts.StatusOK, utils.H{
		"code":    0,
		"message": "更新成功",
		"data":    user,
	})
}

// DeleteUser 删除用户
// @Summary 删除用户
// @Description 删除用户
// @Tags 用户
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Success 200 {object} utils.H
// @Router /api/v1/users/{id} [delete]
func (h *UserHandler) DeleteUser(ctx context.Context, c *app.RequestContext) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(consts.StatusBadRequest, utils.H{
			"code":    400,
			"message": "用户ID不能为空",
		})
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{
			"code":    400,
			"message": "用户ID格式错误",
		})
		return
	}

	userID := uint(id)
	if err := h.userService.DeleteUser(ctx, userID); err != nil {
		c.JSON(consts.StatusInternalServerError, utils.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(consts.StatusOK, utils.H{
		"code":    0,
		"message": "删除成功",
	})
}

// GetUser 获取用户详情
// @Summary 获取用户详情
// @Description 根据ID获取用户详情
// @Tags 用户
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Success 200 {object} utils.H
// @Router /api/v1/users/{id} [get]
func (h *UserHandler) GetUser(ctx context.Context, c *app.RequestContext) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(consts.StatusBadRequest, utils.H{
			"code":    400,
			"message": "用户ID不能为空",
		})
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{
			"code":    400,
			"message": "用户ID格式错误",
		})
		return
	}

	userID := uint(id)
	user, err := h.userService.GetUser(ctx, userID)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, utils.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(consts.StatusOK, utils.H{
		"code":    0,
		"message": "获取成功",
		"data":    user,
	})
}

// GetUserList 获取用户列表
// @Summary 获取用户列表
// @Description 获取用户列表
// @Tags 用户
// @Accept json
// @Produce json
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Param keyword query string false "关键词"
// @Param status query int false "状态"
// @Success 200 {object} utils.H
// @Router /api/v1/users [get]
func (h *UserHandler) GetUserList(ctx context.Context, c *app.RequestContext) {
	var req model.ListUserRequest
	if err := c.BindAndValidate(&req); err != nil {
		logger.Error(ctx, "获取用户列表参数错误", logger.ErrorField(err))
		c.JSON(consts.StatusBadRequest, utils.H{
			"code":    400,
			"message": "参数错误",
			"error":   err.Error(),
		})
		return
	}

	users, total, err := h.userService.GetUserList(ctx, &req)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, utils.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(consts.StatusOK, utils.H{
		"code":    0,
		"message": "获取成功",
		"data": utils.H{
			"list":      users,
			"total":     total,
			"page":      req.Page,
			"page_size": req.PageSize,
		},
	})
}

// Login 用户登录
// @Summary 用户登录
// @Description 用户登录
// @Tags 用户
// @Accept json
// @Produce json
// @Param request body model.LoginRequest true "登录请求"
// @Success 200 {object} utils.H
// @Router /api/v1/login [post]
func (h *UserHandler) Login(ctx context.Context, c *app.RequestContext) {
	var req model.LoginRequest
	if err := c.BindAndValidate(&req); err != nil {
		logger.Error(ctx, "登录参数错误", logger.ErrorField(err))
		c.JSON(consts.StatusBadRequest, utils.H{
			"code":    400,
			"message": "参数错误",
			"error":   err.Error(),
		})
		return
	}

	resp, err := h.userService.Login(ctx, &req, h.jwtManager)
	if err != nil {
		c.JSON(consts.StatusUnauthorized, utils.H{
			"code":    401,
			"message": err.Error(),
		})
		return
	}

	c.JSON(consts.StatusOK, utils.H{
		"code":    0,
		"message": "登录成功",
		"data":    resp,
	})
}

// GetCurrentUser 获取当前登录用户信息
// @Summary 获取当前用户信息
// @Description 获取当前登录用户的信息
// @Tags 用户
// @Accept json
// @Produce json
// @Success 200 {object} utils.H
// @Router /api/v1/me [get]
func (h *UserHandler) GetCurrentUser(ctx context.Context, c *app.RequestContext) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		c.JSON(consts.StatusUnauthorized, utils.H{
			"code":    401,
			"message": "未登录",
		})
		return
	}

	user, err := h.userService.GetUser(ctx, userID)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, utils.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(consts.StatusOK, utils.H{
		"code":    0,
		"message": "获取成功",
		"data":    user,
	})
}
