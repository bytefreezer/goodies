package log

import (
	"fmt"
	"log/slog"
	"os"
)

var defaultHandlerOptions = slog.HandlerOptions{AddSource: true, ReplaceAttr: resolveLogLevel}

var defaultLogger = &Logger{logger: slog.New(slog.NewJSONHandler(os.Stdout, &defaultHandlerOptions)), w: os.Stdout}

func SetMinLogLevel(lev slog.Leveler) {
	cpo := defaultHandlerOptions
	cpo.Level = lev
	newSlog := slog.New(slog.NewJSONHandler(os.Stdout, &cpo))
	defaultLogger.Lock()
	defaultLogger.logger = newSlog
	defaultLogger.Unlock()
}

func WithTrace(traceID string) *Logger {
	return defaultLogger.With(TraceKey, traceID)
}

func WithCustomer(customerID string) *Logger {
	return defaultLogger.With(CustomerKey, customerID)
}

func With(args ...any) *Logger {
	defaultLogger.Lock()
	defer defaultLogger.Unlock()
	nl := Logger{w: defaultLogger.w}
	nl.logger = defaultLogger.logger.With(args...)
	return &nl
}

func Debug(msg string) {
	defaultLogger.Lock()
	defer defaultLogger.Unlock()
	defaultLogger.log(slog.LevelDebug, msg)
}

func Debugf(msg string, args ...any) {
	defaultLogger.Lock()
	defer defaultLogger.Unlock()
	msg = fmt.Sprintf(msg, args...)
	defaultLogger.log(slog.LevelDebug, msg)
}

func Info(msg string) {
	defaultLogger.Lock()
	defer defaultLogger.Unlock()
	defaultLogger.log(slog.LevelInfo, msg)
}

func Infof(msg string, args ...any) {
	defaultLogger.Lock()
	defer defaultLogger.Unlock()
	msg = fmt.Sprintf(msg, args...)
	defaultLogger.log(slog.LevelInfo, msg)
}

func Warn(msg string) {
	defaultLogger.Lock()
	defer defaultLogger.Unlock()
	defaultLogger.log(slog.LevelWarn, msg)
}

func Warnf(msg string, args ...any) {
	defaultLogger.Lock()
	defer defaultLogger.Unlock()
	msg = fmt.Sprintf(msg, args...)
	defaultLogger.log(slog.LevelWarn, msg)
}

func Error(msg string) {
	defaultLogger.Lock()
	defer defaultLogger.Unlock()
	defaultLogger.log(slog.LevelError, msg)
}

func Errorf(msg string, args ...any) {
	defaultLogger.Lock()
	defer defaultLogger.Unlock()
	msg = fmt.Sprintf(msg, args...)
	defaultLogger.log(slog.LevelError, msg)
}

func Fatal(msg string) {
	defaultLogger.Lock()
	defer defaultLogger.Unlock()
	defaultLogger.log(LevelFatal, msg)
	os.Exit(1)
}

func Fatalf(msg string, args ...any) {
	defaultLogger.Lock()
	defer defaultLogger.Unlock()
	msg = fmt.Sprintf(msg, args...)
	defaultLogger.log(LevelFatal, msg)
	os.Exit(1)
}

func Print(msg string) {
	defaultLogger.w.Write([]byte(msg + "\n"))
}

func Printf(msg string, args ...any) {
	msg = fmt.Sprintf(msg, args...)
	defaultLogger.w.Write([]byte(msg + "\n"))
}
