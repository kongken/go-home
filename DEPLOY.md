# Docker 部署指南

## 快速开始

### 1. 构建并启动所有服务

```bash
cd /home/node/workspace/go-home
docker-compose up -d --build
```

### 2. 查看服务状态

```bash
docker-compose ps
```

### 3. 查看日志

```bash
# 所有服务
docker-compose logs -f

# 特定服务
docker-compose logs -f backend
docker-compose logs -f frontend
```

### 4. 停止服务

```bash
docker-compose down
```

### 5. 完全清理（包括数据卷）

```bash
docker-compose down -v
```

## 服务说明

| 服务 | 端口 | 说明 |
|------|------|------|
| backend | 2222 | Go 后端 API 服务 |
| frontend | 80 | React 前端（Nginx） |
| mysql | 3306 | MySQL 数据库 |
| redis | 6379 | Redis 缓存 |
| prometheus | 9090 | Prometheus 监控 |
| grafana | 3001 | Grafana 可视化 |

## 访问地址

- 前端: http://localhost
- 后端 API: http://localhost:2222
- Prometheus: http://localhost:9090
- Grafana: http://localhost:3001 (admin/admin)

## 单独构建

### 后端

```bash
cd /home/node/workspace/go-home
docker build -t go-home-backend .
docker run -p 2222:2222 go-home-backend
```

### 前端

```bash
cd /home/node/workspace/go-home/front
docker build -t go-home-frontend .
docker run -p 80:80 go-home-frontend
```

## 环境变量

### 后端

| 变量 | 默认值 | 说明 |
|------|--------|------|
| SERVER_PORT | 2222 | 服务端口 |
| SERVER_MODE | release | 运行模式 |
| DB_HOST | mysql | 数据库主机 |
| DB_PORT | 3306 | 数据库端口 |
| DB_USER | root | 数据库用户 |
| DB_PASSWORD | password | 数据库密码 |
| DB_NAME | go_home | 数据库名 |
| REDIS_HOST | redis | Redis 主机 |
| REDIS_PORT | 6379 | Redis 端口 |
| JWT_SECRET | your-secret-key | JWT 密钥 |

## 生产部署建议

1. 修改默认密码
2. 使用 HTTPS
3. 配置适当的资源限制
4. 设置日志轮转
5. 配置监控告警
