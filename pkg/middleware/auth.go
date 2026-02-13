package middleware

import (
	"context"
	"strings"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"pet-service/pkg/jwt"
	"pet-service/pkg/logger"
)

var jwtManager *jwt.JWTManager

// InitJWTManager 初始化JWT管理器
func InitJWTManager(secretKey string, tokenDuration int) {
	jwtManager = jwt.NewJWTManager(secretKey, time.Duration(int64(tokenDuration)*int64(1e9)))
}

// GetJWTManager 获取JWT管理器
func GetJWTManager() *jwt.JWTManager {
	return jwtManager
}

// JWTAuthMiddleware JWT认证中间件
func JWTAuthMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// 从Header获取token
		authHeader := string(c.GetHeader("Authorization"))
		if authHeader == "" {
			logger.Warn(ctx, "未携带认证token")
			c.JSON(consts.StatusUnauthorized, map[string]interface{}{
				"code":    401,
				"message": "未携带认证token",
			})
			c.Abort()
			return
		}

		// 解析Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			logger.Warn(ctx, "认证token格式错误")
			c.JSON(consts.StatusUnauthorized, map[string]interface{}{
				"code":    401,
				"message": "认证token格式错误",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 验证token
		claims, err := jwtManager.ValidateToken(tokenString)
		if err != nil {
			logger.Warn(ctx, "认证token无效或已过期", logger.ErrorField(err))
			c.JSON(consts.StatusUnauthorized, map[string]interface{}{
				"code":    401,
				"message": "认证token无效或已过期",
			})
			c.Abort()
			return
		}

		// 将用户信息存入上下文
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)

		c.Next(ctx)
	}
}

// GetUserID 从上下文获取用户ID
func GetUserID(c *app.RequestContext) uint {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0
	}
	return userID.(uint)
}

// GetUsername 从上下文获取用户名
func GetUsername(c *app.RequestContext) string {
	username, exists := c.Get("username")
	if !exists {
		return ""
	}
	return username.(string)
}
