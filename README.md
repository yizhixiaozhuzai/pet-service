# Pet Service

基于 Hertz 框架的宠物服务项目

## 项目特性

- ✅ 完善的日志系统，支持链路追踪
- ✅ 优雅退出机制
- ✅ Panic自动恢复和重启
- ✅ Redis缓存支持
- ✅ GORM代码生成
- ✅ 完整的用户管理功能

## 项目结构

```
pet-service/
├── biz/                    # 业务逻辑层
│   ├── handler/            # HTTP处理器
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
├── internal/               # 内部模块
│   └── cron/               # 定时任务
├── script/                 # 脚本文件
├── main.go                 # 主入口
├── router.go               # 路由注册
├── router_gen.go           # 生成的路由
└── go.mod                  # Go模块文件
```

## 快速开始

### 1. 环境准备

- Go 1.25.5+
- MySQL 5.7+
- Redis 6.0+ (可选)

### 2. 克隆项目

```bash
git clone <repository-url>
cd pet-service
```

### 3. 配置环境变量

复制 `.env.example` 为 `.env` 并修改配置：

```bash
cp .env.example .env
```

修改 `.env` 文件中的数据库和Redis配置：

```bash
# 服务器配置
SERVER_ADDR=:8888

# Redis配置（可选，不配置服务也能启动）
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0

# 数据库配置（可选，不配置服务也能启动）
DATABASE_DSN=root:password@tcp(localhost:3306)/pet_service?charset=utf8mb4&parseTime=True&loc=Local
```

### 4. 初始化数据库（可选）

如果需要使用数据库功能，执行初始化脚本：

```bash
mysql -u root -p < script/init.sql
```

### 5. 安装依赖

```bash
go mod tidy
```

### 6. 运行服务

```bash
# Linux/Mac
./script/run.sh

# Windows
script\run.bat

# 或直接运行
go run main.go
```

服务将在 http://localhost:8888 启动

### 7. 测试服务

```bash
# Linux/Mac
./test.sh

# Windows
test.bat
```

## API文档

### 健康检查

```bash
GET /health
```

响应：
```json
{
  "status": "ok"
}
```

### Ping测试

```bash
GET /ping
```

响应：
```json
{
  "message": "pong"
}
```

### 用户管理

#### 创建用户
```bash
POST /api/v1/users
Content-Type: application/json

{
  "username": "testuser",
  "password": "123456",
  "email": "test@example.com",
  "phone": "13800138000",
  "nickname": "测试用户"
}
```

#### 更新用户
```bash
PUT /api/v1/users/{id}
Content-Type: application/json

{
  "email": "newemail@example.com",
  "nickname": "新昵称"
}
```

#### 删除用户
```bash
DELETE /api/v1/users/{id}
```

#### 获取用户详情
```bash
GET /api/v1/users/{id}
```

#### 获取用户列表
```bash
GET /api/v1/users?page=1&page_size=10&keyword=test&status=1
```

#### 用户登录
```bash
POST /api/v1/login
Content-Type: application/json

{
  "username": "testuser",
  "password": "123456"
}
```

## 日志系统

项目使用zap日志库，支持以下功能：
- 结构化日志输出
- 日志轮转
- 链路追踪(trace_id)
- 调用者信息记录
- 多级别日志

日志文件位置：`./logs/app.log`

日志使用示例：
```go
logger.Info(ctx, "日志信息", logger.String("key", "value"))
logger.Error(ctx, "错误信息", logger.ErrorField(err))
logger.WithCaller(ctx, "重要操作", logger.Int("user_id", 123))
```

## 恢复机制

- Panic自动恢复
- Panic次数统计
- 超过阈值自动退出
- 优雅关闭支持

## 中间件

- TraceIDMiddleware: 链路追踪
- CORSMiddleware: 跨域支持
- RequestLogMiddleware: 请求日志
- RecoveryMiddleware: 错误恢复

## 开发指南

### 添加新功能

1. 在 `biz/model` 中定义数据模型
2. 在 `biz/repository` 中创建仓储接口和实现
3. 在 `biz/service` 中创建业务逻辑
4. 在 `biz/handler` 中创建HTTP处理器
5. 在 `main.go` 的 `registerRoutes` 函数中注册路由

### 数据库代码生成

修改 `script/gen.go`，添加需要生成的表：

```go
g.ApplyBasic(
    g.GenerateModel("users"),
    g.GenerateModel("products"), // 添加新表
)
```

然后运行生成脚本：
```bash
go run script/gen.go
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

## 故障排查

### 编译错误

```bash
# 清理缓存并重新下载依赖
go clean -modcache
go mod tidy
go build
```

### 数据库连接失败

1. 检查数据库是否启动
2. 检查 `DATABASE_DSN` 配置是否正确
3. 确认数据库用户权限

### Redis连接失败

Redis是可选的，连接失败不会阻止服务启动。如果需要使用Redis：
1. 检查Redis是否启动
2. 检查 `REDIS_ADDR` 配置是否正确
3. 确认Redis密码配置

### 端口被占用

修改 `.env` 文件中的 `SERVER_ADDR` 配置：

```bash
SERVER_ADDR=:9999
```

## 许可证

MIT License
