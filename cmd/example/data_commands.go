package main

import (
	"bytes"
	"fmt"
	"go-consolekit/console"
)

type CollectionDemoCommand struct{}

func (c *CollectionDemoCommand) Name() string {
	return "collection:demo"
}

func (c *CollectionDemoCommand) Description() string {
	return "Demonstrate fluent collection operations"
}

func (c *CollectionDemoCommand) Configure(config *console.CommandConfig) {}

func (c *CollectionDemoCommand) Handle(ctx *console.Context) error {
	nums := console.Collect([]any{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})

	ctx.Output().Success("--- Filter + Map chain ---")
	result := nums.
		Filter(func(v any) bool { return v.(int)%2 == 0 }).
		Map(func(v any) any { return v.(int) * 10 })
	ctx.Line(result.Implode(", "))

	ctx.Output().Success("--- Reduce sum ---")
	sum := nums.Reduce(func(c, v any) any { return c.(int) + v.(int) }, 0)
	ctx.Line("Sum: " + itoa(sum.(int)))

	ctx.Output().Success("--- Averages ---")
	ctx.Line(fmt.Sprintf("Avg: %.1f", console.Collect([]any{10, 20, 30, 40, 50}).Avg()))

	ctx.Output().Success("--- Pluck names ---")
	users := console.Collect([]any{
		map[string]any{"name": "Alice", "role": "admin"},
		map[string]any{"name": "Bob", "role": "user"},
		map[string]any{"name": "Charlie", "role": "admin"},
	})
	admins := users.Where("role", "admin").Pluck("name")
	ctx.Line("Admins: " + admins.Implode(", "))

	ctx.Output().Success("--- JSON ---")
	json, _ := console.Collect([]any{1, 2, 3}).ToJSON()
	ctx.Line(json)

	ctx.Output().Success("--- Unique + Sort + Reverse ---")
	vals := console.Collect([]any{3, 1, 2, 3, 1, 4, 2}).
		Unique().
		Sort(func(a, b any) bool { return a.(int) < b.(int) }).
		Reverse()
	ctx.Line(vals.Implode(", "))

	return nil
}

type LogDemoCommand struct{}

func (c *LogDemoCommand) Name() string {
	return "log:demo"
}

func (c *LogDemoCommand) Description() string {
	return "Demonstrate fluent logging"
}

func (c *LogDemoCommand) Configure(config *console.CommandConfig) {}

func (c *LogDemoCommand) Handle(ctx *console.Context) error {
	var buf bytes.Buffer
	logger := console.NewLogger().WithWriter(&buf)

	logger.Info("This is an info message")
	logger.Warning("This is a warning")
	logger.Error("This is an error")

	logger.With("user_id", 42).With("role", "admin").Info("User login")

	logger.WithMap(map[string]interface{}{"order_id": 1001, "total": 49.99}).
		Notice("Order placed")

	logger.WithName("api").Debug("Debugging", "endpoint", "/users")

	ctx.Output().Success("--- Logger output ---")
	ctx.Line(buf.String())

	ctx.Output().Success("--- Channel logging ---")
	var chBuf bytes.Buffer
	ch := console.NewLogger().AddChannel("file", &chBuf)
	ch.Channel("file").Info("Channel message")
	ctx.Line(chBuf.String())

	return nil
}
