package logging

import (
	"context"
)

type Logger interface {
	Named(name string) Logger
	With(args ...interface{}) Logger
	WithContext(ctx context.Context) Logger
	Debug(message string, args ...interface{})
	Info(message string, args ...interface{})
	Warn(message string, args ...interface{})
	Error(message string, args ...interface{})
	Fatal(message string, args ...interface{})
}
