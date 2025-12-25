package middleware

import (
	"context"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/google/uuid"
	"pet-service/pkg/logger"
)

// TraceIDMiddleware 链路追踪中间件
func TraceIDMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// 生成或获取trace_id
		traceID := string(c.GetHeader("X-Trace-ID"))
		if traceID == "" {
			traceID = uuid.New().String()
		}

		// 设置到context中
		c.Set("trace_id", traceID)
		c.Header("X-Trace-ID", traceID)

		// 记录请求开始时间
		startTime := time.Now()

		// 处理请求
		c.Next(ctx)

		// 记录请求日志
		duration := time.Since(startTime)
		logger.Info(ctx, "请求处理完成",
			logger.String("method", string(c.Request.Method())),
			logger.String("path", string(c.Request.Path())),
			logger.Int("status", c.Response.StatusCode()),
			logger.String("duration", duration.String()),
		)
	}
}

// CORSMiddleware 跨域中间件
func CORSMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Trace-ID")
		c.Header("Access-Control-Expose-Headers", "Content-Length, X-Trace-ID")
		c.Header("Access-Control-Allow-Credentials", "true")

		// 处理OPTIONS预检请求
		if string(c.Request.Method()) == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next(ctx)
	}
}

// RequestLogMiddleware 请求日志中间件
func RequestLogMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		startTime := time.Now()

		logger.Info(ctx, "收到请求",
			logger.String("method", string(c.Request.Method())),
			logger.String("path", string(c.Request.Path())),
			logger.String("query", string(c.Request.QueryString())),
			logger.String("client_ip", c.ClientIP()),
		)

		c.Next(ctx)

		duration := time.Since(startTime)
		logger.Info(ctx, "请求完成",
			logger.Int("status", c.Response.StatusCode()),
			logger.String("duration", duration.String()),
		)
	}
}

// RecoveryMiddleware 错误恢复中间件
func RecoveryMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error(ctx, "请求处理panic",
					logger.Any("error", err),
					logger.String("path", string(c.Request.Path())),
				)
				c.JSON(500, map[string]interface{}{
					"code":    500,
					"message": "服务器内部错误",
				})
				c.Abort()
			}
		}()

		c.Next(ctx)
	}
}
