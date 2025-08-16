package middleware

import (
	"context"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
)

// 定义一个接口，任何实现了 Validate() error 的结构体都满足这个接口
type validator interface {
	Validate() error
}

// Validation 是一个中间件，用于自动校验请求参数
func Validation() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			// 检查请求参数是否实现了 validator 接口
			if v, ok := req.(validator); ok {
				// 如果实现了，就调用 Validate 方法
				if err := v.Validate(); err != nil {
					// 如果校验失败，返回一个 BadRequest 错误
					// 这样客户端会收到清晰的错误信息和正确的 HTTP 状态码（400）
					return nil, errors.BadRequest("INVALID_ARGUMENT", err.Error())
				}
			}
			// 校验通过或请求参数没有实现 validator 接口，则继续执行下一个处理函数
			return handler(ctx, req)
		}
	}
}
