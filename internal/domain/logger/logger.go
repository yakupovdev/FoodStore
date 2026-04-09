package logger

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

func NewLogger(config Config) (*Logger, error) {
	zapLvl := zap.NewAtomicLevel()
	if err := zapLvl.UnmarshalText([]byte(config.LogLevel)); err != nil {
		return nil, fmt.Errorf("unmarshal log level: %w", err)
	}
	if err := os.MkdirAll(config.LogFolder, 0755); err != nil {
		return nil, fmt.Errorf("create log folder: %w", err)
	}
	timestamp := time.Now().UTC().Format("2006-01-02T15-04-05.00000")
	logFilePath := filepath.Join(
		config.LogFolder,
		fmt.Sprintf("%s.log", timestamp),
	)

	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("open log file: %w", err)
	}
	zapConfig := zap.NewDevelopmentEncoderConfig()
	zapConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02T15:04:05.000000")
	zapEncoder := zapcore.NewJSONEncoder(zapConfig)

	core := zapcore.NewTee(
		zapcore.NewCore(zapEncoder, zapcore.AddSync(logFile), zapLvl),
	)

	zapLogger := zap.New(core, zap.AddCaller())

	return &Logger{
		Logger: zapLogger,
		file:   logFile,
	}, nil
}

func (l *Logger) Close() {
	if err := l.file.Close(); err != nil {
		fmt.Printf("failed to close log file: %v\n", err)
	}
}

func (l *Logger) With(field ...zap.Field) *Logger {
	return &Logger{
		Logger: l.Logger.With(field...),
		file:   l.file,
	}
}

func FromContext(ctx context.Context) *Logger {
	log, ok := ctx.Value("log").(*Logger)
	if !ok {
		panic("no log in context")
	}
	return log
}
