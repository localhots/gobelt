package log

import (
	"github.com/lytics/logrus"
)

// Logger is a Logrus logger used by the package.
var Logger *logrus.Logger

// F is a short alias for a string to interface map.
type F map[string]interface{}

func init() {
	Logger = logrus.New()
}

// Debug prints a debug message with given fields attached.
func Debug(msg string, fields ...F) {
	withFields(fields).Debug(msg)
}

// Debugf prints a formatted debug message.
func Debugf(format string, args ...interface{}) {
	logrus.Debugf(format, args...)
}

// Info prints an info message with given fields attached.
func Info(msg string, fields ...F) {
	withFields(fields).Info(msg)
}

// Infof prints a formatted info message.
func Infof(format string, args ...interface{}) {
	logrus.Infof(format, args...)
}

// Error prints an error message with given fields attached.
func Error(msg string, fields ...F) {
	withFields(fields).Error(msg)
}

// Errorf prints a formatted error message.
func Errorf(format string, args ...interface{}) {
	logrus.Errorf(format, args...)
}

// Fatal prints an error message with given fields attached and then exits.
func Fatal(msg string, fields ...F) {
	withFields(fields).Fatal(msg)
}

// Fatalf prints a formatted error message and then exits.
func Fatalf(format string, args ...interface{}) {
	logrus.Fatalf(format, args...)
}

func withFields(fields []F) *logrus.Entry {
	switch len(fields) {
	case 0:
		return logrus.NewEntry(Logger)
	case 1:
		return logrus.WithFields(logrus.Fields(fields[0]))
	default:
		f := F{}
		for i := 0; i < len(fields); i++ {
			for k, v := range fields[i] {
				f[k] = v
			}
		}
		return logrus.WithFields(logrus.Fields(f))
	}
}
