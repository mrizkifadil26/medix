package logger

import (
	"fmt"
	"io"
	"time"

	"github.com/sirupsen/logrus"
)

// LogrusLogger wraps a logrus.Logger to implement Logger interface.
type LogrusLogger struct {
	enabled bool
	level   Level
	context string
	logger  *logrus.Logger
}

func NewLogrusLogger() *LogrusLogger {
	baseLogger := logrus.New()
	baseLogger.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: time.RFC3339,
	})

	l := &LogrusLogger{
		enabled: true,
		level:   LevelInfo,
		logger:  baseLogger,
	}

	l.syncLogrusLevel()

	return l
}

func (l *LogrusLogger) WithContext(ctx string) Logger {
	return &LogrusLogger{
		enabled: l.enabled,
		level:   l.level,
		logger:  l.logger,
		context: ctx,
	}
}

func (l *LogrusLogger) SetEnabled(enabled bool) { l.enabled = enabled }
func (l *LogrusLogger) GetLevel() Level         { return l.level }
func (l *LogrusLogger) SetLevel(level Level) {
	l.level = level
	l.syncLogrusLevel()
}

func (l *LogrusLogger) SetOutput(w io.Writer) { l.logger.SetOutput(w) }

func (l *LogrusLogger) Log(level Level, msg string, detail ...any) {
	if !l.enabled || level > l.level {
		return
	}

	entry := l.logger.WithFields(logrus.Fields{})
	if l.context != "" {
		entry = entry.WithField("context", l.context)
	}

	if len(detail) > 0 && detail[0] != nil {
		// Try to flatten if single map[string]interface{}
		if len(detail) == 1 {
			if fields, ok := detail[0].(map[string]interface{}); ok {
				entry = entry.WithFields(fields)
			} else {
				// fallback: join all details as string
				entry = entry.WithField("detail", fmt.Sprint(detail...))
			}
		} else {
			entry = entry.WithField("detail", fmt.Sprint(detail...))
		}
	}

	switch level {
	case LevelError:
		entry.Error(msg)
	case LevelWarn:
		entry.Warn(msg)
	case LevelInfo:
		entry.Info(msg)
	case LevelDebug:
		entry.Debug(msg)
	case LevelTrace:
		entry.Trace(msg)
	default:
		entry.Info(msg)
	}
}

func (l *LogrusLogger) Error(msg string, detail ...any) {
	l.Log(LevelError, msg, detail...)
}

func (l *LogrusLogger) Warn(msg string, detail ...any) {
	l.Log(LevelWarn, msg, detail...)
}

func (l *LogrusLogger) Info(msg string, detail ...any) {
	l.Log(LevelInfo, msg, detail...)
}

func (l *LogrusLogger) Debug(msg string, detail ...any) {
	l.Log(LevelDebug, msg, detail...)
}

func (l *LogrusLogger) Trace(msg string, detail ...any) {
	l.Log(LevelTrace, msg, detail...)
}

// syncLogrusLevel adjusts the internal logrus.Logger level to match l.level
func (l *LogrusLogger) syncLogrusLevel() {
	switch l.level {
	case LevelError:
		l.logger.SetLevel(logrus.ErrorLevel)
	case LevelWarn:
		l.logger.SetLevel(logrus.WarnLevel)
	case LevelInfo:
		l.logger.SetLevel(logrus.InfoLevel)
	case LevelDebug:
		l.logger.SetLevel(logrus.DebugLevel)
	case LevelTrace:
		l.logger.SetLevel(logrus.TraceLevel)
	default:
		l.logger.SetLevel(logrus.InfoLevel)
	}
}
