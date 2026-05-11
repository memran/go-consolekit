# Daemon Mode

Run CLI commands as background daemon processes with lifecycle management.

## Enabling Daemon Mode

```go
app := console.New("myservice").
    Version("1.0.0").
    EnableDaemon().
    PIDFile("/var/run/myservice.pid").
    LogFile("/var/log/myservice.log")
```

Adding `EnableDaemon()` registers `--daemon`, `--pid-file`, and `--log-file` persistent flags.

## Lifecycle

1. User runs `myservice serve --daemon --pid-file /tmp/myservice.pid`
2. Parent process spawns a child with `_CONSOLEKIT_DAEMON_CHILD=1`
3. Parent writes the child PID and exits
4. Child redirects stdout/stderr to the log file
5. Child runs the command handler

## App Methods

```go
func (a *App) EnableDaemon() *App
func (a *App) PIDFile(path string) *App
func (a *App) LogFile(path string) *App
func (a *App) Status() (int, error)
func (a *App) Stop() error
func (a *App) Restart() error
```

### Status

Returns the daemon PID if running, or an error if not.

```go
pid, err := app.Status()
if err != nil {
    fmt.Println("Daemon is not running")
} else {
    fmt.Printf("Running with PID %d\n", pid)
}
```

### Stop

Sends SIGTERM to the daemon process and removes the PID file.

```go
if err := app.Stop(); err != nil {
    fmt.Println("Failed to stop:", err)
}
```

### Restart

Stops the current daemon and spawns a new child process.

```go
if err := app.Restart(); err != nil {
    fmt.Println("Failed to restart:", err)
}
```

## Context Daemon Methods

```go
func (c *Context) IsDaemon() bool
func (c *Context) PID() int
func (c *Context) WritePID(path string) error
func (c *Context) StopDaemon() error
func (c *Context) RunWorker(fn func() error)
```

### IsDaemon

Returns `true` if the process is a daemon child.

```go
if ctx.IsDaemon() {
    ctx.Info("Running as daemon")
}
```

### RunWorker

Runs a function in a loop until shutdown is triggered. Panics in the worker are recovered and logged.

```go
app.Command("serve").
    Handle(func(ctx *console.Context) error {
        ctx.Info("Server started")

        go ctx.RunWorker(func() error {
            // do periodic work
            return nil
        })

        <-ctx.Done()
        ctx.Info("Shutting down")
        return nil
    })
```

## Graceful Shutdown

The framework handles SIGINT/SIGTERM. On first signal:
- Cancels the root context (`ctx.Done()` closes)
- Runs all registered shutdown hooks

On second signal (double Ctrl+C):
- Forces `os.Exit(1)`

```go
ctx.OnShutdown(func() error {
    ctx.Info("Cleaning up resources...")
    return nil
})
```

## Daemon-Specific Context Methods

```go
// Write PID file
ctx.WritePID("/var/run/app.pid")

// Get current PID
pid := ctx.PID()

// Stop daemon from within command
ctx.StopDaemon()
```

## Internal Functions

```go
func IsDaemonChild() bool
func writePID(path string) error
func readPID(path string) (int, error)
func removePID(path string) error
func redirectOutput(path string) error
```

## Complete Example

```go
func main() {
    app := console.New("myservice").
        Version("1.0.0").
        Description("Background service").
        EnableDaemon().
        PIDFile("/tmp/myservice.pid").
        LogFile("/tmp/myservice.log")

    app.Command("serve").
        Handle(func(ctx *console.Context) error {
            ctx.Info("Service started (PID: %d)", ctx.PID())

            ctx.OnShutdown(func() error {
                ctx.Info("Performing cleanup")
                return nil
            })

            go ctx.RunWorker(func() error {
                // background work
                return nil
            })

            <-ctx.Done()
            return nil
        })

    app.Run()
}
```

### Usage

```bash
# Start daemon
myservice serve --daemon

# Check status
myservice status

# Stop daemon
myservice stop

# Restart
myservice restart
```
