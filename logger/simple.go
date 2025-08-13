package logger

import (
	"fmt"
	"io"
	"os"
	"time"
)

type SimpleLogger struct {
	enabled bool
	level   Level
	context string
	output  io.Writer
}

// NewSimpleLogger creates a SimpleLogger with default level INFO and enabled.
func NewSimpleLogger() *SimpleLogger {
	return &SimpleLogger{
		enabled: true,
		level:   LevelInfo,
		output:  os.Stdout,
	}
}

// WithContext returns a new logger with the given default context bound.
func (l *SimpleLogger) WithContext(ctx string) Logger {
	return &SimpleLogger{
		enabled: l.enabled,
		level:   l.level,
		output:  l.output,
		context: ctx,
	}
}

func (l *SimpleLogger) SetEnabled(enabled bool) { l.enabled = enabled }
func (l *SimpleLogger) SetLevel(level Level)    { l.level = level }
func (l *SimpleLogger) GetLevel() Level         { return l.level }
func (l *SimpleLogger) SetOutput(w io.Writer)   { l.output = w } // NEW

// Log writes a log entry.
// If ctxOverride is not empty, it will override the default context.
func (l *SimpleLogger) Log(level Level, msg string, detail ...any) {
	if !l.enabled || level > l.level {
		return
	}

	timestamp := time.Now().Format(time.RFC3339)
	output := fmt.Sprintf("[%s] [%s]", timestamp, level.String())
	if l.context != "" {
		output += fmt.Sprintf(" [%s]", l.context)
	}

	output += " " + msg

	if len(detail) > 0 {
		for _, d := range detail {
			switch v := d.(type) {
			case map[string]interface{}:
				// Print each key=value pair from map
				for key, val := range v {
					output += fmt.Sprintf(" %s=%v", key, val)
				}
			default:
				// fallback print detail normally
				output += fmt.Sprintf(" %v", v)
			}
		}
	}

	fmt.Fprintln(l.output, output)
}

func (l *SimpleLogger) Error(msg string, detail ...any) {
	l.Log(LevelError, msg, detail...)
}

func (l *SimpleLogger) Warn(msg string, detail ...any) {
	l.Log(LevelWarn, msg, detail...)
}

func (l *SimpleLogger) Info(msg string, detail ...any) {
	l.Log(LevelInfo, msg, detail...)
}

func (l *SimpleLogger) Debug(msg string, detail ...any) {
	l.Log(LevelDebug, msg, detail...)
}

func (l *SimpleLogger) Trace(msg string, detail ...any) {
	l.Log(LevelTrace, msg, detail...)
}
