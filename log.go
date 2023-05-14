// Provide the wrapper over uber zap structured logging , so that new logger instance dont  have to deal with zap core configuration\
// User can either consume the default logger with package level function  DefaultLogger
// Or define their own logger by using NewLogger
// ResetDefault should be called if the need to replace default logger , so that subsequent log will go to the right io.Writer taget
// Also come with multi location logger and log with rotate to complete the log stack requirement
package logStack

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"time"
)

type Loglevel = zapcore.Level
type LogOption = zap.Option

// This will be use  as the enabler function when register core
type LevelEnablerFunc func(lvl Loglevel) bool

type RotateOptions struct {

	// Maximum number of Megabytes before log old become it get rotated , default is 100 j
	MaxSize int
	// Maximum number of days to retain the old logs
	MaxAge int
	// Max number of olds log file  to retain, the default is to retain all log file , the default is all will be retain , even though it may get delete by maxAge
	MaxBackups int
	// Determine if log rotate will be compressed
	Compress bool
}

// Will be use as the destination to write multile logfile
// with w is the io.Writer , and level is minimum level of log for that particular writer
type MultiOption struct {
	W io.Writer
	// func(lvl LogLevel) bool
	Level LevelEnablerFunc
}

// This option also add rotate functionality to Multi File logging
type MultiOptionWithRotate struct {
	Filename string
	Ropt     RotateOptions
	// func(lvl LogLevel) bool
	Level LevelEnablerFunc
}

var (
	// option configure the logger to anotate each message with calling info or not
	WithCaller = zap.WithCaller
	// option configure to record stack trace message or not
	AddStacktrack = zap.AddStacktrace
	// add fileds to the logger
	AddFields = zap.Fields
)

const (
	//default level priority=0
	InfoLevel Loglevel = zap.InfoLevel
	// Dont need human review priority=1
	WarnLevel Loglevel = zap.WarnLevel
	// Indicate the application not run smoothly priority=2
	ErrorLevel Loglevel = zap.ErrorLevel
	//development panic , priority=3
	DPanicLevel Loglevel = zap.DPanicLevel
	// Log message then panic , priority=4
	PanicLevel Loglevel = zap.PanicLevel
	// Log message then exit with 1 , priority=5
	FatalLevel Loglevel = zap.FatalLevel
	// Usually not used in production , priority=-1
	DebugLevel Loglevel = zap.DebugLevel
)

// Used to add key value pair to the logger context
type LogField = zapcore.Field

type Logger struct {
	l     *zap.Logger
	level Loglevel
}

func (l *Logger) Debug(msg string, fields ...LogField) {
	l.l.Debug(msg, fields...)
}

func (l *Logger) Info(msg string, fields ...LogField) {
	l.l.Info(msg, fields...)
}

func (l *Logger) Warn(msg string, fields ...LogField) {
	l.l.Warn(msg, fields...)
}

func (l *Logger) Error(msg string, fields ...LogField) {
	l.l.Error(msg, fields...)
}

func (l *Logger) DPanic(msg string, fields ...LogField) {
	l.l.DPanic(msg, fields...)
}

func (l *Logger) Panic(msg string, fields ...LogField) {
	l.l.Panic(msg, fields...)
}

func (l *Logger) Fatal(msg string, fields ...LogField) {
	l.l.Fatal(msg, fields...)
}

// Flushed any buffered log entries.
// Should be called before exiting
func (l *Logger) Sync() error {
	return l.l.Sync()
}

// Return new logging instance with logging destination and set the minimum level of message to be record
func NewLogger(writer io.Writer, level Loglevel, opts ...LogOption) *Logger {
	if writer == nil {
		panic("the writer is nil")
	}
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format(time.RFC1123))
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(config.EncoderConfig),
		zapcore.AddSync(writer),
		zapcore.Level(level),
	)

	logger := &Logger{
		l:     zap.New(core, opts...),
		level: level,
	}
	return logger
}

var defaultLogger = NewLogger(os.Stderr, InfoLevel, WithCaller(true))

// Return the  default logger with error to stderr , and with info level
func DefaultLogger() *Logger {
	return defaultLogger
}

// Not safe for concurrent use
func ResetDefault(l *Logger) {
	defaultLogger = l
	Info = defaultLogger.Info
	Warn = defaultLogger.Warn
	Error = defaultLogger.Error
	DPanic = defaultLogger.DPanic
	Panic = defaultLogger.Panic
	Fatal = defaultLogger.Fatal
	Debug = defaultLogger.Debug
}

// Sync simplifies access to the logger's sync
func Sync() error {
	if defaultLogger != nil {
		return defaultLogger.Sync()
	}
	return nil
}

// This return logger that duplicate log entries to multiple log files
// Also base on the level field in the multi option , we can limit the minimum level of log to particular output
func NewMultiLoger(logs []MultiOption, opts ...LogOption) *Logger {
	var cores []zapcore.Core
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format(time.RFC1123))
	}
	for _, log := range logs {
		log := log
		if log.W == nil {
			panic("The writer is nil")
		}
		// this help to constraint the log  to particular output with minimum level
		lv := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return log.Level(lvl)
		})

		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(config.EncoderConfig),
			zapcore.AddSync(log.W),
			lv,
		)
		cores = append(cores, core)
	}
	logger := &Logger{
		l: zap.New(zapcore.NewTee(cores...), opts...),
	}
	return logger
}

func NewMultiWithRotate(logs []MultiOptionWithRotate, opts ...LogOption) *Logger {
	var cores []zapcore.Core

	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format(time.RFC1123))
	}

	for _, log := range logs {
		log := log
		lv := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return log.Level(Loglevel(lvl))

		})
		w := zapcore.AddSync(&lumberjack.Logger{
			Filename:   log.Filename,
			MaxSize:    log.Ropt.MaxSize,
			MaxBackups: log.Ropt.MaxBackups,
			MaxAge:     log.Ropt.MaxAge,
			Compress:   log.Ropt.Compress,
		})
		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(config.EncoderConfig),
			zapcore.AddSync(w),
			lv,
		)
		cores = append(cores, core)
	}
	logger := &Logger{
		l: zap.New(zapcore.NewTee(cores...), opts...),
	}
	return logger
}
