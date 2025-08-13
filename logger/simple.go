package logger

import (
	"fmt"
	"io"
	"os"
	"time"
)

type SimpleLogger struct {
	Enable bool
	Level  Level
	Out    io.Writer
}

// NewSimpleLogger creates a SimpleLogger with default level WARN and enabled.
func NewSimpleLogger() *SimpleLogger {
	return &SimpleLogger{
		Enable: true,
		Level:  LevelInfo,
		Out:    os.Stdout,
	}
}

func (s *SimpleLogger) Log(level Level, context, msg string, detail any) {
	if !s.Enable || level > s.Level {
		return
	}

	timestamp := time.Now().Format(time.RFC3339)
	fmt.Fprintf(s.Out, "[%s] %s | %s - %s\n", level.String(), timestamp, context, msg)
	if detail != nil {
		fmt.Fprintf(s.Out, "Detail: %+v\n", detail)
	}
}

func (s *SimpleLogger) Error(context, msg string, detail any) {
	s.Log(LevelError, context, msg, detail)
}

func (s *SimpleLogger) Warn(context, msg string, detail any) { s.Log(LevelWarn, context, msg, detail) }
func (s *SimpleLogger) Info(context, msg string, detail any) { s.Log(LevelInfo, context, msg, detail) }
func (s *SimpleLogger) Debug(context, msg string, detail any) {
	s.Log(LevelDebug, context, msg, detail)
}

func (s *SimpleLogger) Trace(context, msg string, detail any) {
	s.Log(LevelTrace, context, msg, detail)
}

func (s *SimpleLogger) SetEnabled(enabled bool) { s.Enable = enabled }
func (s *SimpleLogger) SetLevel(level Level)    { s.Level = level }
func (s *SimpleLogger) GetLevel() Level         { return s.Level }
func (s *SimpleLogger) SetOutput(w io.Writer)   { s.Out = w } // NEW
