package log

import (
	"context"
	"testing"
)

func TestLog(t *testing.T) {
	ctx := context.WithValue(context.Background(), "trace_id", "tick-tock")
	Debug(ctx, "debug message", "hello", "world")
	Info(ctx, "info message", "hello", "world")
	Warn(ctx, "warn message", "hello", "world")
	Error(ctx, "error message", "hello", "world")
	Fatal(ctx, "fatal message", "hello", "world")
}
