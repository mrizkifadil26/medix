package logger

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
	Log(level Level, context, msg string, detail any)
	Error(context, msg string, detail any)
	Warn(context, msg string, detail any)
	Info(context, msg string, detail any)
	Debug(context, msg string, detail any)
	Trace(context, msg string, detail any)

	// Optional: expose config getters/setters
	SetEnabled(enabled bool)
	SetLevel(level Level)
	GetLevel() Level
}
