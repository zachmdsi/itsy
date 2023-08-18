package itsy

import (
	"bufio"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// setupLogger sets up the logger for the Itsy instance.
func setupLogger() *zap.Logger {
	// Encoder Configuration
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.EpochTimeEncoder // Optimized time encoding
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	logWriter := zapcore.AddSync(bufio.NewWriter(os.Stdout))

	core := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), logWriter), zapcore.DebugLevel)
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	return logger
}

