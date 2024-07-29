package loggo

import (
	"context"
	"log"
	"log/slog"
	"os"

	slogmulti "github.com/samber/slog-multi"
)

type LoggerInterface interface {
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, err error, args ...interface{})
	Fatal(msg string, err error, args ...interface{})
	WithOperation(operationID string) LoggerInterface
	IsDebugEnabled() bool
}

type Logger struct {
	logger      *slog.Logger
	operationID string
	level       slog.Level
}

var _ LoggerInterface = &Logger{}

func NewLogger(logFilePath string, level slog.Level) (LoggerInterface, error) {
	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("Error opening log file: %v", err)
		return nil, err
	}

	fileHandler := slog.NewJSONHandler(file, &slog.HandlerOptions{
		Level: level,
	})

	stdoutHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo, // Set stdout to INFO level
	})

	multiHandler := slogmulti.Fanout(fileHandler, stdoutHandler)

	slogLogger := slog.New(multiHandler)

	return &Logger{logger: slogLogger, level: level}, nil
}

func (l *Logger) log(level slog.Level, msg string, args ...interface{}) {
	if l.operationID != "" {
		args = append(args, "operationID", l.operationID)
	}
	l.logger.Log(context.TODO(), level, msg, args...)
}

func (l *Logger) Debug(msg string, args ...interface{}) {
	l.log(slog.LevelDebug, msg, args...)
}

func (l *Logger) Info(msg string, args ...interface{}) {
	l.log(slog.LevelInfo, msg, args...)
}

func (l *Logger) Warn(msg string, args ...interface{}) {
	l.log(slog.LevelWarn, msg, args...)
}

func (l *Logger) Error(msg string, err error, args ...interface{}) {
	args = append(args, "error", err)
	l.log(slog.LevelError, msg, args...)
}

func (l *Logger) Fatal(msg string, err error, args ...interface{}) { // Implement Fatal method
	args = append(args, "error", err)
	l.logger.Log(context.TODO(), slog.LevelError, msg, args...)
	os.Exit(1)
}

func (l *Logger) WithOperation(operationID string) LoggerInterface {
	return &Logger{
		logger:      l.logger,
		operationID: operationID,
		level:       l.level,
	}
}

func (l *Logger) IsDebugEnabled() bool {
	return l.level <= slog.LevelDebug
}
