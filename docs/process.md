# Process Management

Run, control, and capture external processes with a fluent API.

## Constructors

```go
func NewProcess(name string, args ...string) *Process
func Run(name string, args ...string) *Process
```

`Run` is a shortcut for `NewProcess`.

## Configuration

```go
proc := console.NewProcess("git", "clone", "repo").
    WithWorkingDir("/tmp").
    WithEnv("GIT_TERMINAL_PROMPT", "0").
    WithEnvs(map[string]string{"KEY": "val"}).
    WithInputString("stdin content").
    Timeout(30 * time.Second)
```

| Method | Description |
|--------|-------------|
| `WithWorkingDir(dir)` | Set working directory |
| `WithEnv(key, value)` | Add environment variable |
| `WithEnvs(map)` | Add multiple env vars |
| `WithInput(reader)` | Set stdin from io.Reader |
| `WithInputString(s)` | Set stdin from string |
| `Timeout(duration)` | Kill process after duration |

## Running

### Synchronous

```go
result := console.Run("ls", "-la").Run()
```

### Start / Wait (background)

```go
proc := console.NewProcess("long-task")
proc.Start()
// do other work
result := proc.Wait()
```

### MustRun (panic on failure)

```go
result := console.Run("git", "status").MustRun()
```

## Process Control

```go
proc.Start()
proc.PID()          // 0 before Start, PID after
proc.IsRunning()    // true between Start and Wait
proc.Wait()         // returns *Result
proc.Stop()         // kill the process
proc.Signal(sig)    // send custom signal
```

## Result

```go
result := proc.Run()

result.ExitCode()      // int
result.Output()        // stdout as string
result.Error()         // stderr as string
result.Err()           // underlying error
result.IsSuccessful()  // exit code == 0
result.IsFailed()      // exit code != 0
result.Lines()         // split stdout by newline
result.ErrorLines()    // split stderr by newline
```

## Examples

```go
// Capture command output
result := console.Run("go", "version").Run()
fmt.Println(result.Output())

// Check success
if result.IsSuccessful() {
    fmt.Println("OK:", result.Output())
}

// With timeout
result := console.NewProcess("sleep", "10").
    Timeout(1 * time.Second).
    Run()
// result.Err() contains timeout error

// Stdin pipe
result := console.Run("grep", "error").
    WithInputString("line1\nline with error\nline3\n").
    Run()
fmt.Println(result.Output()) // "line with error"

// Background process
proc := console.NewProcess("server")
proc.Start()
fmt.Printf("PID: %d\n", proc.PID())
time.Sleep(5 * time.Second)
proc.Stop()

// Chaining
result := console.Run("go", "test", "./...").
    WithEnv("CGO_ENABLED", "0").
    Timeout(60 * time.Second).
    Run()
```
