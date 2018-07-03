package log

import (
	"context"

	"github.com/sirupsen/logrus"
)

// Logger is a Logrus logger used by the package.
var Logger *logrus.Logger

// F is a short alias for a string to interface map.
type F map[string]interface{}

func init() {
	Logger = logrus.New()
}

// Debug prints a debug message with given fields attached.
func Debug(ctx context.Context, msg string, fields ...F) {
	withFields(mergeFields(ctx, fields)).Debug(msg)
}

// Debugf prints a formatted debug message.
func Debugf(ctx context.Context, format string, args ...interface{}) {
	withFields(contextFields(ctx)).Debugf(format, args...)
}

// Info prints an info message with given fields attached.
func Info(ctx context.Context, msg string, fields ...F) {
	withFields(mergeFields(ctx, fields)).Info(msg)
}

// Infof prints a formatted info message.
func Infof(ctx context.Context, format string, args ...interface{}) {
	withFields(contextFields(ctx)).Infof(format, args...)
}

// Error prints an error message with given fields attached.
func Error(ctx context.Context, msg string, fields ...F) {
	withFields(mergeFields(ctx, fields)).Error(msg)
}

// Errorf prints a formatted error message.
func Errorf(ctx context.Context, format string, args ...interface{}) {
	withFields(contextFields(ctx)).Errorf(format, args...)
}

// Fatal prints an error message with given fields attached and then exits.
func Fatal(ctx context.Context, msg string, fields ...F) {
	withFields(mergeFields(ctx, fields)).Fatal(msg)
}

// Fatalf prints a formatted error message and then exits.
func Fatalf(ctx context.Context, format string, args ...interface{}) {
	withFields(contextFields(ctx)).Fatalf(format, args...)
}

func mergeFields(ctx context.Context, fields []F) F {
	ctxf := contextFields(ctx)
	if len(ctxf) == 0 && len(fields) == 0 {
		return nil
	}
	for i := 0; i < len(fields); i++ {
		for k, v := range fields[i] {
			ctxf[k] = v
		}
	}
	return ctxf
}

func withFields(f F) *logrus.Entry {
	if len(f) == 0 {
		return logrus.NewEntry(Logger)
	}
	return Logger.WithFields(logrus.Fields(f))
}
