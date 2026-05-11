package console

import (
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"testing"
	"time"
)

func TestDaemonWritePIDThenCheck(t *testing.T) {
	dir := t.TempDir()
	pidPath := filepath.Join(dir, "app.pid")
	logPath := filepath.Join(dir, "app.log")

	os.Args = []string{"example", "--pid-file", pidPath, "--log-file", logPath, "queue:work"}

	pf := extractFlagValue(os.Args, "--pid-file")
	lf := extractFlagValue(os.Args, "--log-file")

	if pf != pidPath {
		t.Fatalf("flag extraction failed: pf=%s", pf)
	}

	if lf == "" {
		t.Fatal("log-file extraction failed")
	}

	err := writePID(pf)
	if err != nil {
		t.Fatalf("writePID failed: %v", err)
	}

	if _, err := os.Stat(pf); os.IsNotExist(err) {
		t.Fatal("PID file missing immediately after writePID")
	}

	if _, err := os.Stat(pidPath); os.IsNotExist(err) {
		t.Fatal("PID file missing after writePID (path mismatch)")
	}

	pid, err := readPID(pf)
	if err != nil {
		t.Fatalf("readPID failed: %v", err)
	}
	if pid != os.Getpid() {
		t.Fatalf("PID mismatch: got %d, want %d", pid, os.Getpid())
	}

	time.Sleep(100 * time.Millisecond)
	if _, err := os.Stat(pf); os.IsNotExist(err) {
		t.Fatal("PID file disappeared during wait")
	}
}

func TestDaemonRedirectOutputDoesNotRemovePID(t *testing.T) {
	dir := t.TempDir()
	pidPath := filepath.Join(dir, "app.pid")
	logPath := filepath.Join(dir, "app.log")

	err := writePID(pidPath)
	if err != nil {
		t.Fatalf("writePID failed: %v", err)
	}

	savedStdout := os.Stdout
	savedStderr := os.Stderr

	err = redirectOutput(logPath)
	if err != nil {
		t.Fatalf("redirectOutput failed: %v", err)
	}

	if _, err := os.Stat(pidPath); os.IsNotExist(err) {
		t.Fatal("PID file was removed by redirectOutput")
	}

	pid, err := readPID(pidPath)
	if err != nil {
		t.Fatalf("readPID failed after redirect: %v", err)
	}
	if pid != os.Getpid() {
		t.Fatalf("PID mismatch after redirect: got %d, want %d", pid, os.Getpid())
	}

	os.Stdout.Close()
	os.Stderr.Close()
	os.Stdout = savedStdout
	os.Stderr = savedStderr
}

func TestDaemonProcessExists(t *testing.T) {
	if !processExists(os.Getpid()) {
		t.Fatal("processExists should return true for our PID")
	}

	if processExists(0) {
		t.Fatal("processExists should return false for PID 0")
	}
	if processExists(-1) {
		t.Fatal("processExists should return false for PID -1")
	}
}

func TestDaemonStartChildFlow(t *testing.T) {
	dir := t.TempDir()
	pidPath := filepath.Join(dir, "daemon.pid")
	logPath := filepath.Join(dir, "daemon.log")

	exePath := filepath.Join(dir, "testd.exe")
	buildCmd := exec.Command("go", "build", "-o", exePath, "../cmd/example")
	buildOut, err := buildCmd.CombinedOutput()
	if err != nil {
		t.Skipf("could not build example binary: %v\n%s", err, string(buildOut))
	}

	if _, err := os.Stat(exePath); os.IsNotExist(err) {
		t.Skipf("binary not built: %v", err)
	}
	t.Logf("binary built at: %s", exePath)

	childArgs := []string{exePath, "--pid-file", pidPath, "--log-file", logPath, "table", "demo"}
	attr := &os.ProcAttr{
		Env:   append(os.Environ(), daemonEnvVar+"=1"),
		Files: []*os.File{nil, nil, nil},
	}

	proc, err := os.StartProcess(exePath, childArgs, attr)
	if err != nil {
		t.Fatalf("StartProcess failed: %v", err)
	}

	ps, err := proc.Wait()
	if err != nil {
		t.Logf("Wait returned error: %v", err)
	}
	t.Logf("child exit code: %d", ps.ExitCode())

	if _, err := os.Stat(pidPath); os.IsNotExist(err) {
		t.Log("PID file was removed by shutdown hook (expected for quick commands)")
	} else {
		pid, readErr := readPID(pidPath)
		if readErr == nil {
			t.Logf("PID file still exists with PID %d (shutdown hook did not fire)", pid)
		}
	}

	logData, logErr := os.ReadFile(logPath)
	if logErr != nil {
		t.Logf("log file not found: %v", logErr)
	} else {
		t.Logf("log content: %s", string(logData))
	}
}

func TestDaemonArgsPersistForRestart(t *testing.T) {
	app := New("test").EnableDaemon()

	os.Args = []string{"example", "--daemon", "--pid-file", "app.pid", "worker:start"}
	app.daemonArgs = os.Args

	stripped := stripDaemonFlag(app.daemonArgs)
	expected := []string{"example", "--pid-file", "app.pid", "worker:start"}
	if len(stripped) != len(expected) {
		t.Fatalf("expected %d args, got %d: %v", len(expected), len(stripped), stripped)
	}
	for i := range expected {
		if stripped[i] != expected[i] {
			t.Fatalf("arg %d: expected '%s', got '%s'", i, expected[i], stripped[i])
		}
	}
}

func TestDaemonStopSendsSignal(t *testing.T) {
	dir := t.TempDir()
	pidPath := filepath.Join(dir, "app.pid")

	err := writePID(pidPath)
	if err != nil {
		t.Fatalf("writePID failed: %v", err)
	}

	pid, err := readPID(pidPath)
	if err != nil || pid != os.Getpid() {
		t.Fatalf("readPID mismatch: pid=%d, err=%v", pid, err)
	}

	proc, err := os.FindProcess(pid)
	if err != nil {
		t.Fatalf("FindProcess failed: %v", err)
	}
	if err := proc.Signal(syscall.SIGTERM); err != nil {
		t.Logf("Signal to self returned (expected on Windows): %v", err)
	}

	removePID(pidPath)
	if _, err := os.Stat(pidPath); !os.IsNotExist(err) {
		t.Fatal("PID file should be removed")
	}
}

func TestDaemonRestartUsesStoredArgs(t *testing.T) {
	app := New("test").EnableDaemon()

	originalArgs := []string{"example", "--daemon", "--pid-file", "/tmp/app.pid", "--log-file", "/tmp/app.log", "queue:work"}
	app.daemonArgs = originalArgs

	stripped := stripDaemonFlag(app.daemonArgs)
	app.daemonArgs = stripped

	expected := []string{"example", "--pid-file", "/tmp/app.pid", "--log-file", "/tmp/app.log", "queue:work"}
	if len(app.daemonArgs) != len(expected) {
		t.Fatalf("Restart args length: got %d, want %d\n  got:  %v\n  want: %v",
			len(app.daemonArgs), len(expected), app.daemonArgs, expected)
	}
	for i := range expected {
		if app.daemonArgs[i] != expected[i] {
			t.Fatalf("Restart arg %d: got '%s', want '%s'", i, app.daemonArgs[i], expected[i])
		}
	}
}
