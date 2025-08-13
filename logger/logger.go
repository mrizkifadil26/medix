package logger

import "io"

// Level defines log severity levels.
type Level int

const (
	LevelError Level = iota + 1
	LevelWarn
	LevelInfo
	LevelDebug
	LevelTrace
)

func (l Level) String() string {
	switch l {
	case LevelError:
		return "ERROR"
	case LevelWarn:
		return "WARN"
	case LevelInfo:
		return "INFO"
	case LevelDebug:
		return "DEBUG"
	case LevelTrace:
		return "TRACE"
	default:
		return "UNKNOWN"
	}
}

// Logger interface supports the specified severities.
type Logger interface {
	WithContext(ctx string) Logger

	Log(level Level, msg string, detail ...any)
	Error(msg string, detail ...any)
	Warn(msg string, detail ...any)
	Info(msg string, detail ...any)
	Debug(msg string, detail ...any)
	Trace(msg string, detail ...any)

	// Optional: expose config getters/setters
	SetEnabled(enabled bool)
	SetLevel(level Level)
	GetLevel() Level
	SetOutput(w io.Writer)
}
