package log

import (
	"context"
	"fmt"
	"github.com/lmittmann/tint"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
)

var (
	logger  *slog.Logger
	handler slog.Handler
)

type DynamicHandler struct {
	errWriter slog.Handler
	appWriter slog.Handler
	level     slog.LevelVar
}

func newDynamicHandler(errWriter, appWriter slog.Handler, level slog.Level) *DynamicHandler {
	h := &DynamicHandler{
		errWriter: errWriter,
		appWriter: appWriter,
	}
	h.SetLevel(level)
	return h
}

// 设置日志级别
func (h *DynamicHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.level.Level() >= level
}

func (h *DynamicHandler) Handle(ctx context.Context, record slog.Record) error {
	if record.Level >= slog.LevelError {
		return h.errWriter.Handle(ctx, record)
	}
	return h.appWriter.Handle(ctx, record)
}

func (h *DynamicHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	// TODO
	return nil
}

func (h *DynamicHandler) WithGroup(name string) slog.Handler {
	// TODO
	return nil
}

func (h *DynamicHandler) SetLevel(level slog.Level) {
	h.level.Set(level)
}

// 开发模式
func devHandler() slog.Handler {
	appWriter := os.Stdout
	return tint.NewHandler(appWriter, &tint.Options{
		Level:      slog.LevelDebug,
		TimeFormat: "2006-01-02 15:04:05",
	})
}

// 生产模式
func prodHandler() slog.Handler {
	// TODO 修改成配置
	// 多个输出
	errWriter := &lumberjack.Logger{
		Filename:   "",
		MaxSize:    0,
		MaxBackups: 0,
		MaxAge:     0, // days
	}

	// 普通日志
	appWriter := &lumberjack.Logger{
		Filename:   "",
		MaxSize:    0,
		MaxBackups: 0,
		MaxAge:     0, // days
	}

	return newDynamicHandler(
		slog.NewJSONHandler(errWriter, &slog.HandlerOptions{
			Level: slog.LevelError,
		}),
		slog.NewJSONHandler(io.MultiWriter(appWriter, os.Stdout), &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}),
		slog.LevelInfo,
	)
}

// 追加个性化参数，例如 trace_id
func addParams(ctx context.Context, args ...any) []any {
	// 添加 trace_id
	if ctx == nil {
		ctx = context.Background()
	}
	if traceID := ctx.Value("trace_id"); traceID != nil {
		args = append(args, "trace_id", traceID)
	}
	// TODO 如果是 gorm.go则需要调整 skip
	args = append(args, "caller", caller(3))
	return args
}

// 获取调用者信息，输出文件名和行号
func caller(skip int) string {
	pc, file, line, _ := runtime.Caller(skip)
	if filepath.Base(file) == "log_db.go" {
		pc, file, line, _ = runtime.Caller(skip + 6)
	}
	fn := runtime.FuncForPC(pc)
	return fmt.Sprintf("%s %s:%d", fn.Name(), filepath.Base(file), line)
}

func Info(ctx context.Context, msg string, args ...any) {
	logger.InfoContext(ctx, msg, addParams(ctx, args...)...)
}

func Warn(ctx context.Context, msg string, args ...any) {
	logger.WarnContext(ctx, msg, addParams(ctx, args...)...)
}

func Error(ctx context.Context, msg string, args ...any) {
	logger.ErrorContext(ctx, msg, addParams(ctx, args...)...)
}

func Debug(ctx context.Context, msg string, args ...any) {
	logger.DebugContext(ctx, msg, addParams(ctx, args...)...)
}

func Fatal(ctx context.Context, msg string, args ...any) {
	logger.ErrorContext(ctx, msg, addParams(ctx, args...)...)
	// 退出程序
	os.Exit(1)
}

// 根据不同环境返回不同的日志记录器
func newSlogLogger() *slog.Logger {
	var (
		env string
	)
	// TODO 修改成配置
	env = os.Getenv("ENV")

	switch env {
	case "dev":
		handler = devHandler()
	case "prod":
		handler = prodHandler()
	default:
		handler = devHandler()
	}

	return slog.New(handler)
}

func init() {
	logger = newSlogLogger()
}
