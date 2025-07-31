package log

import (
	"context"
	glogger "gorm.io/gorm/logger"
	"log/slog"
	"reflect"
	"time"
)

type DBLogger struct{}

var (
	m = map[glogger.LogLevel]slog.Level{glogger.Info: slog.LevelInfo, glogger.Warn: slog.LevelWarn, glogger.Error: slog.LevelError}
)

func (logger *DBLogger) LogMode(level glogger.LogLevel) glogger.Interface {
	dynamicHandler, ok := handler.(*DynamicHandler)
	if !ok {
		Warn(nil, "fail to set log mode.", "type of handler", reflect.TypeOf(handler).String())
		return logger
	}
	dynamicHandler.SetLevel(m[level])
	Info(nil, "set log mode successfully.")
	return logger
}

func (logger *DBLogger) Info(ctx context.Context, s string, i ...interface{}) {
	Info(ctx, s, i...)
}

func (logger *DBLogger) Warn(ctx context.Context, s string, i ...interface{}) {
	Warn(ctx, s, i...)
}

func (logger *DBLogger) Error(ctx context.Context, s string, i ...interface{}) {
	Error(ctx, s, i...)
}

func (logger *DBLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	sql, rowsAffected := fc()
	elapsed := time.Since(begin)
	if elapsed > time.Millisecond*200 {
		Warn(ctx, "slow query",
			"sql", sql,
			"rows_affected", rowsAffected,
			"duration_ms", elapsed.Milliseconds(),
			"error", err)
		return
	}
	Info(ctx, "database query.",
		"sql", sql,
		"rows_affected", rowsAffected,
		"duration_ms", elapsed.Milliseconds(),
		"error", err)
}

func Logger() *DBLogger {
	return &DBLogger{}
}
