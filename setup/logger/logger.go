package logger

import (
	"time"

	nwConfig "github.com/gennesseaux/NotionWatcher/setup/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger : Global variable to store the logger
var Logger *zap.Logger

// config : config
var config = nwConfig.Config

func init() {
	if config.Environment == "development" {
		zapConfig := zap.NewDevelopmentConfig()
		zapConfig.Level, _ = zap.ParseAtomicLevel(config.LogLevel)
		zapConfig.DisableCaller = false
		zapConfig.Encoding = "console"
		zapConfig.EncoderConfig.TimeKey = "timestamp"
		zapConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)
		zapLogger, _ := zapConfig.Build()
		Logger = zapLogger
	} else if config.Environment == "production" {
		zapConfig := zap.NewProductionConfig()
		zapConfig.Level, _ = zap.ParseAtomicLevel(config.LogLevel)
		zapConfig.DisableCaller = true
		zapConfig.Encoding = "console"
		zapConfig.EncoderConfig.TimeKey = "timestamp"
		zapConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)
		zapLogger, _ := zapConfig.Build()
		Logger = zapLogger
	}
}
