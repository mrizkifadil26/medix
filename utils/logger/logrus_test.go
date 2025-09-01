package logger_test

import (
	"bytes"
	"regexp"
	"strings"
	"testing"

	"github.com/mrizkifadil26/medix/utils/logger"
)

var ansi = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func TestLogrusLogger_BasicLogging(t *testing.T) {
	var buf bytes.Buffer
	l := logger.NewLogrusLogger()
	l.SetEnabled(true)
	l.SetLevel(logger.LevelTrace)
	l.SetOutput(&buf) // redirect logrus output

	l.Error("error message")
	l.Warn("warn message")
	l.Info("info message")
	l.Debug("debug message")
	l.Trace("trace message")

	output := buf.String()
	expectedLogs := []struct {
		level   string
		message string
	}{
		{"ERRO", "error message"},
		{"WARN", "warn message"},
		{"INFO", "info message"},
		{"DEBU", "debug message"},
		{"TRAC", "trace message"},
	}

	for _, exp := range expectedLogs {
		if !strings.Contains(output, exp.level) {
			t.Errorf("Expected log output to contain level %q", exp.level)
		}

		if !strings.Contains(output, exp.message) {
			t.Errorf("Expected log output to contain message %q", exp.message)
		}
	}
}

func TestLogrusLogger_WithContext(t *testing.T) {
	var buf bytes.Buffer
	l := logger.NewLogrusLogger()
	l.SetEnabled(true)
	l.SetLevel(logger.LevelTrace)
	l.SetOutput(&buf)

	ctxLogger := l.WithContext("TestCtx")

	ctxLogger.Error("error message")
	ctxLogger.Warn("warn message")
	ctxLogger.Info("info message")
	ctxLogger.Debug("debug message")
	ctxLogger.Trace("trace message")

	output := buf.String()
	expectedKeywords := []string{"ERRO", "WARN", "INFO", "DEBU", "TRAC"}
	for _, kw := range expectedKeywords {
		if !strings.Contains(output, kw) {
			t.Errorf("Expected log output to contain %q, got:\n%s", kw, output)
		}
	}

	// Check if context is present in the output
	cleanOutput := stripAnsi(output)
	expectedContext := "context=TestCtx"
	if !strings.Contains(cleanOutput, expectedContext) {
		t.Errorf("Expected log output to contain %q, got:\n%s", expectedContext, output)
	}
}

func TestLogrusLogger_Disabled(t *testing.T) {
	var buf bytes.Buffer
	l := logger.NewLogrusLogger()
	l.SetEnabled(false)
	l.SetLevel(logger.LevelDebug)
	l.SetOutput(&buf)

	l.Info("should not appear")

	if buf.Len() > 0 {
		t.Errorf("Expected no log output, got: %s", buf.String())
	}
}

func TestLogrusLogger_LevelFiltering(t *testing.T) {
	var buf bytes.Buffer
	l := logger.NewLogrusLogger()
	l.SetLevel(logger.LevelWarn)
	l.SetOutput(&buf)

	l.Debug("debug message")
	if buf.Len() > 0 {
		t.Errorf("Expected no log output for debug when level is WARN, got: %s", buf.String())
	}

	l.Error("error message")
	if !strings.Contains(buf.String(), "error message") {
		t.Errorf("Expected error log to appear, got: %s", buf.String())
	}
}

func stripAnsi(s string) string {
	return ansi.ReplaceAllString(s, "")
}
