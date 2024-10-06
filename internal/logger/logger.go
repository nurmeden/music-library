package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var _ Logger = (*loggerWrapper)(nil) // Подтверждаем, что zapLogger реализует интерфейс Logger

type loggerWrapper struct {
	sugaredLogger *zap.SugaredLogger
}

func NewZapLogger() Logger {
	writerSyncer := getLogWriter()
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writerSyncer, zapcore.DebugLevel)
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zap.FatalLevel))
	sugaredLogger := logger.Sugar()

	return &loggerWrapper{sugaredLogger: sugaredLogger}
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter() zapcore.WriteSyncer {
	return zapcore.AddSync(os.Stdout)
}

// Реализация методов интерфейса Logger

func (l *loggerWrapper) Debug(args ...interface{}) {
	l.sugaredLogger.Debug(args...)
}

func (l *loggerWrapper) Debugf(template string, args ...interface{}) {
	l.sugaredLogger.Debugf(template, args...)
}

func (l *loggerWrapper) Info(args ...interface{}) {
	l.sugaredLogger.Info(args...)
}

func (l *loggerWrapper) Infof(template string, args ...interface{}) {
	l.sugaredLogger.Infof(template, args...)
}

func (l *loggerWrapper) Warn(args ...interface{}) {
	l.sugaredLogger.Warn(args...)
}

func (l *loggerWrapper) Warnf(template string, args ...interface{}) {
	l.sugaredLogger.Warnf(template, args...)
}

func (l *loggerWrapper) Error(args ...interface{}) {
	l.sugaredLogger.Error(args...)
}

func (l *loggerWrapper) Errorf(template string, args ...interface{}) {
	l.sugaredLogger.Errorf(template, args...)
}

func (l *loggerWrapper) Panic(args ...interface{}) {
	l.sugaredLogger.Panic(args...)
}

func (l *loggerWrapper) Panicf(template string, args ...interface{}) {
	l.sugaredLogger.Panicf(template, args...)
}

func (l *loggerWrapper) Fatal(args ...interface{}) {
	l.sugaredLogger.Fatal(args...)
}

func (l *loggerWrapper) Fatalf(template string, args ...interface{}) {
	l.sugaredLogger.Fatalf(template, args...)
}
