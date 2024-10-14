// AnhCao 2024
package logger

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

// initialize logger
func Init() {
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	cfg.EncoderConfig.EncodeTime = syslogTimeEncoder

	// Disable JSON encoding
	cfg.Encoding = "console"

	l, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	Logger = l
}

func syslogTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}
