package console

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestRegistryAddFind(t *testing.T) {
	r := NewRegistry()
	cmd := &testCommand{name: "test"}
	r.Add(cmd)

	found, err := r.Find("test")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if found.Name() != "test" {
		t.Fatalf("expected 'test', got '%s'", found.Name())
	}
}

func TestRegistryFindNotFound(t *testing.T) {
	r := NewRegistry()
	_, err := r.Find("nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent command")
	}
	if _, ok := err.(*CommandNotFoundError); !ok {
		t.Fatalf("expected CommandNotFoundError, got %T", err)
	}
}

func TestRegistryAll(t *testing.T) {
	r := NewRegistry()
	r.Add(&testCommand{name: "a"})
	r.Add(&testCommand{name: "b"})

	all := r.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 commands, got %d", len(all))
	}
}

func TestCommandConfigArguments(t *testing.T) {
	cfg := NewCommandConfig("test")
	cfg.Argument("name").Required().Description("The name")
	cfg.Argument("age").Default("30")

	if len(cfg.Arguments) != 2 {
		t.Fatalf("expected 2 arguments, got %d", len(cfg.Arguments))
	}
	if !cfg.Arguments[0].required {
		t.Fatal("expected first argument to be required")
	}
	if cfg.Arguments[1].defaultVal != "30" {
		t.Fatalf("expected default '30', got '%s'", cfg.Arguments[1].defaultVal)
	}
}

func TestCommandConfigOptions(t *testing.T) {
	cfg := NewCommandConfig("test")
	cfg.Option("db").Shortcut("d").Default("sqlite").Description("Database")
	cfg.Option("verbose").Shortcut("v").Required()

	if len(cfg.Options) != 2 {
		t.Fatalf("expected 2 options, got %d", len(cfg.Options))
	}
	if cfg.Options[0].shortcut != "d" {
		t.Fatalf("expected shortcut 'd', got '%s'", cfg.Options[0].shortcut)
	}
	if cfg.Options[0].defaultVal != "sqlite" {
		t.Fatalf("expected default 'sqlite', got '%s'", cfg.Options[0].defaultVal)
	}
	if !cfg.Options[1].required {
		t.Fatal("expected second option to be required")
	}
}

func TestContextArgAndOption(t *testing.T) {
	ctx := newContext(context.Background())
	ctx.argsMap["name"] = "test"
	ctx.optionsMap["db"] = "postgres"

	if ctx.Arg("name") != "test" {
		t.Fatalf("expected 'test', got '%s'", ctx.Arg("name"))
	}
	if ctx.Option("db") != "postgres" {
		t.Fatalf("expected 'postgres', got '%s'", ctx.Option("db"))
	}
}

func TestContextMissingArg(t *testing.T) {
	ctx := newContext(context.Background())
	if ctx.Arg("missing") != "" {
		t.Fatalf("expected empty string, got '%s'", ctx.Arg("missing"))
	}
}

func TestContextDone(t *testing.T) {
	ctx := newContext(context.Background())
	if ctx.IsCancelled() {
		t.Fatal("expected context not to be cancelled initially")
	}
	select {
	case <-ctx.Done():
		t.Fatal("Done() should not fire before cancellation")
	default:
	}
}

func TestContextCancellation(t *testing.T) {
	ctx := newContext(context.Background())
	ctx.cancel()
	if !ctx.IsCancelled() {
		t.Fatal("expected context to be cancelled")
	}
	select {
	case <-ctx.Done():
	default:
		t.Fatal("Done() should fire after cancellation")
	}
}

func TestContextShutdownCallbackRegistration(t *testing.T) {
	ctx := newContext(context.Background())
	called := false
	ctx.OnShutdown(func() error {
		called = true
		return nil
	})
	if called {
		t.Fatal("callback should not be called before shutdown")
	}
}

func TestContextShutdownCallbackExecution(t *testing.T) {
	ctx := newContext(context.Background())
	called := false
	ctx.OnShutdown(func() error {
		called = true
		return nil
	})
	ctx.runShutdown()
	if !called {
		t.Fatal("callback should have been called after shutdown")
	}
}

func TestContextShutdownCallbackMultiple(t *testing.T) {
	ctx := newContext(context.Background())
	order := []string{}
	ctx.OnShutdown(func() error {
		order = append(order, "first")
		return nil
	})
	ctx.OnShutdown(func() error {
		order = append(order, "second")
		return nil
	})
	ctx.runShutdown()
	if len(order) != 2 {
		t.Fatalf("expected 2 callbacks, got %d", len(order))
	}
	if order[0] != "first" {
		t.Fatalf("expected first callback to run first, got '%s'", order[0])
	}
	if order[1] != "second" {
		t.Fatalf("expected second callback to run second, got '%s'", order[1])
	}
}

func TestContextShutdownCallbackOnce(t *testing.T) {
	ctx := newContext(context.Background())
	count := 0
	ctx.OnShutdown(func() error {
		count++
		return nil
	})
	ctx.runShutdown()
	ctx.runShutdown()
	if count != 1 {
		t.Fatalf("expected callback to run exactly once, got %d", count)
	}
}

func TestContextShutdownTriggersCancel(t *testing.T) {
	ctx := newContext(context.Background())
	ctx.output = NewOutput(NewCLIRenderer())
	ctx.Shutdown("test shutdown")
	if !ctx.IsCancelled() {
		t.Fatal("expected context to be cancelled after Shutdown()")
	}
}

func TestContextShutdownWithError(t *testing.T) {
	ctx := newContext(context.Background())
	var gotErr error
	ctx.OnShutdown(func() error {
		return errors.New("cleanup error")
	})
	ctx.OnShutdown(func() error {
		gotErr = errors.New("second error")
		return gotErr
	})
	ctx.runShutdown()
	if gotErr == nil {
		t.Fatal("expected callback errors to be captured")
	}
}

func TestContextIsDaemon(t *testing.T) {
	ctx := newContext(context.Background())
	if ctx.IsDaemon() {
		t.Fatal("expected IsDaemon to be false in tests")
	}
}

func TestContextPID(t *testing.T) {
	ctx := newContext(context.Background())
	if ctx.PID() != os.Getpid() {
		t.Fatalf("expected PID %d, got %d", os.Getpid(), ctx.PID())
	}
}

func TestContextWritePID(t *testing.T) {
	ctx := newContext(context.Background())
	dir := t.TempDir()
	path := filepath.Join(dir, "test.pid")
	if err := ctx.WritePID(path); err != nil {
		t.Fatalf("WritePID failed: %v", err)
	}
	pid, err := readPID(path)
	if err != nil {
		t.Fatalf("readPID failed: %v", err)
	}
	if pid != os.Getpid() {
		t.Fatalf("expected PID %d, got %d", os.Getpid(), pid)
	}
}

func TestContextStopDaemon(t *testing.T) {
	ctx := newContext(context.Background())
	called := false
	ctx.OnShutdown(func() error {
		called = true
		return nil
	})
	ctx.StopDaemon()
	if !ctx.IsCancelled() {
		t.Fatal("expected context to be cancelled after StopDaemon")
	}
	if !called {
		t.Fatal("expected shutdown callback to be called")
	}
}

func TestContextRunWorkerCancellation(t *testing.T) {
	ctx := newContext(context.Background())
	done := make(chan struct{})
	go func() {
		ctx.RunWorker(func() error {
			return nil
		})
		close(done)
	}()
	ctx.cancel()
	<-done
}

func TestEnableDaemon(t *testing.T) {
	app := New("test")
	app.EnableDaemon()
	if !app.daemonEnabled {
		t.Fatal("expected daemon to be enabled")
	}
}

func TestPIDFile(t *testing.T) {
	app := New("test")
	app.PIDFile("storage/app.pid")
	if app.pidFile != "storage/app.pid" {
		t.Fatalf("expected 'storage/app.pid', got '%s'", app.pidFile)
	}
}

func TestLogFile(t *testing.T) {
	app := New("test")
	app.LogFile("storage/app.log")
	if app.logFile != "storage/app.log" {
		t.Fatalf("expected 'storage/app.log', got '%s'", app.logFile)
	}
}

func TestWriteAndReadPID(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "pid")
	if err := writePID(path); err != nil {
		t.Fatalf("writePID failed: %v", err)
	}
	pid, err := readPID(path)
	if err != nil {
		t.Fatalf("readPID failed: %v", err)
	}
	if pid != os.Getpid() {
		t.Fatalf("expected PID %d, got %d", os.Getpid(), pid)
	}
}

func TestRemovePID(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "pid")
	writePID(path)
	if err := removePID(path); err != nil {
		t.Fatalf("removePID failed: %v", err)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatal("expected PID file to be removed")
	}
}

func TestHasDaemonFlag(t *testing.T) {
	if !hasDaemonFlag([]string{"app", "--daemon", "cmd"}) {
		t.Fatal("expected --daemon flag to be detected")
	}
	if hasDaemonFlag([]string{"app", "cmd"}) {
		t.Fatal("expected no --daemon flag")
	}
}

func TestStripDaemonFlag(t *testing.T) {
	result := stripDaemonFlag([]string{"app", "--daemon", "cmd"})
	if len(result) != 2 || result[0] != "app" || result[1] != "cmd" {
		t.Fatal("expected --daemon to be stripped")
	}
}

func TestExtractFlagValue(t *testing.T) {
	args := []string{"app", "--pid-file", "test.pid", "--log-file=app.log"}
	if v := extractFlagValue(args, "--pid-file"); v != "test.pid" {
		t.Fatalf("expected 'test.pid', got '%s'", v)
	}
	if v := extractFlagValue(args, "--log-file"); v != "app.log" {
		t.Fatalf("expected 'app.log', got '%s'", v)
	}
}

func TestCommandBuilderRegistration(t *testing.T) {
	app := New("test")

	app.Command("greet").
		Description("Greet someone").
		Argument("name").
		Required().
		Description("Name to greet").
		Handle(func(ctx *Context) error {
			ctx.Success("Hello " + ctx.Arg("name"))
			return nil
		})

	cmd, err := app.registry.Find("greet")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if cmd.Name() != "greet" {
		t.Fatalf("expected 'greet', got '%s'", cmd.Name())
	}
	if cmd.Description() != "Greet someone" {
		t.Fatalf("expected 'Greet someone', got '%s'", cmd.Description())
	}

	cfg := NewCommandConfig(cmd.Name())
	cmd.Configure(cfg)
	if len(cfg.Arguments) != 1 {
		t.Fatalf("expected 1 argument, got %d", len(cfg.Arguments))
	}
	if cfg.Arguments[0].Name != "name" {
		t.Fatalf("expected 'name', got '%s'", cfg.Arguments[0].Name)
	}
	if !cfg.Arguments[0].required {
		t.Fatal("expected argument to be required")
	}

	ctx := newContext(context.Background())
	ctx.output = NewOutput(NewCLIRenderer())
	ctx.argsMap["name"] = "World"
	err = cmd.Handle(ctx)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestRendererLine(t *testing.T) {
	r := NewCLIRenderer()
	defer func() {
		if err := recover(); err != nil {
			t.Fatalf("Line() panicked: %v", err)
		}
	}()
	r.Line("test")
}

func TestRendererInfo(t *testing.T) {
	r := NewCLIRenderer()
	defer func() {
		if err := recover(); err != nil {
			t.Fatalf("Info() panicked: %v", err)
		}
	}()
	r.Info("test")
}

func TestRendererSuccess(t *testing.T) {
	r := NewCLIRenderer()
	defer func() {
		if err := recover(); err != nil {
			t.Fatalf("Success() panicked: %v", err)
		}
	}()
	r.Success("test")
}

func TestRendererWarning(t *testing.T) {
	r := NewCLIRenderer()
	defer func() {
		if err := recover(); err != nil {
			t.Fatalf("Warning() panicked: %v", err)
		}
	}()
	r.Warning("test")
}

func TestRendererError(t *testing.T) {
	r := NewCLIRenderer()
	defer func() {
		if err := recover(); err != nil {
			t.Fatalf("Error() panicked: %v", err)
		}
	}()
	r.Error("test")
}

func TestRendererTitle(t *testing.T) {
	r := NewCLIRenderer()
	defer func() {
		if err := recover(); err != nil {
			t.Fatalf("Title() panicked: %v", err)
		}
	}()
	r.Title("Test")
}

type testCommand struct {
	name string
}

func (c *testCommand) Name() string        { return c.name }
func (c *testCommand) Description() string  { return "test command" }
func (c *testCommand) Configure(cfg *CommandConfig) {}
func (c *testCommand) Handle(ctx *Context) error { return nil }
