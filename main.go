package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"gorm.io/gorm"
	"pet-service/biz/handler"
	"pet-service/biz/repository"
	"pet-service/biz/service"
	"pet-service/config"
	"pet-service/pkg/database"
	"pet-service/pkg/logger"
	"pet-service/pkg/middleware"
	"pet-service/pkg/recovery"
	"pet-service/pkg/redis"
)

var (
	db          *gorm.DB
	cfg         *config.Config
	userHandler *handler.UserHandler
)

func main() {
	// 加载配置
	cfg = config.Load()

	// 初始化日志
	logLevel := logger.LevelInfo
	if cfg.Log.Level != "" {
		logLevel = logger.LogLevel(cfg.Log.Level)
	}
	if err := logger.Init(logLevel, cfg.Log.MaxSize, cfg.Log.MaxBackups, cfg.Log.MaxAge, cfg.Log.Compress, cfg.Log.OutputPath); err != nil {
		fmt.Printf("初始化日志失败: %v\n", err)
		os.Exit(1)
	}

	logger.Info(context.Background(), "===== 服务启动 =====")

	// 设置全局panic恢复
	defer recovery.HandlePanic()

	// 初始化Redis
	if err := redis.Init(cfg); err != nil {
		logger.Error(context.Background(), "Redis初始化失败", logger.ErrorField(err))
		// Redis初始化失败不阻止服务启动，只记录错误
		logger.Warn(context.Background(), "服务将在无Redis的情况下运行")
	}

	// 初始化数据库
	var err error
	db, err = database.Connect(cfg)
	if err != nil {
		logger.Error(context.Background(), "数据库连接失败", logger.ErrorField(err))
		logger.Warn(context.Background(), "服务将在无数据库的情况下运行")
	} else {
		logger.Info(context.Background(), "数据库连接成功")

		// 初始化仓储和服务
		userRepo := repository.NewUserRepository(db)
		userService := service.NewUserService(userRepo)
		userHandler = handler.NewUserHandler(userService)
	}

	h := server.Default(
		server.WithHostPorts(cfg.Server.Addr),
		server.WithReadTimeout(cfg.Server.ReadTimeout),
		server.WithWriteTimeout(cfg.Server.WriteTimeout),
		server.WithIdleTimeout(cfg.Server.IdleTimeout),
	)

	// 注册中间件
	h.Use(
		middleware.CORSMiddleware(),
		middleware.TraceIDMiddleware(),
		middleware.RequestLogMiddleware(),
		recovery.RecoveryMiddleware(),
	)

	// 初始化JWT管理器
	middleware.InitJWTManager(cfg.JWT.Secret, cfg.JWT.TokenDuration)
	logger.Info(context.Background(), "JWT管理器初始化成功")

	// 注册路由
	registerRoutes(h)

	// 优雅关闭
	go gracefulShutdown(h)

	// 启动服务
	logger.Info(context.Background(), "服务启动成功", logger.String("addr", cfg.Server.Addr))
	h.Spin()

	// 清理资源
	cleanup()
}

// registerRoutes 注册路由
func registerRoutes(h *server.Hertz) {
	// 健康检查
	h.GET("/health", func(ctx context.Context, c *app.RequestContext) {
		c.JSON(200, map[string]interface{}{
			"status": "ok",
		})
	})

	// Ping接口
	h.GET("/ping", func(ctx context.Context, c *app.RequestContext) {
		c.JSON(200, map[string]interface{}{
			"message": "pong",
		})
	})

	// API v1 路由组
	v1 := h.Group("/api/v1")
	{
		// 用户路由
		if userHandler != nil {
			// 公开路由 - 不需要认证
			v1.POST("/login", userHandler.Login)
			v1.POST("/users", userHandler.CreateUser)

			// 需要认证的路由
			authGroup := v1.Group("")
			authGroup.Use(middleware.JWTAuthMiddleware())
			{
				authGroup.GET("/me", userHandler.GetCurrentUser)

				userGroup := authGroup.Group("/users")
				{
					userGroup.PUT("/:id", userHandler.UpdateUser)
					userGroup.DELETE("/:id", userHandler.DeleteUser)
					userGroup.GET("/:id", userHandler.GetUser)
					userGroup.GET("", userHandler.GetUserList)
				}
			}
		}
	}
}

// gracefulShutdown 优雅关闭
func gracefulShutdown(h *server.Hertz) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info(context.Background(), "接收到关闭信号,开始优雅关闭...")

	// 关闭HTTP服务
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := h.Shutdown(ctx); err != nil {
		logger.Error(context.Background(), "服务关闭失败", logger.ErrorField(err))
	} else {
		logger.Info(context.Background(), "服务关闭成功")
	}
}

// cleanup 清理资源
func cleanup() {
	logger.Info(context.Background(), "开始清理资源...")

	// 关闭Redis连接
	if err := redis.Close(); err != nil {
		logger.Error(context.Background(), "关闭Redis连接失败", logger.ErrorField(err))
	}

	// 关闭日志
	if err := logger.Sync(); err != nil {
		logger.Error(context.Background(), "关闭日志失败", logger.ErrorField(err))
	}

	// 关闭数据库连接
	if db != nil {
		sqlDB, err := db.DB()
		if err == nil {
			sqlDB.Close()
		}
	}

	logger.Info(context.Background(), "资源清理完成")
}
