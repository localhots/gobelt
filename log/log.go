package log

import (
	"context"

	"github.com/sirupsen/logrus"
)

// F is a short alias for a string to interface map.
type F map[string]interface{}

// Level defines logging level.
type Level byte

const (
	// LevelDebug is a logging level of debug messages.
	LevelDebug Level = iota
	// LevelInfo is a logging level of info messages.
	LevelInfo
	// LevelWarn is a logging level of warning messages.
	LevelWarn
	// LevelError is a logging level of error messages.
	LevelError
	// LevelFatal is a logging level of fatal messages.
	LevelFatal
)

var (
	// Logger is a Logrus logger used by the package.
	Logger = logrus.New()

	level = LevelDebug
)

// SetLevel sets up minimum logging level.
func SetLevel(l Level) {
	level = l
}

// Debug prints a debug message with given fields attached.
func Debug(ctx context.Context, msg string, fields ...F) {
	if level >= LevelDebug {
		withFields(mergeFields(ctx, fields)).Debug(msg)
	}
}

// Debugf prints a formatted debug message.
func Debugf(ctx context.Context, format string, args ...interface{}) {
	if level >= LevelDebug {
		withFields(contextFields(ctx)).Debugf(format, args...)
	}
}

// Info prints an info message with given fields attached.
func Info(ctx context.Context, msg string, fields ...F) {
	if level >= LevelInfo {
		withFields(mergeFields(ctx, fields)).Info(msg)
	}
}

// Infof prints a formatted info message.
func Infof(ctx context.Context, format string, args ...interface{}) {
	if level >= LevelInfo {
		withFields(contextFields(ctx)).Infof(format, args...)
	}
}

// Warn prints an warning message with given fields attached.
func Warn(ctx context.Context, msg string, fields ...F) {
	if level >= LevelWarn {
		withFields(mergeFields(ctx, fields)).Warn(msg)
	}
}

// Warnf prints a formatted warning message.
func Warnf(ctx context.Context, format string, args ...interface{}) {
	if level >= LevelWarn {
		withFields(contextFields(ctx)).Warnf(format, args...)
	}
}

// Error prints an error message with given fields attached.
func Error(ctx context.Context, msg string, fields ...F) {
	if level >= LevelError {
		withFields(mergeFields(ctx, fields)).Error(msg)
	}
}

// Errorf prints a formatted error message.
func Errorf(ctx context.Context, format string, args ...interface{}) {
	if level >= LevelError {
		withFields(contextFields(ctx)).Errorf(format, args...)
	}
}

// Fatal prints an error message with given fields attached and then exits.
func Fatal(ctx context.Context, msg string, fields ...F) {
	if level >= LevelFatal {
		withFields(mergeFields(ctx, fields)).Fatal(msg)
	}
}

// Fatalf prints a formatted error message and then exits.
func Fatalf(ctx context.Context, format string, args ...interface{}) {
	if level >= LevelFatal {
		withFields(contextFields(ctx)).Fatalf(format, args...)
	}
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
