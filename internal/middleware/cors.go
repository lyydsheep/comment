package middleware

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	ktHttp "github.com/go-kratos/kratos/v2/transport/http"
)

// CORSOptions 定义CORS配置选项
type CORSOptions struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	AllowCredentials bool
	MaxAge           int
}

// DefaultCORSOptions 默认CORS配置，包含常见的本地开发地址
var DefaultCORSOptions = CORSOptions{
	AllowOrigins: []string{
		"*",
		"http://localhost:*",
		"http://127.0.0.1:*",
		"https://localhost:*",
		"https://127.0.0.1:*",
	},
	AllowMethods: []string{
		"GET",
		"HEAD",
		"PUT",
		"PATCH",
		"POST",
		"DELETE",
		"OPTIONS",
	},
	AllowHeaders: []string{
		"Origin",
		"Content-Length",
		"Content-Type",
		"Authorization",
		"X-Requested-With",
		"Accept",
		"X-CSRF-Token",
	},
	AllowCredentials: true,
	MaxAge:           86400,
}

// CORS 创建CORS中间件
func CORS(opts ...CORSOptions) middleware.Middleware {
	options := DefaultCORSOptions
	if len(opts) > 0 {
		options = opts[0]
	}

	// 预处理允许的源，将通配符模式转换为可比较的形式
	allowOrigins := make([]string, 0, len(options.AllowOrigins))
	allowOriginPatterns := make([]string, 0, len(options.AllowOrigins))
	
	for _, origin := range options.AllowOrigins {
		if origin == "*" {
			allowOrigins = append(allowOrigins, origin)
		} else if strings.Contains(origin, "*") {
			// 将 "http://localhost:*" 转换为 "http://localhost:"
			pattern := strings.TrimSuffix(origin, "*")
			allowOriginPatterns = append(allowOriginPatterns, pattern)
		} else {
			allowOrigins = append(allowOrigins, origin)
		}
	}

	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			// 检查是否是HTTP传输
			if tr, ok := transport.FromServerContext(ctx); ok {
				if ht, ok := tr.(ktHttp.Transporter); ok {
					// 获取请求的Origin
					requestOrigin := ht.Request().Header.Get("Origin")
					
					// 设置CORS响应头
					header := ht.ReplyHeader()
					
					// 确定应该设置哪个Access-Control-Allow-Origin值
					allowOrigin := ""
					if len(allowOrigins) > 0 && allowOrigins[0] == "*" {
						// 如果允许所有源，直接设置为请求的Origin或者"*"
						if options.AllowCredentials {
							// 当允许凭证时，不能使用"*"，需要使用具体的Origin
							allowOrigin = requestOrigin
						} else {
							allowOrigin = "*"
						}
					} else if requestOrigin != "" {
						// 检查请求的Origin是否在允许列表中
						for _, origin := range allowOrigins {
							if origin == requestOrigin {
								allowOrigin = requestOrigin
								break
							}
						}
						
						// 如果没有精确匹配，检查模式匹配
						if allowOrigin == "" {
							for _, pattern := range allowOriginPatterns {
								if strings.HasPrefix(requestOrigin, pattern) {
									allowOrigin = requestOrigin
									break
								}
							}
						}
					}
					
					// 只有在找到匹配的Origin时才设置Access-Control-Allow-Origin
					if allowOrigin != "" {
						header.Set("Access-Control-Allow-Origin", allowOrigin)
					}
					
					// 设置其他CORS头部
					allowMethodsStr := strings.Join(options.AllowMethods, ", ")
					allowHeadersStr := strings.Join(options.AllowHeaders, ", ")
					
					header.Set("Access-Control-Allow-Methods", allowMethodsStr)
					header.Set("Access-Control-Allow-Headers", allowHeadersStr)
					
					if options.AllowCredentials {
						header.Set("Access-Control-Allow-Credentials", "true")
					}
					
					if options.MaxAge > 0 {
						header.Set("Access-Control-Max-Age", strconv.Itoa(options.MaxAge))
					}

					// 如果是OPTIONS请求，直接返回
					if ht.Request().Method == http.MethodOptions {
						return nil, nil
					}
				}
			}

			// 继续处理请求
			return handler(ctx, req)
		}
	}
}