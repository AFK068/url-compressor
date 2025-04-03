package logger

import (
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New() *zap.Logger {
	zapConfig := zap.NewProductionConfig()

	// Setup config.
	zapConfig.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	zapConfig.DisableStacktrace = true
	zapConfig.DisableCaller = true

	logger, err := zapConfig.Build()
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}

	return logger
}
