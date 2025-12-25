# 项目快速开始指南

## 项目初始化已完成

本Hertz框架项目已按照您的要求完成初始化，包含以下功能：

### 1. 完善的日志系统
- ✅ 使用zap日志库
- ✅ 支持结构化日志
- ✅ 自动日志轮转(lumberjack)
- ✅ 链路追踪能力
- ✅ 记录代码调用位置
- ✅ 日志级别控制(Debug/Info/Warn/Error/Fatal)

### 2. 完善的优雅退出机制
- ✅ Panic自动恢复
- ✅ Panic次数统计和限制
- ✅ 超过阈值自动退出
- ✅ SIGINT/SIGTERM信号处理
- ✅ 资源清理机制

### 3. 完整的项目结构
```
pet-service/
├── biz/                    # 业务逻辑层
│   ├── handler/            # HTTP处理器 (含User示例)
│   ├── service/            # 业务服务
│   ├── repository/         # 数据仓储
│   ├── model/              # 数据模型
│   └── router/             # 路由定义
├── config/                 # 配置管理
├── pkg/                    # 公共包
│   ├── logger/             # 日志系统
│   ├── middleware/         # 中间件
│   ├── redis/              # Redis客户端
│   ├── database/           # 数据库连接
│   └── recovery/           # 恢复机制
├── script/                 # 脚本文件
├── main.go                 # 主入口
└── README.md               # 项目文档
```

### 4. Redis缓存能力
- ✅ go-redis/v9集成
- ✅ 连接池配置
- ✅ 常用操作封装(Get/Set/Del/Hash等)
- ✅ 错误处理和日志记录

### 5. GORM Gen支持
- ✅ gorm.io/gen集成
- ✅ 代码生成脚本(script/gen.go)
- ✅ 数据库连接管理
- ✅ 配置文件参考: https://gorm.io/zh_CN/gen/index.html

## 开始使用

### 步骤1: 安装依赖
```bash
cd c:/project/study/pet-service
go mod tidy
```

### 步骤2: 配置数据库
修改 `config/config.go` 或设置环境变量：
```bash
DATABASE_DSN=root:password@tcp(localhost:3306)/pet_service?charset=utf8mb4&parseTime=True&loc=Local
REDIS_ADDR=localhost:6379
```

### 步骤3: 初始化数据库
```bash
mysql -u root -p < script/init.sql
```

### 步骤4: 生成GORM代码
```bash
go run script/gen.go
```

### 步骤5: 运行服务
```bash
go run main.go
```

服务将在 http://localhost:8888 启动

## API测试

### 健康检查
```bash
curl http://localhost:8888/health
```

### 创建用户
```bash
curl -X POST http://localhost:8888/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "123456",
    "email": "test@example.com",
    "nickname": "测试用户"
  }'
```

### 获取用户列表
```bash
curl http://localhost:8888/api/v1/users?page=1&page_size=10
```

### 用户登录
```bash
curl -X POST http://localhost:8888/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "123456"
  }'
```

## 日志示例

日志会自动记录到 `./logs/app.log`，同时输出到控制台：

```json
{
  "time": "2025-12-25T10:30:45+08:00",
  "level": "info",
  "caller": "main.go:65",
  "function": "main.registerRoutes",
  "trace_id": "550e8400-e29b-41d4-a716-446655440000",
  "msg": "收到请求",
  "method": "POST",
  "path": "/api/v1/users",
  "client_ip": "127.0.0.1"
}
```

## 核心功能说明

### 日志系统使用
```go
// 基础使用
logger.Info(ctx, "日志信息", logger.String("key", "value"))
logger.Error(ctx, "错误信息", logger.ErrorField(err))

// 带调用者信息
logger.WithCaller(ctx, "重要操作", logger.Int("user_id", 123))
```

### Redis使用
```go
// 设置缓存
redis.Set(ctx, "user:123", userData, 5*time.Minute)

// 获取缓存
val, _ := redis.Get(ctx, "user:123")

// 删除缓存
redis.Del(ctx, "user:123")
```

### 中间件使用
```go
// 在main.go中已自动注册
h.Use(
    middleware.CORSMiddleware(),      // 跨域
    middleware.TraceIDMiddleware(),   // 链路追踪
    middleware.RequestLogMiddleware(), // 请求日志
    recovery.RecoveryMiddleware(),    // 错误恢复
)
```

## 下一步开发

参考User模块的实现，添加新的业务模块：

1. 在 `biz/model` 定义模型
2. 在 `biz/repository` 创建仓储
3. 在 `biz/service` 实现业务逻辑
4. 在 `biz/handler` 创建处理器
5. 在 `main.go` 注册路由

## 注意事项

1. 确保MySQL和Redis服务已启动
2. 修改数据库连接字符串为实际配置
3. 首次运行需要执行初始化SQL脚本
4. 日志文件会自动轮转，无需手动清理
5. Panic会自动恢复，但超过10次会自动退出

如有问题，请查看日志文件 `./logs/app.log`
