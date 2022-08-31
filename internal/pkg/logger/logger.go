package logger

import "go.uber.org/zap"

var globalLogger *zap.SugaredLogger

func init() {
	globalLogger = zap.NewExample().Sugar()
}

func GetLogger() *zap.SugaredLogger {
	return globalLogger
}

func SetLogger(logger *zap.SugaredLogger) {
	globalLogger = logger
}

func Debug(args ...interface{}) {
	globalLogger.Debug(args...)
}

func Info(args ...interface{}) {
	globalLogger.Info(args...)
}

func Warn(args ...interface{}) {
	globalLogger.Warn(args...)
}

func Error(args ...interface{}) {
	globalLogger.Error(args...)
}

func Fatal(args ...interface{}) {
	globalLogger.Fatal(args...)
}

func Debugf(template string, args ...interface{}) {
	globalLogger.Debugf(template, args...)
}

func Infof(template string, args ...interface{}) {
	globalLogger.Infof(template, args...)
}

func Warnf(template string, args ...interface{}) {
	globalLogger.Warnf(template, args...)
}

func Errorf(template string, args ...interface{}) {
	globalLogger.Errorf(template, args...)
}

func Fatalf(template string, args ...interface{}) {
	globalLogger.Fatalf(template, args...)
}
