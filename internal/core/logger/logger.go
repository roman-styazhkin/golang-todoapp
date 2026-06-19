package core_logger

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.Logger
	file *os.File
}

type loggerContextKey = struct{}

var (
	key = loggerContextKey{}
)

func ToContext(ctx context.Context, logger *Logger) context.Context {
	return context.WithValue(ctx, key, logger)
}

func NewLogger(config Config) (*Logger, error) {
	zapLevel := zap.NewAtomicLevel()

	if err := zapLevel.UnmarshalText([]byte(config.Level)); err != nil {
		return nil, fmt.Errorf("failed to unmarshal text, %w", err)
	}

	if err := os.MkdirAll(config.Folder, 0755); err != nil {
		return nil, fmt.Errorf("failed to make folder, %w", err)
	}

	timeNow := time.Now().UTC().Format("2006-01-02T15-04-05.000000")
	filePath := filepath.Join(
		config.Folder,
		fmt.Sprintf("%s.log", timeNow),
	)

	logFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open file, %w", err)
	}

	zapConfig := zap.NewDevelopmentEncoderConfig()
	zapConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02T15:04:05.000000")

	zapEncoder := zapcore.NewConsoleEncoder(zapConfig)

	core := zapcore.NewTee(
		zapcore.NewCore(zapEncoder, zapcore.AddSync(logFile), zapLevel),
		zapcore.NewCore(zapEncoder, zapcore.AddSync(os.Stdout), zapLevel),
	)

	logger := zap.New(core, zap.AddCaller())
	return &Logger{
		Logger: logger,
		file:   logFile,
	}, nil
}

func (l *Logger) Close() {
	if err := l.file.Close(); err != nil {
		panic(err)
	}
}

func (l *Logger) With(fields ...zap.Field) *Logger {
	return &Logger{
		Logger: l.Logger.With(fields...),
		file:   l.file,
	}
}

func FromContext(ctx context.Context) *Logger {
	logger, ok := ctx.Value(key).(*Logger)

	if !ok {
		panic("failed to get logger from context")
	}

	return logger
}
