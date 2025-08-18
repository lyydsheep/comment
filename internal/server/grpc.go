package server

import (
	v1 "comment/api/comment/v1"
	"comment/internal/conf"
	"comment/internal/middleware"
	"comment/internal/service"

	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

// NewGRPCServer 创建并配置一个新的 gRPC 服务器实例
// 参数：
//
//	c - 服务器配置，包含 gRPC 相关设置
//	comment - 评论服务实现实例
//
// 返回：
//
//	配置好的 gRPC 服务器实例
func NewGRPCServer(c *conf.Server, comment *service.CommentService) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
			middleware.Validation(),
		),
	}
	if c.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Grpc.Network))
	}
	if c.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.Grpc.Addr))
	}
	if c.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Grpc.Timeout.AsDuration()))
	}

	// 创建 gRPC 服务器
	srv := grpc.NewServer(opts...)

	// 注册评论服务到服务器
	v1.RegisterCommentServiceServer(srv, comment)

	return srv
}
