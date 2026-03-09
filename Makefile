.PHONY: build run test clean proto docker

# 变量
BINARY_NAME=go-home
MAIN_FILE=cmd/server/main.go
BUILD_DIR=build

# 构建
build:
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_FILE)

# 运行
run:
	go run $(MAIN_FILE)

# 测试
test:
	go test -v ./...

# 清理
clean:
	rm -rf $(BUILD_DIR)

# 生成 protobuf
tproto:
	buf generate

# 下载依赖
deps:
	go mod download
	go mod tidy

# 格式化代码
fmt:
	go fmt ./...

# 代码检查
lint:
	golangci-lint run

# Docker 构建
docker-build:
	docker build -t $(BINARY_NAME):latest .

# Docker 运行
docker-run:
	docker-compose up -d

# Docker 停止
docker-stop:
	docker-compose down

# 数据库迁移
migrate:
	go run $(MAIN_FILE) migrate

# 开发模式（带热重载）
dev:
	air -c .air.toml
