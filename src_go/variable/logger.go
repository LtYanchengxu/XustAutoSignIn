package variable

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func getLoggerEncoder() zapcore.Encoder {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.TimeEncoder(func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.UTC().Local().Format(time.RFC1123Z))
	})
	config.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(config)
}

func NewConsoleLogger() *zap.SugaredLogger {
	writer := zapcore.AddSync(os.Stdout)
	core := zapcore.NewCore(getLoggerEncoder(), writer, zapcore.DebugLevel)
	return zap.New(core).Sugar()
}

func NewFileLogger(path string, level zapcore.Level) *zap.SugaredLogger {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   path,
		MaxSize:    512,
		MaxBackups: 12,
		MaxAge:     365,
		Compress:   true,
	}
	writer := zapcore.AddSync(lumberJackLogger)
	core := zapcore.NewCore(getLoggerEncoder(), writer, level)
	return zap.New(core).Sugar()
}
