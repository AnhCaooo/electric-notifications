// Created by Anh Cao on 27.08.2024.

package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

// todo: add defer to sync logger
// initialize logger
func Init() {
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // Use ISO8601 time format

	// Disable JSON encoding
	cfg.Encoding = "console"

	l, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	Logger = l
}
