package pgxlog

import (
	"context"

	"github.com/diontr00/logStack"
	"github.com/jackc/pgx/v5/tracelog"
)

type Logger struct {
	logger *logStack.Logger
}

func NewPgxLogger(logger *logStack.Logger) *Logger {
	return &Logger{logger: logger}
}

func (l *Logger) Log(
	ctx context.Context,
	tracelv tracelog.LogLevel,
	msg string,
	log map[string]any,
) {
	log_fields := make([]logStack.LogField, len(log))

	var i = 0
	for key, val := range log {
		log_fields[i] = logStack.Any(key, val)
		i++
	}

	switch tracelv {
	case tracelog.LogLevelDebug:
		l.logger.Debug(msg, log_fields...)
	case tracelog.LogLevelInfo:
		l.logger.Info(msg, log_fields...)
	case tracelog.LogLevelWarn:
		l.logger.Warn(msg, log_fields...)
	case tracelog.LogLevelError:
		l.logger.Error(msg, log_fields...)
	case tracelog.LogLevelNone:
		l.logger.Debug(msg, log_fields...)
	case tracelog.LogLevelTrace:
		l.logger.Debug(msg, append(log_fields, logStack.Stringer("PGX_LOG_LEVEL", tracelv))...)
	default:
		l.logger.Debug(msg, append(log_fields, logStack.Stringer("PGX_LOG_LEVEL", tracelv))...)

	}

}
