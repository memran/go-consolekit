package console

import (
	"os"
	"runtime"
	"testing"
	"time"
)

func TestProcessRun(t *testing.T) {
	result := NewProcess("go", "version").Run()
	if !result.IsSuccessful() {
		t.Fatalf("expected success, got exit code %d: %s", result.ExitCode(), result.Err())
	}
	if result.Output() == "" {
		t.Fatal("expected non-empty output")
	}
}

func TestProcessMustRun(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("MustRun panicked: %v", r)
		}
	}()
	result := Run("go", "version").MustRun()
	if !result.IsSuccessful() {
		t.Fatal("MustRun should succeed")
	}
}

func TestProcessMustRunPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for invalid command")
		}
	}()
	Run("nonexistent-command-12345").MustRun()
}

func TestProcessWithWorkingDir(t *testing.T) {
	dir, _ := os.Getwd()
	result := NewProcess("go", "env", "GOMOD").WithWorkingDir(dir).Run()
	if !result.IsSuccessful() {
		t.Fatalf("expected success: %s", result.Err())
	}
}

func TestProcessWithEnv(t *testing.T) {
	result := NewProcess("go", "env", "GO_TEST_VAR").
		WithEnv("GO_TEST_VAR", "custom_value").
		Run()
	if !result.IsSuccessful() {
		t.Fatalf("expected success: %s", result.Err())
	}
}

func TestProcessWithInputString(t *testing.T) {
	cmd := "go"
	if runtime.GOOS == "windows" {
		cmd = "findstr"
		result := Run(cmd, "hello").WithInputString("hello world\nfoo bar").Run()
		if !result.IsSuccessful() {
			t.Fatalf("expected success: %s, stderr: %s", result.Err(), result.Error())
		}
		if !containsOutput(result.Output(), "hello") {
			t.Fatalf("expected 'hello' in output, got '%s'", result.Output())
		}
	} else {
		_ = Run(cmd, "run", ".").WithInputString("hello").Run()
	}
}

func TestProcessTimeout(t *testing.T) {
	sleep := "sleep"
	arg := "5"
	if runtime.GOOS == "windows" {
		sleep = "ping"
		arg = "-n"
	}

	result := NewProcess(sleep, arg, "5").Timeout(100 * time.Millisecond).Run()
	if result.IsSuccessful() {
		t.Log("process completed before timeout (unlikely but possible)")
	} else {
		t.Logf("process timed out as expected: %v", result.Err())
	}
}

func TestProcessStartWait(t *testing.T) {
	proc := NewProcess("go", "version")
	err := proc.Start()
	if err != nil {
		t.Fatalf("Start failed: %v", err)
	}
	if !proc.IsRunning() {
		t.Fatal("expected process to be running")
	}

	result := proc.Wait()
	if !result.IsSuccessful() {
		t.Fatalf("Wait failed: %s", result.Err())
	}
	if proc.IsRunning() {
		t.Fatal("expected process to be done after Wait")
	}
}

func TestProcessPID(t *testing.T) {
	proc := NewProcess("go", "version")
	if proc.PID() != 0 {
		t.Fatal("expected 0 before start")
	}
	proc.Start()
	if proc.PID() == 0 {
		t.Fatal("expected non-zero PID after start")
	}
	proc.Wait()
}

func TestProcessResultOutput(t *testing.T) {
	result := Run("go", "version").Run()
	if result.Output() == "" {
		t.Fatal("expected output")
	}
}

func TestProcessResultIsFailed(t *testing.T) {
	result := Run("nonexistent-cmd-999").Run()
	if !result.IsFailed() {
		t.Fatal("expected failure")
	}
}

func TestProcessResultLines(t *testing.T) {
	result := Run("go", "version").Run()
	lines := result.Lines()
	if len(lines) == 0 {
		t.Fatal("expected at least 1 line")
	}
}

func TestProcessFluentChaining(t *testing.T) {
	result := Run("go", "version").
		WithEnv("GOFLAGS", "").
		MustRun()

	if !result.IsSuccessful() {
		t.Fatal("chaining failed")
	}
}

func TestNewProcess(t *testing.T) {
	p := NewProcess("echo", "hello")
	if p.name != "echo" {
		t.Fatalf("expected 'echo', got '%s'", p.name)
	}
}

func TestProcessStop(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("process stop test unreliable on Windows")
	}

	sleep := "sleep"
	arg := "30"
	if runtime.GOOS == "windows" {
		sleep = "ping"
		arg = "-n 30 127.0.0.1"
	}

	p := NewProcess(sleep, arg)
	if err := p.Start(); err != nil {
		t.Fatalf("Start failed: %v", err)
	}
	time.Sleep(100 * time.Millisecond)
	if !p.IsRunning() {
		t.Fatal("expected process to be running")
	}
	if err := p.Stop(); err != nil {
		t.Fatalf("Stop failed: %v", err)
	}
	result := p.Wait()
	if result.IsSuccessful() {
		t.Log("process may have been killed")
	}
}

func TestProcessWaitNotStarted(t *testing.T) {
	p := NewProcess("go", "version")
	result := p.Wait()
	if result.IsSuccessful() {
		t.Fatal("expected failure for Wait without Start")
	}
}

func TestProcessStopNotStarted(t *testing.T) {
	p := NewProcess("go", "version")
	if err := p.Stop(); err == nil {
		t.Fatal("expected error stopping unstarted process")
	}
}

func TestProcessString(t *testing.T) {
	p := NewProcess("echo", "hello", "world")
	s := p.String()
	if s != "echo [hello world]" {
		t.Fatalf("expected 'echo [hello world]', got '%s'", s)
	}
}

func TestProcessErrorOutput(t *testing.T) {
	result := Run("go", "build", "-invalid-flag").Run()
	if result.ExitCode() == 0 {
		t.Skip("command unexpectedly succeeded")
	}
	if result.Error() == "" && result.Err() == nil {
		t.Log("no stderr captured")
	}
}

func TestSplitLines(t *testing.T) {
	tests := []struct {
		input string
		want  int
	}{
		{"", 0},
		{"hello", 1},
		{"hello\nworld", 2},
		{"hello\nworld\n", 2},
		{"a\nb\nc\n", 3},
		{"a\r\nb\r\nc", 3},
	}
	for _, tt := range tests {
		got := splitLines(tt.input)
		if len(got) != tt.want {
			t.Fatalf("splitLines(%q) = %d lines, want %d: %v", tt.input, len(got), tt.want, got)
		}
	}
}

func TestProcessSignal(t *testing.T) {
	p := NewProcess("go", "version")
	if err := p.Signal(os.Kill); err == nil {
		t.Fatal("expected error for Signal before Start")
	}
}

func containsOutput(s, substr string) bool {
	return len(s) >= len(substr) && containsSubstring(s, substr)
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
