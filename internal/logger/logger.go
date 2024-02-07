package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger = zap.SugaredLogger

func New() *Logger {
	var logger *zap.Logger

	if os.Getenv("ENVIRONMENT") == "production" {
		logger, _ = zap.NewProduction()
	} else {
		config := zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		logger, _ = config.Build()
	}

	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()

	return sugar
}
