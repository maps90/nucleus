package log

import (
	"context"

	"github.com/sirupsen/logrus"
)

// Entry struct
type Entry struct {
	*logrus.Entry
	ctx context.Context
}

// Report this should be used to report to sentry
func (l *Entry) Report(ctx context.Context) *Entry {
	l.ctx = ctx
	return l
}

// Error log
func Error(args ...interface{}) {
	New().Error(args...)
}

// Trace log
func Trace(args ...interface{}) {
	New().Trace(args...)
}

// Debug log
func Debug(args ...interface{}) {
	New().Debug(args...)
}

// Info log
func Info(args ...interface{}) {
	New().Info(args...)
}

// Warn log
func Warn(args ...interface{}) {
	New().Warn(args...)
}

// Fatal log
func Fatal(args ...interface{}) {
	New().Fatal(args...)
}

// Panic log
func Panic(args ...interface{}) {
	New().Panic(args...)
}
