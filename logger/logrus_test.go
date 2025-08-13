package logger_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/mrizkifadil26/medix/logger"
)

func TestLogrusLogger_BasicLogging(t *testing.T) {
	var buf bytes.Buffer
	l := logger.NewLogrusLogger()
	l.SetEnabled(true)
	l.SetLevel(logger.LevelTrace)
	l.SetOutput(&buf) // redirect logrus output

	l.Error("TestCtx", "error message", nil)
	l.Warn("TestCtx", "warn message", nil)
	l.Info("TestCtx", "info message", nil)
	l.Debug("TestCtx", "debug message", nil)
	l.Trace("TestCtx", "trace message", nil)

	output := buf.String()
	expectedKeywords := []string{"error message", "warn message", "info message", "debug message", "trace message"}
	for _, kw := range expectedKeywords {
		if !strings.Contains(output, kw) {
			t.Errorf("Expected log output to contain %q, got:\n%s", kw, output)
		}
	}
}

func TestLogrusLogger_Disabled(t *testing.T) {
	var buf bytes.Buffer
	l := logger.NewLogrusLogger()
	l.SetEnabled(false)
	l.SetLevel(logger.LevelDebug)
	l.SetOutput(&buf)

	l.Info("TestCtx", "should not appear", nil)

	if buf.Len() > 0 {
		t.Errorf("Expected no log output, got: %s", buf.String())
	}
}

func TestLogrusLogger_LevelFiltering(t *testing.T) {
	var buf bytes.Buffer
	l := logger.NewLogrusLogger()
	l.SetLevel(logger.LevelWarn)
	l.SetOutput(&buf)

	l.Debug("TestCtx", "debug message", nil)
	if buf.Len() > 0 {
		t.Errorf("Expected no log output for debug when level is WARN, got: %s", buf.String())
	}

	l.Error("TestCtx", "error message", nil)
	if !strings.Contains(buf.String(), "error message") {
		t.Errorf("Expected error log to appear, got: %s", buf.String())
	}
}
