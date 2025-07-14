package logger

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

var Log = logrus.New()

func init() {
	Log.SetFormatter(&logrus.TextFormatter{
		DisableTimestamp: true,
		ForceColors:      true,
		FullTimestamp:    false,
	})
	Log.SetLevel(logrus.InfoLevel)
}

// ANSI colors
const (
	colorReset   = "\033[0m"
	colorRed     = "\033[31m"
	colorGreen   = "\033[32m"
	colorYellow  = "\033[33m"
	colorBlue    = "\033[34m"
	colorMagenta = "\033[35m"
	colorCyan    = "\033[36m"
)

// generic tag printer using logrus
func logWithTag(level logrus.Level, tag, color, msg string) {
	colored := fmt.Sprintf("%s[%s]%s %s", color, tag, colorReset, msg)

	switch level {
	case logrus.DebugLevel:
		Log.Debug(colored)
	case logrus.InfoLevel:
		Log.Info(colored)
	case logrus.WarnLevel:
		Log.Warn(colored)
	case logrus.ErrorLevel:
		Log.Error(colored)
	case logrus.FatalLevel:
		Log.Fatal(colored)
	default:
		Log.Print(colored)
	}
}

// Custom tagged log helpers
func Step(msg string) {
	logWithTag(logrus.InfoLevel, "STEP", colorCyan, msg)
}

func Done(msg string) {
	logWithTag(logrus.InfoLevel, "DONE", colorGreen, msg)
}

func Warn(msg string) {
	logWithTag(logrus.WarnLevel, "WARN", colorYellow, msg)
}

func Error(msg string) {
	logWithTag(logrus.ErrorLevel, "ERROR", colorRed, msg)
}

func Dry(msg string) {
	logWithTag(logrus.InfoLevel, "DRY-RUN", colorMagenta, msg)
}

func Watch(msg string) {
	logWithTag(logrus.InfoLevel, "WATCH", colorYellow, msg)
}

func Info(msg string) {
	logWithTag(logrus.InfoLevel, "INFO", colorBlue, msg)
}

// Flexible for custom tags
func Println(tag string, color string, msg string) {
	tag = strings.ToUpper(tag)
	logWithTag(logrus.InfoLevel, tag, color, msg)
}
