package grpc

import (
	"context"

	"github.com/kongken/go-home/internal/handler"
	pb "github.com/kongken/go-home/pkg/proto/user"
	"google.golang.org/grpc"
)

// UserServer 用户 gRPC 服务
type UserServer struct {
	pb.UnimplementedUserServiceServer
	userHandler *handler.UserHandler
}

// NewUserServer 创建用户 gRPC 服务
func NewUserServer(userHandler *handler.UserHandler) *UserServer {
	return &UserServer{userHandler: userHandler}
}

// RegisterUserService 注册用户 gRPC 服务
func RegisterUserService(s *grpc.Server, userHandler *handler.UserHandler) {
	pb.RegisterUserServiceServer(s, NewUserServer(userHandler))
}

// GetUser 获取用户
func (s *UserServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
	// 这里简化处理，实际应该调用 handler 或 service 的方法
	return &pb.User{
		UserId:   req.UserId,
		Username: "test",
		Nickname: "Test User",
	}, nil
}