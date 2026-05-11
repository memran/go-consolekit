package console

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestLoggerInfo(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger().WithWriter(&buf)
	logger.Info("hello")

	got := buf.String()
	if !strings.Contains(got, "hello") {
		t.Fatalf("expected 'hello' in output, got '%s'", got)
	}
	if !strings.Contains(got, "INFO") {
		t.Fatalf("expected 'INFO' in output, got '%s'", got)
	}
}

func TestLoggerDebug(t *testing.T) {
	var buf bytes.Buffer
	NewLogger().WithWriter(&buf).Debug("debug msg")

	got := buf.String()
	if !strings.Contains(got, "debug msg") {
		t.Fatalf("expected 'debug msg' in output, got '%s'", got)
	}
	if !strings.Contains(got, "DEBUG") {
		t.Fatalf("expected 'DEBUG' in output, got '%s'", got)
	}
}

func TestLoggerWarning(t *testing.T) {
	var buf bytes.Buffer
	NewLogger().WithWriter(&buf).Warning("warning msg")

	got := buf.String()
	if !strings.Contains(got, "warning msg") {
		t.Fatalf("expected 'warning msg' in output, got '%s'", got)
	}
	if !strings.Contains(got, "WARNING") {
		t.Fatalf("expected 'WARNING' in output, got '%s'", got)
	}
}

func TestLoggerError(t *testing.T) {
	var buf bytes.Buffer
	NewLogger().WithWriter(&buf).Error("error msg")

	got := buf.String()
	if !strings.Contains(got, "error msg") {
		t.Fatalf("expected 'error msg' in output, got '%s'", got)
	}
	if !strings.Contains(got, "ERROR") {
		t.Fatalf("expected 'ERROR' in output, got '%s'", got)
	}
}

func TestLoggerNotice(t *testing.T) {
	var buf bytes.Buffer
	NewLogger().WithWriter(&buf).Notice("notice msg")

	got := buf.String()
	if !strings.Contains(got, "NOTICE") {
		t.Fatalf("expected 'NOTICE' in output, got '%s'", got)
	}
}

func TestLoggerCritical(t *testing.T) {
	var buf bytes.Buffer
	NewLogger().WithWriter(&buf).Critical("critical msg")

	got := buf.String()
	if !strings.Contains(got, "CRITICAL") {
		t.Fatalf("expected 'CRITICAL' in output, got '%s'", got)
	}
}

func TestLoggerAlert(t *testing.T) {
	var buf bytes.Buffer
	NewLogger().WithWriter(&buf).Alert("alert msg")

	got := buf.String()
	if !strings.Contains(got, "ALERT") {
		t.Fatalf("expected 'ALERT' in output, got '%s'", got)
	}
}

func TestLoggerEmergency(t *testing.T) {
	var buf bytes.Buffer
	NewLogger().WithWriter(&buf).Emergency("emergency msg")

	got := buf.String()
	if !strings.Contains(got, "EMERGENCY") {
		t.Fatalf("expected 'EMERGENCY' in output, got '%s'", got)
	}
}

func TestLoggerLogLevel(t *testing.T) {
	var buf bytes.Buffer
	NewLogger().WithWriter(&buf).Log(DebugLevel, "level msg")

	got := buf.String()
	if !strings.Contains(got, "DEBUG") {
		t.Fatalf("expected 'DEBUG' in output, got '%s'", got)
	}
	if !strings.Contains(got, "level msg") {
		t.Fatalf("expected 'level msg' in output, got '%s'", got)
	}
}

func TestLoggerWithContextSingle(t *testing.T) {
	var buf bytes.Buffer
	NewLogger().WithWriter(&buf).With("user_id", 123).Info("user login")

	got := buf.String()
	if !strings.Contains(got, "user login") {
		t.Fatalf("expected 'user login' in output, got '%s'", got)
	}
	if !strings.Contains(got, `"user_id"`) || !strings.Contains(got, `123`) {
		t.Fatalf("expected context in output, got '%s'", got)
	}
}

func TestLoggerWithContextChained(t *testing.T) {
	var buf bytes.Buffer
	NewLogger().WithWriter(&buf).
		With("user_id", 123).
		With("role", "admin").
		Info("chained context")

	got := buf.String()
	if !strings.Contains(got, `"user_id"`) || !strings.Contains(got, `"role"`) {
		t.Fatalf("expected both context keys, got '%s'", got)
	}
}

func TestLoggerWithMap(t *testing.T) {
	var buf bytes.Buffer
	NewLogger().WithWriter(&buf).
		WithMap(map[string]interface{}{"key1": "val1", "key2": 42}).
		Info("map context")

	got := buf.String()
	if !strings.Contains(got, `"key1"`) || !strings.Contains(got, `"key2"`) {
		t.Fatalf("expected map keys in output, got '%s'", got)
	}
}

func TestLoggerWithKVContext(t *testing.T) {
	var buf bytes.Buffer
	NewLogger().WithWriter(&buf).Info("order placed", "order_id", 1001, "total", 49.99)

	got := buf.String()
	if !strings.Contains(got, `"order_id"`) || !strings.Contains(got, `1001`) {
		t.Fatalf("expected kv context in output, got '%s'", got)
	}
}

func TestLoggerWithContextNotMutated(t *testing.T) {
	var buf1, buf2 bytes.Buffer
	logger := NewLogger()

	logger.WithWriter(&buf1).With("key", "val1").Info("first")
	logger.WithWriter(&buf2).Info("second")

	if strings.Contains(buf2.String(), "val1") {
		t.Fatal("context from first call should not leak to second")
	}
}

func TestLoggerChannel(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger().AddChannel("file", &buf)

	logger.Channel("file").Info("channel msg")
	got := buf.String()
	if !strings.Contains(got, "channel msg") {
		t.Fatalf("expected 'channel msg' in output, got '%s'", got)
	}
}

func TestLoggerStack(t *testing.T) {
	var buf1, buf2 bytes.Buffer
	logger := NewLogger().
		AddChannel("file", &buf1).
		AddChannel("slack", &buf2)

	logger.Stack("file", "slack").Info("stack msg")

	if !strings.Contains(buf1.String(), "stack msg") {
		t.Fatalf("expected in buf1, got '%s'", buf1.String())
	}
	if !strings.Contains(buf2.String(), "stack msg") {
		t.Fatalf("expected in buf2, got '%s'", buf2.String())
	}
}

func TestLoggerWithName(t *testing.T) {
	var buf bytes.Buffer
	NewLogger().WithWriter(&buf).WithName("app").Info("named")

	got := buf.String()
	if !strings.Contains(got, "app.INFO") {
		t.Fatalf("expected 'app.INFO' in output, got '%s'", got)
	}
}

func TestLoggerDefaultWritesToStdout(t *testing.T) {
	logger := NewLogger()
	if logger.writers == nil || len(logger.writers) == 0 {
		t.Fatal("expected default writers")
	}
}

func TestLoggerLevelFormat(t *testing.T) {
	levels := []struct {
		level LogLevel
		name  string
	}{
		{DebugLevel, "DEBUG"},
		{InfoLevel, "INFO"},
		{NoticeLevel, "NOTICE"},
		{WarningLevel, "WARNING"},
		{ErrorLevel, "ERROR"},
		{CriticalLevel, "CRITICAL"},
		{AlertLevel, "ALERT"},
		{EmergencyLevel, "EMERGENCY"},
	}

	for _, l := range levels {
		if l.level.String() != l.name {
			t.Fatalf("expected '%s', got '%s'", l.name, l.level.String())
		}
	}
}

func TestLoggerFormatTimestamp(t *testing.T) {
	var buf bytes.Buffer
	NewLogger().WithWriter(&buf).Info("timestamp test")

	got := buf.String()
	if len(got) < 20 {
		t.Fatalf("expected timestamp prefix, got too short: '%s'", got)
	}
	if got[0] != '[' {
		t.Fatalf("expected '[' at start, got '%s'", got[:1])
	}
}

func TestLoggerChannelNotFound(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger().WithWriter(&buf)

	logger.Channel("nonexistent").Info("fallback")
	got := buf.String()
	if !strings.Contains(got, "fallback") {
		t.Fatalf("expected fallback to default writer, got '%s'", got)
	}
}

func TestLoggerStackEmpty(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger().WithWriter(&buf)

	logger.Stack().Info("empty stack")
	got := buf.String()
	if !strings.Contains(got, "empty stack") {
		t.Fatalf("expected fallback when stack is empty, got '%s'", got)
	}
}

func TestLoggerAddChannelMultipleWriters(t *testing.T) {
	var buf1, buf2 bytes.Buffer
	logger := NewLogger().AddChannel("multi", &buf1)
	logger.AddChannel("multi", &buf2)

	logger.Channel("multi").Info("multi writer")
	if !strings.Contains(buf1.String(), "multi writer") {
		t.Fatal("expected first writer to receive message")
	}
	if !strings.Contains(buf2.String(), "multi writer") {
		t.Fatal("expected second writer to receive message")
	}
}

func TestLoggerWithAfterChannel(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger().AddChannel("file", &buf)

	logger.Channel("file").With("env", "test").Info("with after channel")

	got := buf.String()
	if !strings.Contains(got, `"env"`) {
		t.Fatalf("expected context after channel, got '%s'", got)
	}
}

func TestLoggerNoPanicOnEmptyMsg(t *testing.T) {
	var buf bytes.Buffer
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("panic on empty message: %v", r)
		}
	}()
	NewLogger().WithWriter(&buf).Info("")
}

func TestLoggerDefaultFunctions(t *testing.T) {
	var buf bytes.Buffer
	old := defaultLogger.writers
	defaultLogger.writers = []io.Writer{&buf}
	defer func() { defaultLogger.writers = old }()

	Info("default info")
	if !strings.Contains(buf.String(), "default info") {
		t.Fatalf("expected default info in output, got '%s'", buf.String())
	}
}

func TestLoggerDefaultError(t *testing.T) {
	var buf bytes.Buffer
	old := defaultLogger.writers
	defaultLogger.writers = []io.Writer{&buf}
	defer func() { defaultLogger.writers = old }()

	Error("default error")
	if !strings.Contains(buf.String(), "default error") {
		t.Fatal("default Error failed")
	}
}
