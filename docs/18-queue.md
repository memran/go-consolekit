# Queue & Worker

The `Queue` provides a file-based task queue with a concurrent worker, retry with backoff, delayed tasks, and fluent API.

## Quick Start

```go
q := console.NewQueue("storage/queues").
    Name("emails").
    MaxAttempts(3).
    RetryDelay(30 * time.Second).
    Concurrency(5)

q.Register("send-email", func(ctx *console.TaskContext) error {
    var payload struct { To string; Body string }
    ctx.Bind(&payload)
    return sendEmail(payload.To, payload.Body)
})

task, _ := q.Push("send-email", struct {
    To   string
    Body string
}{To: "user@example.com", Body: "Welcome!"})

q.Work() // blocks, processes tasks
```

## Async Worker

```go
q.Run()            // non-blocking background goroutine
// ...
q.Stop()           // graceful shutdown
q.IsRunning()      // bool
```

## Delayed Tasks

```go
q.Push("send-email", payload, console.WithDelay(5*time.Minute))
```

## Retry & Failure

```go
q.Push("flaky-task", data,
    console.WithMaxAttempts(5),
    console.WithRetryDelay(10*time.Second),
)
```

## Queue Management

```go
q.Pending()          // []*Task
q.Failed()           // []*Task
q.Completed()        // []*Task
q.Count()            // int (pending count)
q.Retry("task-id")   // move failed back to pending
q.Remove("task-id")  // delete task
q.Flush()            // clear all tasks
```

## Logging

```go
q.Logger(func(format string, args ...any) {
    log.Printf(format, args...)
})
```

## Storage Layout

```
storage/queues/
  <name>/
    pending/<uuid>.json
    processing/<uuid>.json
    failed/<uuid>.json
    completed/<uuid>.json
```

## Task Lifecycle

1. **Push** → written to `pending/`
2. **Worker picks up** → moved to `processing/`, handler called
3. **Success** → moved to `completed/`
4. **Failure + retries left** → moved to `pending/` with `scheduled_at`
5. **Failure + no retries** → moved to `failed/`

## Method Reference

| Method | Description |
|--------|-------------|
| `NewQueue(basePath)` | Create queue with file storage root |
| `Name(name)` | Set queue name (directory suffix) |
| `MaxAttempts(n)` | Max retry attempts (default 3) |
| `RetryDelay(d)` | Delay between retries (default 30s) |
| `Concurrency(n)` | Max concurrent workers (default 1) |
| `Logger(fn)` | Attach logger |
| `Register(name, handler)` | Register task handler |
| `Push(name, payload, opts...)` | Push a task |
| `Work()` | Start blocking worker |
| `Run()` | Start non-blocking worker |
| `Stop()` | Stop worker |
| `IsRunning()` | Check if worker is running |
| `Pending()` | List pending tasks |
| `Failed()` | List failed tasks |
| `Completed()` | List completed tasks |
| `Count()` | Pending task count |
| `Retry(id)` | Retry a failed task |
| `Remove(id)` | Delete a task |
| `Flush()` | Delete all tasks |

## Task Options

| Option | Description |
|--------|-------------|
| `WithDelay(d)` | Schedule task for future execution |
| `WithMaxAttempts(n)` | Override max retry attempts |
| `WithRetryDelay(d)` | Override retry delay |
