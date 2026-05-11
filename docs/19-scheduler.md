# Scheduler

The `Scheduler` manages cron-like scheduled tasks with file-based persistence. Schedule entries survive restarts — defined once, persist to disk, auto-load on start.

## Quick Start

```go
q := console.NewQueue("storage/queues").Concurrency(5)
q.Register("send-digest", func(ctx *console.TaskContext) error {
    return nil
})

s := console.NewScheduler("storage/schedules").Queue(q)

s.Every("5m").Call("send-digest", nil)

s.Start() // non-blocking
// ...
s.Stop()
```

## Schedule Expressions

| Expression | Meaning |
|------------|---------|
| `"5m"` | Every 5 minutes |
| `"1h"` | Every hour |
| `"30s"` | Every 30 seconds |
| `"2h30m"` | Every 2.5 hours |
| `@daily` | Once per day |
| `@hourly` | Once per hour |

## Fluent API

```go
s.Every("5m").Call("task-name", payload)
s.Daily().Call("daily-digest", data)
s.Hourly().Call("hourly-check", nil)
```

## Start / Stop

```go
s.Start() // non-blocking background goroutine
s.Stop()

s.Run()   // blocking
```

## Persistence

Schedule entries are saved as JSON files in the `basePath` directory:

```json
storage/schedules/<uuid>.json
{
  "id": "uuid",
  "expression": "5m",
  "task_name": "send-digest",
  "payload": {...},
  "enabled": true,
  "last_run": "2026-05-11T10:00:00Z"
}
```

Entries are loaded automatically on `Start()`/`Run()`. Only enabled entries are loaded.

## Management

```go
s.Entries()          // []*ScheduleEntry
s.Remove("entry-id") // delete schedule
```

## Logging

```go
s.Logger(func(format string, args ...any) {
    log.Printf(format, args...)
})
```

## Method Reference

| Method | Description |
|--------|-------------|
| `NewScheduler(basePath)` | Create scheduler with file storage |
| `Queue(queue)` | Attach queue |
| `Every(dur)` | Create interval schedule (e.g. `"5m"`) |
| `Daily()` | Create daily schedule |
| `Hourly()` | Create hourly schedule |
| `Run()` | Start blocking scheduler loop |
| `Start()` | Start non-blocking scheduler |
| `Stop()` | Stop scheduler |
| `IsRunning()` | Check if running |
| `Entries()` | List all schedule entries |
| `Remove(id)` | Delete a schedule entry |
| `Logger(fn)` | Attach logger |
