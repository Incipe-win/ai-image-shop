package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger

func Init() {
	var err error
	config := zap.NewProductionConfig()
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	log, err = config.Build(zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}
}

func Sync() {
	log.Sync()
}

func Info(message string, fields ...interface{}) {
	log.Info(message, convertToZapFields(fields)...)
}

func Error(message string, fields ...interface{}) {
	log.Error(message, convertToZapFields(fields)...)
}

func Fatal(message string, fields ...interface{}) {
	log.Fatal(message, convertToZapFields(fields)...)
}

func Debug(message string, fields ...interface{}) {
	log.Debug(message, convertToZapFields(fields)...)
}

func convertToZapFields(fields []interface{}) []zap.Field {
	if len(fields)%2 != 0 {
		return []zap.Field{zap.Any("extra", fields)}
	}

	zapFields := make([]zap.Field, 0, len(fields)/2)
	for i := 0; i < len(fields); i += 2 {
		key, ok := fields[i].(string)
		if !ok {
			continue
		}
		zapFields = append(zapFields, zap.Any(key, fields[i+1]))
	}
	return zapFields
}