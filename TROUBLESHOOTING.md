# 常见问题解决

## 问题1: `invalid array length -delta * delta (constant -256 of type int64)`

### 错误信息
```
C:\Users\Administrator\go\pkg\mod\golang.org\x\tools@v0.21.1-0.20240508182429-e35e4ccd0d2d\internal\tokeninternal\tokeninternal.go:64:9: invalid array length -delta * delta (constant -256 of type int64)
```

### 原因
`golang.org/x/tools` 包的特定版本 `v0.21.1-0.20240508182429-e35e4ccd0d2d` 存在编译器bug。

### 解决方案 ✅

```bash
# 1. 降级 golang.org/x/tools 到稳定版本
go get -u golang.org/x/tools@v0.20.0

# 2. 清理并重新整理依赖
go mod tidy

# 3. 重新编译
go build -o pet-service.exe main.go
```

**验证**:
```bash
# 编译应该成功，无错误输出
go build -o pet-service.exe main.go
echo %errorlevel%
# 应该输出: 0
```

### 其他可选版本
如果 `v0.20.0` 还有问题，可以尝试：
```bash
go get -u golang.org/x/tools@v0.19.0
# 或
go get -u golang.org/x/tools@v0.18.0
```

---

## 问题2: Go版本不存在

### 错误信息
```
go.mod: unknown go version '1.25.5'
```

### 原因
`go.mod` 文件中设置的Go版本不存在（Go目前最新版本是1.23.x）。

### 解决方案

将 `go.mod` 文件中的Go版本改为实际存在的版本：

```go
// 错误的版本
go 1.25.5

// 正确的版本
go 1.21  // 或 1.20, 1.22, 1.23 等实际存在的版本
```

### 可用的Go版本

以下是一些常用的稳定版本：
- 1.20.x
- 1.21.x (LTS)
- 1.22.x
- 1.23.x (最新)

---

## 其他常见问题

### 3. 端口被占用

**错误信息**:
```
bind: address already in use
```

**解决方案**:
- 修改 `.env` 文件中的端口配置
```bash
SERVER_ADDR=:9999
```

### 4. 数据库连接失败

**错误信息**:
```
Error 2003: Can't connect to MySQL server
```

**解决方案**:
1. 检查MySQL服务是否启动
2. 检查 `DATABASE_DSN` 配置
3. 检查数据库用户权限

### 5. Redis连接失败

**注意**: Redis是可选的，连接失败不会阻止服务启动

**解决方案**:
1. 检查Redis服务是否启动
2. 检查 `REDIS_ADDR` 配置
3. 检查Redis密码配置

### 6. 依赖包下载失败

**错误信息**:
```
go: module xxx: Get "https://...": dial tcp: lookup xxx on 127.0.0.1:53: no such host
```

**解决方案**:
1. 检查网络连接
2. 配置Go代理（国内推荐）：
```bash
go env -w GOPROXY=https://goproxy.cn,direct
```

### 7. 编译错误: 找不到依赖包

**解决方案**:
```bash
go mod download
go mod tidy
```

### 8. 日志目录创建失败

**错误信息**:
```
创建日志目录失败: mkdir logs: no such file or directory
```

**解决方案**:
```bash
# Windows
mkdir logs

# Linux/Mac
mkdir -p logs
```

---

## 开发环境配置

### 1. 设置Go代理（国内推荐）

```bash
go env -w GOPROXY=https://goproxy.cn,direct
go env -w GOSUMDB=sum.golang.google.cn
```

### 2. 检查Go环境

```bash
go version
go env GOPATH
go env GOROOT
```

### 3. 验证项目配置

```bash
# 检查依赖
go list -m all

# 检查编译
go build

# 运行测试
go test ./...
```

---

## 完整的依赖问题解决流程

### 步骤1: 清理所有缓存
```bash
go clean -modcache
```

### 步骤2: 修复有问题的依赖包
```bash
go get -u golang.org/x/tools@v0.20.0
```

### 步骤3: 重新整理依赖
```bash
go mod tidy
```

### 步骤4: 验证编译
```bash
go build -o pet-service.exe main.go
```

---

## 性能优化建议

### 1. 减少依赖

定期清理未使用的依赖：
```bash
go mod tidy
```

### 2. 使用Go Modules缓存

Go会自动缓存模块，无需额外配置

### 3. 编译优化

发布时使用优化编译：
```bash
go build -ldflags="-s -w" -o pet-service main.go
```

---

## 调试技巧

### 1. 启用详细日志

```bash
# Linux/Mac
export LOG_LEVEL=debug
go run main.go

# Windows
set LOG_LEVEL=debug
go run main.go
```

### 2. 使用Delve调试器

```bash
go install github.com/go-delve/delve/cmd/dlv@latest
dlv debug main.go
```

### 3. 性能分析

```bash
go test -cpuprofile=cpu.prof -memprofile=mem.prof
go tool pprof cpu.prof
```

---

## 常用命令参考

### 依赖管理
```bash
go mod init        # 初始化模块
go mod tidy        # 整理依赖
go mod download    # 下载依赖
go mod verify      # 验证依赖
go get -u xxx      # 更新依赖
go get xxx@version # 指定版本安装
```

### 构建和运行
```bash
go run main.go                    # 直接运行
go build -o app main.go           # 编译
go build -ldflags="-s -w"         # 优化编译
go install ./...                  # 安装到GOPATH
```

### 测试
```bash
go test ./...                     # 运行所有测试
go test -v ./...                  # 详细输出
go test -race ./...               # 竞态检测
go test -cover ./...              # 覆盖率测试
```

---

## 获取帮助

- [Hertz文档](https://www.cloudwego.io/docs/hertz/)
- [Go官方文档](https://go.dev/doc/)
- [GORM文档](https://gorm.io/zh_CN/docs/)
- [Go社区](https://go.dev/community/)

---

## 项目特定问题

### Hertz相关

如果遇到Hertz相关编译问题：
```bash
# 更新Hertz到最新版本
go get -u github.com/cloudwego/hertz
go mod tidy
```

### GORM相关

如果遇到GORM相关编译问题：
```bash
# 更新GORM相关包
go get -u gorm.io/gorm
go get -u gorm.io/gen
go mod tidy
```

---

## 编译成功验证

运行以下命令验证项目编译成功：

```bash
# 编译主程序
go build -o pet-service.exe main.go

# 编译测试程序
go build -o simple_test.exe simple_test.go

# 检查退出码
echo %errorlevel%
# 应该输出: 0

# 验证文件存在
dir pet-service.exe simple_test.exe
```

如果所有命令都成功执行且退出码为0，说明编译问题已完全解决。
