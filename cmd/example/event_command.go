package main

import (
	"fmt"
	"go-consolekit/console"
	"sync"
)

type EventDemoCommand struct{}

func (c *EventDemoCommand) Name() string {
	return "event:demo"
}

func (c *EventDemoCommand) Description() string {
	return "Demonstrate async event system"
}

func (c *EventDemoCommand) Configure(config *console.CommandConfig) {}

func (c *EventDemoCommand) Handle(ctx *console.Context) error {
	ctx.Output().Success("--- Sync events ---")
	bus := console.NewEventBus()
	bus.On("order.placed", func(e any) error {
		order := e.(map[string]any)
		ctx.Line(fmt.Sprintf("Processing order %d for %s", order["id"], order["email"]))
		return nil
	})
	bus.On("order.placed", func(e any) error {
		ctx.Line("Sending confirmation email...")
		return nil
	})
	bus.Emit("order.placed", map[string]any{"id": 1001, "email": "user@example.com"})

	ctx.Output().Success("--- Async events ---")
	var mu sync.Mutex
	asyncBus := console.NewEventBus()
	asyncBus.On("task", func(e any) error {
		mu.Lock()
		ctx.Line("Worker 1 processing...")
		mu.Unlock()
		return nil
	})
	asyncBus.On("task", func(e any) error {
		mu.Lock()
		ctx.Line("Worker 2 processing...")
		mu.Unlock()
		return nil
	})
	asyncBus.EmitAsync("task", nil)
	errs := asyncBus.Flush()
	if len(errs) > 0 {
		ctx.Warning(fmt.Sprintf("Errors: %v", errs))
	}

	ctx.Output().Success("--- Wildcard listeners ---")
	wildBus := console.NewEventBus()
	wildBus.On("user.*", func(e any) error {
		ctx.Line(fmt.Sprintf("User event: %v", e))
		return nil
	})
	wildBus.Emit("user.registered", "user-1")
	wildBus.Emit("user.deleted", "user-2")

	ctx.Output().Success("--- Priority ---")
	prioBus := console.NewEventBus()
	prioBus.On("event", func(e any) error {
		ctx.Line("Low priority")
		return nil
	}, 0)
	prioBus.On("event", func(e any) error {
		ctx.Line("High priority")
		return nil
	}, 10)
	prioBus.Emit("event", nil)

	return nil
}
