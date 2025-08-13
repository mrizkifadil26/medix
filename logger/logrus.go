package logger

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// LogrusLogger wraps a logrus.Logger to implement Logger interface.
type LogrusLogger struct {
	logger *logrus.Logger
	Enable bool
	Level  Level
}

func NewLogrusLogger() *LogrusLogger {
	baseLogger := logrus.New()
	baseLogger.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: time.RFC3339,
	})

	l := &LogrusLogger{
		logger: baseLogger,
		Enable: true,
		Level:  LevelInfo,
	}

	l.syncLogrusLevel()

	return l
}

func (l *LogrusLogger) Log(level Level, context, msg string, detail any) {
	if !l.Enable || level > l.Level {
		return
	}

	col := colorForLevel(level)
	tag := fmt.Sprintf("%s[%s]%s", col, strings.ToUpper(level.String()), colorReset)
	fullMsg := fmt.Sprintf("%s %s - %s", tag, context, msg)

	entry := l.logger.WithField("context", context)
	if detail != nil {
		entry = entry.WithField("detail", detail)
	}

	switch level {
	case LevelError:
		entry.Error(fullMsg)
	case LevelWarn:
		entry.Warn(fullMsg)
	case LevelInfo:
		entry.Info(fullMsg)
	case LevelDebug:
		entry.Debug(fullMsg)
	case LevelTrace:
		entry.Trace(fullMsg)
	default:
		entry.Info(fullMsg)
	}
}

func (l *LogrusLogger) Error(context, msg string, detail any) {
	l.Log(LevelError, context, msg, detail)
}
func (l *LogrusLogger) Warn(context, msg string, detail any) { l.Log(LevelWarn, context, msg, detail) }
func (l *LogrusLogger) Info(context, msg string, detail any) { l.Log(LevelInfo, context, msg, detail) }
func (l *LogrusLogger) Debug(context, msg string, detail any) {
	l.Log(LevelDebug, context, msg, detail)
}
func (l *LogrusLogger) Trace(context, msg string, detail any) {
	l.Log(LevelTrace, context, msg, detail)
}

func (l *LogrusLogger) SetEnabled(enabled bool) { l.Enable = enabled }
func (l *LogrusLogger) SetLevel(level Level) {
	l.Level = level
	l.syncLogrusLevel()
}
func (l *LogrusLogger) GetLevel() Level { return l.Level }
func (l *LogrusLogger) SetOutput(w io.Writer) {
	l.logger.SetOutput(w)
}

// ANSI color codes for LogrusLogger
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorCyan   = "\033[36m"
	colorGray   = "\033[37m"
)

func colorForLevel(l Level) string {
	switch l {
	case LevelError:
		return colorRed
	case LevelWarn:
		return colorYellow
	case LevelInfo:
		return colorBlue
	case LevelDebug:
		return colorCyan
	case LevelTrace:
		return colorGray
	default:
		return colorReset
	}
}

// syncLogrusLevel adjusts the internal logrus.Logger level to match l.level
func (l *LogrusLogger) syncLogrusLevel() {
	switch l.Level {
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
