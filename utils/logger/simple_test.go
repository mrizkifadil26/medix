package logger_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/mrizkifadil26/medix/utils/logger"
)

func TestSimpleLogger_BasicLogging(t *testing.T) {
	var buf bytes.Buffer
	l := logger.NewSimpleLogger()
	l.SetEnabled(true)
	l.SetLevel(logger.LevelTrace)
	l.SetOutput(&buf)

	// No bound context, just message and details
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
		{"ERROR", "error message"},
		{"WARN", "warn message"},
		{"INFO", "info message"},
		{"DEBUG", "debug message"},
		{"TRACE", "trace message"},
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

func TestSimpleLogger_WithContext(t *testing.T) {
	var buf bytes.Buffer
	l := logger.NewSimpleLogger()
	l.SetEnabled(true)
	l.SetLevel(logger.LevelTrace)
	l.SetOutput(&buf)

	ctxLogger := l.WithContext("TestCtx")

	ctxLogger.Error("error message")
	ctxLogger.Warn("warn message")
	ctxLogger.Info("info message")
	ctxLogger.Debug("debug message")
	ctxLogger.Trace("trace message")

	logOutput := buf.String()
	expectedKeywords := []string{"ERROR", "WARN", "INFO", "DEBUG", "TRACE", "TestCtx"}
	for _, kw := range expectedKeywords {
		if !strings.Contains(logOutput, kw) {
			t.Errorf("Expected log output to contain %q, got:\n%s", kw, logOutput)
		}

		// Context should be inside brackets: [TestCtx]
		if !strings.Contains(logOutput, "[TestCtx]") {
			t.Errorf("Expected log output to contain context '[TestCtx]', got:\n%s", logOutput)
		}
	}
}

func TestSimpleLogger_Disabled(t *testing.T) {
	var buf bytes.Buffer
	l := logger.NewSimpleLogger()
	l.SetEnabled(false)
	l.SetOutput(&buf)

	l.Info("should not appear")

	if buf.Len() > 0 {
		t.Errorf("Expected no log output, got: %s", buf.String())
	}
}

func TestSimpleLogger_LevelFiltering(t *testing.T) {
	var buf bytes.Buffer
	l := logger.NewSimpleLogger()
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
