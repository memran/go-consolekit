# Logger

Structured logger with log levels, channels, context, and key-value pairs.

## Constructor

```go
func NewLogger() *Logger
```

## Log Levels

```go
const (
    DebugLevel LogLevel = iota
    InfoLevel
    NoticeLevel
    WarningLevel
    ErrorLevel
    CriticalLevel
    AlertLevel
    EmergencyLevel
)
```

## Level Methods

```go
func (l *Logger) Debug(msg string, kv ...interface{})
func (l *Logger) Info(msg string, kv ...interface{})
func (l *Logger) Notice(msg string, kv ...interface{})
func (l *Logger) Warning(msg string, kv ...interface{})
func (l *Logger) Error(msg string, kv ...interface{})
func (l *Logger) Critical(msg string, kv ...interface{})
func (l *Logger) Alert(msg string, kv ...interface{})
func (l *Logger) Emergency(msg string, kv ...interface{})
```

## Context Methods

```go
func (l *Logger) With(key string, value interface{}) *Logger
func (l *Logger) WithMap(data map[string]interface{}) *Logger
func (l *Logger) WithWriter(writer io.Writer) *Logger
func (l *Logger) WithName(name string) *Logger
```

`With` and `WithMap` return a new `Logger` with added context. Original is not modified.

## Channel Methods

```go
func (l *Logger) Channel(name string) *Logger
func (l *Logger) Stack(names ...string) *Logger
func (l *Logger) AddChannel(name string, writer io.Writer) *Logger
```

- `Channel` returns a logger writing to the named channel's writers (or default if not found)
- `Stack` combines writers from multiple channels
- `AddChannel` registers a writer to a named channel

## Generic Log

```go
func (l *Logger) Log(level LogLevel, msg string, kv ...interface{})
```

## Output Format

```
[2025-05-11 14:30:00] local.INFO: message {"key":"value"}
```

Format: `[timestamp] name.LEVEL: message {"context":"json"}`

## Package-Level Functions

Default logger writes to `os.Stdout`.

```go
func Debug(msg string, kv ...interface{})
func Info(msg string, kv ...interface{})
func Notice(msg string, kv ...interface{})
func Warning(msg string, kv ...interface{})
func Error(msg string, kv ...interface{})
func Critical(msg string, kv ...interface{})
func Alert(msg string, kv ...interface{})
func Emergency(msg string, kv ...interface{})
```

## Examples

```go
log := console.NewLogger()

// Basic logging
log.Info("Server started")
log.Debug("Connecting to database")
log.Error("Connection failed")

// With key-value context
log.Info("User logged in", "user_id", 42, "ip", "192.168.1.1")

// With context
log = log.With("service", "api").
    WithMap(map[string]interface{}{
        "version": "1.0.0",
        "env":     "production",
    })

log.Info("Request started", "method", "GET", "path", "/users")

// Named logger
apiLog := log.WithName("api")
apiLog.Info("API endpoint called") // [timestamp] api.INFO: ...

// Custom writer
var buf bytes.Buffer
fileLog := log.WithWriter(&buf)
fileLog.Warning("Written to buffer")

// Channels
log.AddChannel("errors", os.Stderr)
log.AddChannel("audit", auditFile)

errorLog := log.Channel("errors")
errorLog.Error("Something went wrong")

// Stack channels
stackedLog := log.Stack("errors", "audit")
stackedLog.Warning("Important event") // writes to both channels

// Package-level functions
console.Info("Application started")
console.Warning("Low disk space", "percent", 85)
console.Error("Fatal error", "code", 500)

// With specific level
log.Log(console.NoticeLevel, "Maintenance scheduled", "time", "02:00")
```
