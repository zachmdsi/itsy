package itsy

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// SetupLogger creates a new logger instance.
func SetupLogger() *zap.Logger {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoder := zapcore.NewConsoleEncoder(encoderConfig)
	writeSyncer := zapcore.Lock(os.Stdout)
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)
	logger := zap.New(core)
	return logger
}