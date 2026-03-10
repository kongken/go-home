# Butterfly 框架迁移任务

## 目标
将当前 go-home 项目从原生 Gin + MongoDB/Redis 架构迁移到 butterfly 微服务框架

## Butterfly 框架特性
- **配置管理**: 文件/Consul 配置中心
- **服务运行时**: 应用生命周期管理、优雅关闭
- **传输层**: HTTP (Gin)、gRPC、Twirp RPC
- **数据存储**: GORM、MongoDB v2、Redis、原生 SQL
- **可观测性**: Prometheus 指标、OpenTelemetry 链路追踪、结构化日志
- **依赖注入**: Google Wire 支持

## 当前项目状态
- 10 个服务: User, Blog, Feed, Friend, Group, Message, Notification, Comment, Album, Settings
- 60+ API 端点
- MongoDB + Redis 存储
- JWT 认证
- 分层架构: Handler → Service → Repository

## 迁移计划

### Phase 1: 基础架构迁移
- [ ] 添加 butterfly 依赖
- [ ] 创建 butterfly 配置文件 (config.yaml)
- [ ] 重构 main.go 使用 butterfly app
- [ ] 配置 MongoDB/Redis 连接

### Phase 2: 存储层迁移
- [ ] 替换 MongoDB 连接为 butterfly mongo store
- [ ] 替换 Redis 连接为 butterfly redis store
- [ ] 更新 Repository 层使用 butterfly store API

### Phase 3: 服务层迁移
- [ ] 集成 butterfly 日志系统
- [ ] 添加 OpenTelemetry 链路追踪
- [ ] 添加 Prometheus 指标暴露

### Phase 4: 依赖注入重构
- [ ] 使用 Google Wire 重构依赖注入
- [ ] 创建 provider 函数
- [ ] 生成 wire_gen.go

### Phase 5: 传输层增强
- [ ] 添加 gRPC 服务支持
- [ ] 添加 Twirp RPC 支持
- [ ] 配置多协议端口

### Phase 6: 可观测性
- [ ] 配置 OpenTelemetry 导出
- [ ] 配置 Prometheus 指标
- [ ] 添加健康检查端点

## 新目录结构
```
go-home/
├── cmd/
│   └── server/
│       ├── main.go
│       └── wire.go           # Wire 注入器
├── internal/
│   ├── config/               # 配置结构
│   ├── handler/              # HTTP Handlers
│   ├── service/              # 业务逻辑
│   ├── repository/           # 数据访问
│   ├── model/                # 数据模型
│   └── wire/                 # Wire providers
├── pkg/                      # 公共包
├── proto/                    # Protobuf
├── config.yaml               # Butterfly 配置
├── Dockerfile
└── docker-compose.yml
```

## 关键变更

### 1. 依赖替换
```go
// 旧
"github.com/kongken/go-home/pkg/database"
"github.com/kongken/go-home/pkg/cache"

// 新
"butterfly.orx.me/core/store/mongo"
"butterfly.orx.me/core/store/redis"
"butterfly.orx.me/core/app"
"butterfly.orx.me/core/log"
```

### 2. Repository 模式变更
```go
// 旧
func (r *userRepository) GetByID(ctx context.Context, id string) (*model.User, error)

// 新 - 使用 butterfly mongo API
func (r *userRepository) GetByID(ctx context.Context, id string) (*model.User, error) {
    collection := mongo.GetClient("primary").Database("gohome").Collection("users")
    // ...
}
```

### 3. 日志系统变更
```go
// 旧
"go.uber.org/zap"

// 新
"butterfly.orx.me/core/log"
"log/slog"

logger := log.FromContext(ctx)
logger.Info("user created", "user_id", user.ID)
```

### 4. 主程序重构
```go
// 旧
r := gin.New()
// ... 手动初始化所有依赖

// 新
config := &app.Config{
    Service: "go-home",
    Router: setupRoutes,
    InitFunc: []func() error{
        initDatabase,
        initCache,
    },
}
application := app.New(config)
application.Run()
```

## 环境变量配置
```bash
# 配置类型
export BUTTERFLY_CONFIG_TYPE=file
export BUTTERFLY_CONFIG_FILE_PATH=/path/to/config.yaml

# 链路追踪
export BUTTERFLY_TRACING_ENDPOINT=localhost:4318
export BUTTERFLY_TRACING_PROVIDER=http

# Prometheus
export BUTTERFLY_PROMETHEUS_PUSH_ENDPOINT=http://pushgateway:9091
```

## 配置示例 (config.yaml)
```yaml
store:
  mongo:
    primary:
      uri: "mongodb://localhost:27017"
  redis:
    cache:
      addr: "localhost:6379"
      password: ""
      db: 0

otel:
  # OpenTelemetry 配置
```

## 端口规划
- 8080: HTTP (Gin)
- 9090: gRPC
- 2223: Prometheus Metrics

## 开始时间
2026-03-10

## 预计完成
待定
