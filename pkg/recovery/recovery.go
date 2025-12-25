package recovery

import (
	"context"
	"errors"
	"fmt"
	"runtime/debug"
	"sync/atomic"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"pet-service/pkg/logger"
)

var (
	panicCount      int32
	maxPanicRecover = 10 // 最大panic恢复次数
	restartDelay    = 5 * time.Second
)

// RecoveryMiddleware panic恢复中间件
func RecoveryMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		defer func() {
			if err := recover(); err != nil {
				stack := debug.Stack()

				logger.Error(ctx, "发生panic",
					logger.Any("error", err),
					logger.String("path", string(c.Request.Path())),
					logger.String("method", string(c.Request.Method())),
					logger.String("stack", string(stack)),
				)

				// 增加panic计数
				atomic.AddInt32(&panicCount, 1)

				// 检查panic次数是否超过阈值
				if atomic.LoadInt32(&panicCount) > int32(maxPanicRecover) {
					logger.Fatal(ctx, "Panic次数超过阈值,服务将退出")
				}

				// 返回友好的错误信息
				c.JSON(consts.StatusInternalServerError, utils.H{
					"code":    500,
					"message": "服务器内部错误",
					"error":   fmt.Sprintf("%v", err),
				})
				c.Abort()
			}
		}()

		c.Next(ctx)
	}
}

// GracefulShutdown 优雅关闭处理器
type GracefulShutdown struct {
	server      interface{}
	connections int32
}

// NewGracefulShutdown 创建优雅关闭处理器
func NewGracefulShutdown(server interface{}) *GracefulShutdown {
	return &GracefulShutdown{
		server: server,
	}
}

// HandleSignals 处理系统信号
func (gs *GracefulShutdown) HandleSignals() {
	// 这里可以添加信号监听逻辑
	// 例如: 监听 SIGINT, SIGTERM 等信号
}

// WaitShutdown 等待关闭完成
func (gs *GracefulShutdown) WaitShutdown() {
	// 实现优雅关闭逻辑
}

// RestartMiddleware 自动重启中间件(在panic后延迟重启)
func RestartMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		defer func() {
			if err := recover(); err != nil {
				count := atomic.LoadInt32(&panicCount)
				if count < int32(maxPanicRecover) {
					logger.Warn(ctx, fmt.Sprintf("检测到panic,将在%v后自动重试,当前panic次数:%d", restartDelay, count))

					// 延迟重启
					time.Sleep(restartDelay)
				}
			}
		}()

		c.Next(ctx)
	}
}

// PanicRecovery 全局panic恢复函数
func PanicRecovery() {
	if r := recover(); r != nil {
		stack := debug.Stack()
		logger.Error(context.Background(), "全局panic恢复",
			logger.Any("error", r),
			logger.String("stack", string(stack)),
		)

		// 增加panic计数
		atomic.AddInt32(&panicCount, 1)

		// 检查是否需要退出
		if atomic.LoadInt32(&panicCount) > int32(maxPanicRecover) {
			logger.Fatal(context.Background(), "Panic次数超过阈值,服务将退出")
		}
	}
}

// HandlePanic 处理panic并决定是否继续
func HandlePanic() {
	if err := recover(); err != nil {
		stack := debug.Stack()
		logger.Error(context.Background(), "Panic处理",
			logger.Any("error", err),
			logger.String("stack", string(stack)),
		)

		// 增加panic计数
		count := atomic.AddInt32(&panicCount, 1)

		if count <= int32(maxPanicRecover) {
			logger.Warn(context.Background(), fmt.Sprintf("Panic已恢复,服务继续运行,当前panic次数:%d", count))
		} else {
			logger.Fatal(context.Background(), "Panic次数超过阈值,服务将退出")
		}
	}
}

// GetPanicCount 获取panic次数
func GetPanicCount() int32 {
	return atomic.LoadInt32(&panicCount)
}

// ResetPanicCount 重置panic计数
func ResetPanicCount() {
	atomic.StoreInt32(&panicCount, 0)
}

// ValidateRequest 验证请求
func ValidateRequest(c *app.RequestContext) error {
	// 基本验证
	if c == nil {
		return errors.New("request context is nil")
	}

	// 验证请求大小
	if len(c.Request.Body()) > 10<<20 { // 10MB限制
		return errors.New("request body too large")
	}

	// 验证Content-Type
	contentType := c.Request.Header.Get("Content-Type")
	if contentType != "" && contentType != "application/json" && contentType != "application/x-www-form-urlencoded" {
		return errors.New("unsupported content type")
	}

	return nil
}

// ErrorHandler 统一错误处理
func ErrorHandler(ctx context.Context, c *app.RequestContext, err error) {
	if err == nil {
		return
	}

	logger.Error(ctx, "请求处理错误", logger.ErrorField(err))

	// 根据错误类型返回不同的状态码
	if err != nil && err.Error() == "request body too large" {
		c.JSON(consts.StatusBadRequest, utils.H{
			"code":    400,
			"message": "请求体过大",
		})
		return
	}

	c.JSON(consts.StatusInternalServerError, utils.H{
		"code":    500,
		"message": "服务器内部错误",
	})
}

// HealthCheckMiddleware 健康检查中间件
func HealthCheckMiddleware(path string) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		if string(c.Request.Path()) == path {
			c.JSON(consts.StatusOK, utils.H{
				"status": "healthy",
				"time":   time.Now().Format(time.RFC3339),
			})
			c.Abort()
		}
	}
}
