package main

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
)

func main() {
	fmt.Println("开始测试Hertz服务...")

	h := server.Default(
		server.WithHostPorts(":8888"),
	)

	h.GET("/health", func(ctx context.Context, c *app.RequestContext) {
		c.JSON(200, map[string]interface{}{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	h.GET("/ping", func(ctx context.Context, c *app.RequestContext) {
		c.JSON(200, map[string]interface{}{
			"message": "pong",
		})
	})

	fmt.Println("服务启动成功，监听 :8888")
	fmt.Println("访问 http://localhost:8888/health 测试")
	fmt.Println("按 Ctrl+C 停止服务")

	h.Spin()
}
