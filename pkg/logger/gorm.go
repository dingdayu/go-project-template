// pkg/logger/gorm.go
package logger

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	gormLogger "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

// GormLogger 是一个适配 slog 的 GORM 日志实现
type GormLogger struct {
	log   *slog.Logger
	level gormLogger.LogLevel
}

// NewGormLogger 创建一个新的 GORM 日志适配器
func NewGormLogger(log *slog.Logger) gormLogger.Interface {
	if log == nil {
		log = logger
	}
	return &GormLogger{
		log:   log,
		level: gormLogger.Warn, // 默认级别
	}
}

// LogMode 实现 gorm.Logger 接口：设置日志级别
func (l *GormLogger) LogMode(level gormLogger.LogLevel) gormLogger.Interface {
	newLogger := *l
	newLogger.level = level
	return &newLogger
}

// Info 实现接口
func (l *GormLogger) Info(ctx context.Context, msg string, args ...interface{}) {
	if l.level >= gormLogger.Info {
		l.log.InfoContext(ctx, fmt.Sprintf(msg, args...))
	}
}

// Warn 实现接口
func (l *GormLogger) Warn(ctx context.Context, msg string, args ...interface{}) {
	if l.level >= gormLogger.Warn {
		l.log.WarnContext(ctx, fmt.Sprintf(msg, args...))
	}
}

// Error 实现接口
func (l *GormLogger) Error(ctx context.Context, msg string, args ...interface{}) {
	if l.level >= gormLogger.Error {
		l.log.ErrorContext(ctx, fmt.Sprintf(msg, args...))
	}
}

// Trace 实现接口：记录 SQL 执行详情
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()

	switch {
	case err != nil && l.level >= gormLogger.Error:
		l.log.ErrorContext(ctx, sql,
			slog.String("source", utils.FileWithLineNum()),
			slog.Duration("du", elapsed),
			slog.Int64("rows", rows),
			slog.Any("error", err),
		)
	case l.level >= gormLogger.Info:
		l.log.InfoContext(ctx, sql,
			slog.String("source", utils.FileWithLineNum()),
			slog.Duration("du", elapsed),
			slog.Int64("rows", rows),
		)
	}
}
