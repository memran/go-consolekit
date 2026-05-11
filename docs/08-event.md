# Event System

The `EventBus` provides a fluent pub/sub system with sync and async dispatch, wildcard matching, priority ordering, and panic recovery.

## Quick Start

```go
bus := console.NewEventBus()

bus.On("order.placed", func(e any) error {
    order := e.(map[string]any)
    fmt.Printf("Processing order %d\n", order["id"])
    return nil
})

bus.Emit("order.placed", map[string]any{"id": 1001})
```

## Sync Dispatch

```go
bus.On("event", handler1)
bus.On("event", handler2)
bus.Emit("event", payload)  // handlers run in order
```

## Async Dispatch

```go
bus.On("task", func(e any) error {
    time.Sleep(100 * time.Millisecond)
    return nil
})

bus.EmitAsync("task", payload)
errs := bus.Flush()  // wait for all, collect errors
```

## Wildcards

```go
bus.On("user.*", func(e any) error {
    // catches: user.registered, user.deleted, user.updated
    return nil
})

bus.On("app.*", func(e any) error { return nil })
// matches: app.user.created, app.user.deleted, app.anything
```

## Priority

```go
bus.On("event", handler, 10)  // high priority (runs first)
bus.On("event", handler, 0)   // low priority (runs last)
```

## Subscribers

```go
type UserSubscriber struct{}
func (s *UserSubscriber) Subscribe(bus *console.EventBus) {
    bus.On("user.registered", s.onRegistered)
    bus.On("user.deleted", s.onDeleted)
}

bus.Subscribe(&UserSubscriber{})
```

## Removing Listeners

```go
bus.Off("event")    // remove all listeners
bus.Off("event", h) // remove specific handler
```

## Checking

```go
bus.HasListeners("order.placed") // bool
```

## Error & Panic Safety

```go
bus.On("failing", func(e any) error {
    return errors.New("something went wrong")
})
bus.On("panicking", func(e any) error {
    panic("unexpected!")
})

errs := bus.Flush()
// errors are collected, panics are recovered
```

## Method Reference

| Method | Description |
|--------|-------------|
| `NewEventBus()` | Create a new event bus |
| `On(event, handler, priority...)` | Register a listener |
| `Off(event, handler...)` | Remove listener(s) |
| `Emit(event, payload)` | Dispatch synchronously |
| `EmitAsync(event, payload)` | Dispatch concurrently |
| `Flush()` | Wait for async completions, return errors |
| `HasListeners(event)` | Check if any listener matches |
| `Subscribe(subscriber)` | Register a Subscriber struct |
