# Go Home

基于 Go/Gin/Protobuf 的社交网络后端服务。

## 技术栈

- **框架**: Gin
- **数据库**: PostgreSQL + GORM
- **缓存**: Redis
- **认证**: JWT
- **API 协议**: RESTful + Protobuf

## 项目结构

```
.
├── cmd/server/          # 应用程序入口
├── internal/            # 内部代码
│   ├── config/         # 配置管理
│   ├── handler/        # HTTP 处理器
│   ├── middleware/     # 中间件
│   ├── model/          # 数据模型
│   ├── repository/     # 数据访问层
│   └── service/        # 业务逻辑层
├── pkg/                # 公共包
│   ├── cache/          # Redis 缓存
│   ├── database/       # 数据库连接
│   └── proto/          # Protobuf 生成的代码
├── proto/              # Protobuf 定义文件
├── config.yaml         # 配置文件
├── docker-compose.yml  # Docker Compose 配置
└── Dockerfile          # Docker 镜像构建
```

## 快速开始

### 1. 安装依赖

```bash
go mod download
```

### 2. 配置环境

编辑 `config.yaml` 文件，配置数据库和 Redis 连接信息。

### 3. 启动依赖服务

```bash
docker-compose up -d postgres redis
```

### 4. 运行应用

```bash
make run
# 或
go run cmd/server/main.go
```

### 5. 使用 Docker 运行完整环境

```bash
docker-compose up -d
```

## API 接口

### 认证

- `POST /api/v1/auth/register` - 用户注册
- `POST /api/v1/auth/login` - 用户登录
- `POST /api/v1/auth/refresh` - 刷新 Token

### 用户

- `GET /api/v1/users/:id` - 获取用户信息
- `PUT /api/v1/users/me` - 更新当前用户信息

### 博客

- `GET /api/v1/blogs` - 获取博客列表
- `GET /api/v1/blogs/:id` - 获取博客详情
- `GET /api/v1/users/:user_id/blogs` - 获取用户博客列表
- `POST /api/v1/blogs` - 创建博客（需认证）
- `PUT /api/v1/blogs/:id` - 更新博客（需认证）
- `DELETE /api/v1/blogs/:id` - 删除博客（需认证）

## 开发

### 生成 Protobuf 代码

```bash
make proto
```

### 运行测试

```bash
make test
```

### 代码格式化

```bash
make fmt
```

## 配置说明

| 配置项 | 说明 | 默认值 |
|--------|------|--------|
| server.host | 服务器监听地址 | 0.0.0.0 |
| server.port | 服务器端口 | 8080 |
| server.mode | Gin 模式 (debug/release) | debug |
| database.host | 数据库主机 | localhost |
| database.port | 数据库端口 | 5432 |
| redis.host | Redis 主机 | localhost |
| redis.port | Redis 端口 | 6379 |
| jwt.secret | JWT 密钥 | your-secret-key |
| jwt.access_expiry | Access Token 过期时间(秒) | 3600 |
| jwt.refresh_expiry | Refresh Token 过期时间(秒) | 604800 |
