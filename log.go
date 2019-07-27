package elog

import "context"

// 对外提供统一接口，可自定义替换
// 默认使用elog
type Logger interface {
	Debug(ctx context.Context, msg interface{}, kv ...interface{})
	Info(ctx context.Context, msg interface{}, kv ...interface{})
	Warn(ctx context.Context, msg interface{}, kv ...interface{})
	Error(ctx context.Context, msg interface{}, kv ...interface{})
	Fatal(ctx context.Context, msg interface{}, kv ...interface{}) // 这个里面会panic
	With(ctx context.Context, kv ...interface{}) context.Context   // 增量附加字段 以后的日志都会带上这个日志
	Close() error
	EnableDebug(b bool)
}

var __eLogger Logger

var __eLoggerErrorAbove Logger

func SetLogger(l Logger) {
	__eLogger = l
}

func GetLogger() Logger {
	if __eLogger == nil {
		SetLogger(GetJsonELog())
	}
	return __eLogger
}

func SetLoggerErrorAbove(l Logger) {
	__eLoggerErrorAbove = l
}

func GetLoggerErrorAbove() Logger {
	if __eLoggerErrorAbove == nil {
		SetLogger(GetJsonELogErrorAbove())
	}
	return __eLoggerErrorAbove
}

func Debug(ctx context.Context, msg interface{}, kv ...interface{}) {
	GetLogger().Debug(ctx, msg, kv...)
}

func Info(ctx context.Context, msg interface{}, kv ...interface{}) {
	GetLogger().Info(ctx, msg, kv...)
}

func Warn(ctx context.Context, msg interface{}, kv ...interface{}) {
	GetLogger().Warn(ctx, msg, kv...)
}

func Error(ctx context.Context, msg interface{}, kv ...interface{}) {
	GetLogger().Error(ctx, msg, kv...)
	GetLoggerErrorAbove().Error(ctx, msg, kv...)
}

func Fatal(ctx context.Context, msg interface{}, kv ...interface{}) {
	GetLogger().Fatal(ctx, msg, kv...)
	GetLoggerErrorAbove().Error(ctx, msg, kv...)
}

func With(ctx context.Context, kv ...interface{}) context.Context {
	return GetLogger().With(ctx, kv...)
}

// 这个方法以后不要用了，请使用Close()
func Flush() error {
	return Close()
}

func Close() error {
	GetLoggerErrorAbove().Close()
	return GetLogger().Close()
}

func EnableDebug(b bool) {
	GetLogger().EnableDebug(b)
}
