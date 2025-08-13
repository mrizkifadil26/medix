package logger

import (
	"io"
)

// defaultLogger holds the global logger instance
var defaultLogger Logger = NewSimpleLogger()

// InitSimple sets the default logger to SimpleLogger
func InitSimple(enabled bool, level Level, out io.Writer) {
	l := NewSimpleLogger()
	l.SetEnabled(enabled)
	l.SetLevel(level)
	if out != nil {
		l.SetOutput(out)
	}

	defaultLogger = l
}

// InitLogrus sets the default logger to LogrusLogger
func InitLogrus(enabled bool, level Level, out io.Writer) {
	l := NewLogrusLogger()
	l.SetEnabled(enabled)
	l.SetLevel(level)
	if out != nil {
		l.SetOutput(out)
	}

	defaultLogger = l
}

// SetLogger allows injecting any custom Logger
func SetLogger(l Logger) {
	if l != nil {
		defaultLogger = l
	}
}

// --- Wrapper functions for easy calls ---
func Error(msg string, detail ...any) { defaultLogger.Error(msg, detail...) }
func Warn(msg string, detail ...any)  { defaultLogger.Warn(msg, detail...) }
func Info(msg string, detail ...any)  { defaultLogger.Info(msg, detail...) }
func Debug(msg string, detail ...any) { defaultLogger.Debug(msg, detail...) }
func Trace(msg string, detail ...any) { defaultLogger.Trace(msg, detail...) }
