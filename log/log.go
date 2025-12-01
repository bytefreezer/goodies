package log

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"runtime"
	"sync"
	"time"
)

const (
	TraceKey    = "traceId"
	CustomerKey = "customerId"
)

type Record struct {
	Time       time.Time    `json:"time"`
	Level      string       `json:"level"`
	Source     *slog.Source `json:"source"`
	Message    string       `json:"msg"`
	TraceID    string       `json:"traceId"`
	CustomerID string       `json:"customerId"`
}

type Logger struct {
	sync.Mutex
	logger *slog.Logger
	w      io.Writer
}

func New(w io.Writer) *Logger {
	return &Logger{logger: slog.New(slog.NewJSONHandler(w, &defaultHandlerOptions)), w: w}
}

func (l *Logger) WithTrace(traceID string) *Logger {
	return l.With(TraceKey, traceID)
}

func (l *Logger) WithCustomer(customerID string) *Logger {
	return l.With(CustomerKey, customerID)
}

func (l *Logger) With(args ...any) *Logger {
	l.Lock()
	defer l.Unlock()
	nl := Logger{w: l.w}
	nl.logger = l.logger.With(args...)
	return &nl
}

func (l *Logger) Debug(msg string) {
	l.Lock()
	defer l.Unlock()
	l.log(slog.LevelDebug, msg)
}

func (l *Logger) Debugf(msg string, args ...any) {
	l.Lock()
	defer l.Unlock()
	msg = fmt.Sprintf(msg, args...)
	l.log(slog.LevelDebug, msg)
}

func (l *Logger) Info(msg string) {
	l.Lock()
	defer l.Unlock()
	l.log(slog.LevelInfo, msg)
}

func (l *Logger) Infof(msg string, args ...any) {
	l.Lock()
	defer l.Unlock()
	msg = fmt.Sprintf(msg, args...)
	l.log(slog.LevelInfo, msg)
}

func (l *Logger) Warn(msg string) {
	l.Lock()
	defer l.Unlock()
	l.log(slog.LevelWarn, msg)
}

func (l *Logger) Warnf(msg string, args ...any) {
	l.Lock()
	defer l.Unlock()
	msg = fmt.Sprintf(msg, args...)
	l.log(slog.LevelWarn, msg)
}

func (l *Logger) Error(msg string) {
	l.Lock()
	defer l.Unlock()
	l.log(slog.LevelError, msg)
}

func (l *Logger) Errorf(msg string, args ...any) {
	l.Lock()
	defer l.Unlock()
	msg = fmt.Sprintf(msg, args...)
	l.log(slog.LevelError, msg)
}

func (l *Logger) Fatal(msg string) {
	l.Lock()
	defer l.Unlock()
	l.log(LevelFatal, msg)
	os.Exit(1)
}

func (l *Logger) Fatalf(msg string, args ...any) {
	l.Lock()
	defer l.Unlock()
	msg = fmt.Sprintf(msg, args...)
	l.log(LevelFatal, msg)
	os.Exit(1)
}

func (l *Logger) Print(msg string) {
	l.w.Write([]byte(msg + "\n"))
}

func (l *Logger) Printf(msg string, args ...any) {
	msg = fmt.Sprintf(msg, args...)
	l.w.Write([]byte(msg + "\n"))
}

func (l *Logger) SetMinLogLevel(lev slog.Leveler) {
	cpo := defaultHandlerOptions
	cpo.Level = lev
	newSlog := slog.New(slog.NewJSONHandler(l.w, &cpo))
	l.Lock()
	l.logger = newSlog
	l.Unlock()
}

func (l *Logger) log(level slog.Level, msg string, args ...any) {
	ctx := context.Background()
	if !l.logger.Enabled(ctx, level) {
		return
	}
	var pc uintptr
	var pcs [1]uintptr
	// skip [runtime.Callers, this function, this function's caller]
	runtime.Callers(3, pcs[:])
	pc = pcs[0]

	r := slog.NewRecord(time.Now(), level, msg, pc)
	r.Add(args...)
	_ = l.logger.Handler().Handle(ctx, r)
}
