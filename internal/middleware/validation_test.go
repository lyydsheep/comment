package middleware

import (
	"context"
	"errors"
	"testing"

	kratoserrors "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/stretchr/testify/assert"
)

// 定义一个测试用的实现了validator接口的结构体
// 这个结构体用于测试验证通过的场景
type validRequest struct {
	Field string
}

func (v *validRequest) Validate() error {
	return nil // 验证通过
}

// 定义一个测试用的实现了validator接口的结构体
// 这个结构体用于测试验证失败的场景
type invalidRequest struct {
	Field string
}

func (i *invalidRequest) Validate() error {
	return errors.New("validation failed") // 验证失败
}

// 定义一个测试用的没有实现validator接口的结构体
type nonValidatorRequest struct {
	Field string
}

// 定义一个用于测试自定义错误消息的结构体
type customErrorRequest struct {
	ErrorMessage string
}

func (c *customErrorRequest) Validate() error {
	return errors.New(c.ErrorMessage)
}

// MockReply 是一个简单的回复结构体
// 用于模拟RPC调用的返回值
type MockReply struct {
	Success bool
	Message string
}

// MockService 是一个模拟的服务接口
// 用于测试RPC中间件的调用

type MockService interface {
	// TestRPCMethod 是一个模拟的RPC方法
	TestRPCMethod(ctx context.Context, req interface{}) (interface{}, error)
}

// MockServiceImpl 是MockService接口的实现
// 用于模拟RPC服务的处理逻辑
type MockServiceImpl struct{}

// TestRPCMethod 实现MockService接口的方法
// 简单返回一个成功的响应
func (s *MockServiceImpl) TestRPCMethod(ctx context.Context, req interface{}) (interface{}, error) {
	return &MockReply{
		Success: true,
		Message: "RPC call successful",
	}, nil
}

func TestValidationMiddlewareWithRPC(t *testing.T) {
	// 创建模拟服务实例
	service := &MockServiceImpl{}

	// 创建一个处理函数，模拟RPC方法调用
	handler := middleware.Handler(func(ctx context.Context, req interface{}) (interface{}, error) {
		return service.TestRPCMethod(ctx, req)
	})

	// 创建验证中间件
	validationMw := Validation()
	handlerWithMw := validationMw(handler)

	// 创建一个上下文
	ctx := context.Background()

	// 测试场景1：请求对象实现了validator接口并且验证通过
	t.Run("ValidRequestWithRPC", func(t *testing.T) {
		req := &validRequest{Field: "valid"}
		reply, err := handlerWithMw(ctx, req)
		assert.NoError(t, err)
		assert.NotNil(t, reply)

		mockReply, ok := reply.(*MockReply)
		assert.True(t, ok)
		assert.True(t, mockReply.Success)
		assert.Equal(t, "RPC call successful", mockReply.Message)
	})

	// 测试场景2：请求对象实现了validator接口但验证失败
	t.Run("InvalidRequestWithRPC", func(t *testing.T) {
		req := &invalidRequest{Field: "invalid"}
		reply, err := handlerWithMw(ctx, req)
		assert.Error(t, err)
		assert.Nil(t, reply)

		// 验证错误类型和内容
		kratosErr, ok := err.(*kratoserrors.Error)
		assert.True(t, ok)
		assert.Equal(t, int32(400), kratosErr.Code)
		assert.Equal(t, "INVALID_ARGUMENT", kratosErr.Reason)
		assert.Contains(t, kratosErr.Message, "validation failed")
	})

	// 测试场景3：请求对象没有实现validator接口
	t.Run("NonValidatorRequestWithRPC", func(t *testing.T) {
		req := &nonValidatorRequest{Field: "test"}
		reply, err := handlerWithMw(ctx, req)
		assert.NoError(t, err)
		assert.NotNil(t, reply)

		mockReply, ok := reply.(*MockReply)
		assert.True(t, ok)
		assert.True(t, mockReply.Success)
	})

	// 测试场景4：nil请求
	t.Run("NilRequestWithRPC", func(t *testing.T) {
		reply, err := handlerWithMw(ctx, nil)
		assert.NoError(t, err)
		assert.NotNil(t, reply)

		mockReply, ok := reply.(*MockReply)
		assert.True(t, ok)
		assert.True(t, mockReply.Success)
	})

	// 测试场景5：非结构体类型的请求（如基本类型）
	t.Run("BasicTypeRequestWithRPC", func(t *testing.T) {
		req := "string request"
		reply, err := handlerWithMw(ctx, req)
		assert.NoError(t, err)
		assert.NotNil(t, reply)

		mockReply, ok := reply.(*MockReply)
		assert.True(t, ok)
		assert.True(t, mockReply.Success)
	})

	// 测试场景6：带有不同验证错误消息的请求
	t.Run("CustomErrorMessageWithRPC", func(t *testing.T) {
		customErrMsg := "custom validation error"
		req := &customErrorRequest{ErrorMessage: customErrMsg}
		reply, err := handlerWithMw(ctx, req)
		assert.Error(t, err)
		assert.Nil(t, reply)

		// 验证错误类型和内容
		kratosErr, ok := err.(*kratoserrors.Error)
		assert.True(t, ok)
		assert.Equal(t, int32(400), kratosErr.Code)
		assert.Equal(t, "INVALID_ARGUMENT", kratosErr.Reason)
		assert.Contains(t, kratosErr.Message, customErrMsg)
	})

	// 测试场景7：中间件链的正确执行顺序
	t.Run("MiddlewareChainWithRPC", func(t *testing.T) {
		// 创建一个捕获调用顺序的中间件链
		callOrder := []string{}

		// 第一个中间件
		mw1 := func(handler middleware.Handler) middleware.Handler {
			return func(ctx context.Context, req interface{}) (interface{}, error) {
				callOrder = append(callOrder, "mw1")
				return handler(ctx, req)
			}
		}

		// 第二个中间件是我们的验证中间件
		mw2 := Validation()

		// 第三个中间件
		mw3 := func(handler middleware.Handler) middleware.Handler {
			return func(ctx context.Context, req interface{}) (interface{}, error) {
				callOrder = append(callOrder, "mw3")
				return handler(ctx, req)
			}
		}

		// 创建中间件链
		chain := middleware.Chain(mw1, mw2, mw3)(handler)

		// 使用验证通过的请求调用中间件链
		req := &validRequest{Field: "valid"}
		reply, err := chain(ctx, req)

		// 验证调用顺序和结果
		assert.NoError(t, err)
		assert.NotNil(t, reply)
		// 注意：验证中间件在验证通过时不会添加到调用顺序中
		assert.Equal(t, []string{"mw1", "mw3"}, callOrder)

		mockReply, ok := reply.(*MockReply)
		assert.True(t, ok)
		assert.True(t, mockReply.Success)
	})

	// 测试场景8：使用普通上下文测试验证中间件的通用性
	t.Run("GRPCCtxWithRPC", func(t *testing.T) {
		// 使用标准上下文替代gRPC上下文，因为我们主要测试中间件逻辑
		// 验证中间件在不同类型的上下文中都能正常工作
		grpcCtx := context.WithValue(context.Background(), "grpc", "test-context")
		req := &validRequest{Field: "grpc context test"}
		reply, err := handlerWithMw(grpcCtx, req)
		assert.NoError(t, err)
		assert.NotNil(t, reply)

		mockReply, ok := reply.(*MockReply)
		assert.True(t, ok)
		assert.True(t, mockReply.Success)
	})

	// 测试场景9：验证中间件在验证失败时不会调用下一个处理器
	t.Run("ValidationFailureDoesNotCallNextHandler", func(t *testing.T) {
		called := false

		// 创建一个处理器，如果被调用就设置called为true
		handler := middleware.Handler(func(ctx context.Context, req interface{}) (interface{}, error) {
			called = true
			return nil, nil
		})

		// 创建验证中间件和处理器
		handlerWithMw := Validation()(handler)

		// 使用验证失败的请求
		req := &invalidRequest{Field: "invalid"}
		handlerWithMw(ctx, req)

		// 验证下一个处理器没有被调用
		assert.False(t, called, "Next handler should not be called when validation fails")
	})
}
