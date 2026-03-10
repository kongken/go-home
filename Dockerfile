# 构建阶段
FROM golang:1.25-alpine AS builder

WORKDIR /app

# 安装 git 和 ca-certificates
RUN apk add --no-cache git ca-certificates

# 复制 go mod 文件
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 构建
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o server ./cmd/server

# 运行阶段
FROM alpine:latest

WORKDIR /app

# 安装 ca-certificates
RUN apk --no-cache add ca-certificates

# 从构建阶段复制二进制文件
COPY --from=builder /app/server .

# 暴露端口
EXPOSE 2222

# 运行
CMD ["./server"]
