package logger

import "log/slog"

type CronLogger struct {
	logger *slog.Logger
}

func NewCronLogger(logger *slog.Logger) *CronLogger {
	return &CronLogger{
		logger: logger,
	}
}

func (l *CronLogger) Info(msg string, args ...any) {
	l.logger.Info(msg, args...)
}

func (l *CronLogger) Error(err error, msg string, args ...any) {
	l.logger.Error(msg, append(args, slog.Any("error", err))...)
}
