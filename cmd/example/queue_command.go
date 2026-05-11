package main

import (
	"fmt"
	"go-consolekit/console"
	"os"
	"path/filepath"
	"time"
)

type QueueDemoCommand struct{}

func (c *QueueDemoCommand) Name() string {
	return "queue:demo"
}

func (c *QueueDemoCommand) Description() string {
	return "Demonstrate file-based queue worker"
}

func (c *QueueDemoCommand) Configure(config *console.CommandConfig) {}

func (c *QueueDemoCommand) Handle(ctx *console.Context) error {
	dir := filepath.Join(os.TempDir(), "consolekit-queue-demo")
	defer os.RemoveAll(dir)

	q := console.NewQueue(dir).
		Name("demo").
		MaxAttempts(3).
		RetryDelay(1 * time.Second).
		Concurrency(5).
		Logger(func(format string, args ...any) {
			ctx.Line(fmt.Sprintf(format, args...))
		})

	q.Register("send-email", func(tc *console.TaskContext) error {
		var payload struct {
			To   string
			Body string
		}
		tc.Bind(&payload)
		ctx.Line(fmt.Sprintf("Sending email to %s: %s", payload.To, payload.Body))
		return nil
	})

	q.Register("process-payment", func(tc *console.TaskContext) error {
		var payload struct {
			OrderID int
			Amount  float64
		}
		tc.Bind(&payload)
		ctx.Line(fmt.Sprintf("Processing payment $%.2f for order %d", payload.Amount, payload.OrderID))
		return nil
	})

	ctx.Output().Success("Pushing sync email task")
	q.Push("send-email", struct {
		To   string
		Body string
	}{To: "user@example.com", Body: "Welcome!"})

	ctx.Output().Success("Pushing payment task with delay")
	q.Push("process-payment", struct {
		OrderID int
		Amount  float64
	}{OrderID: 1001, Amount: 49.99}, console.WithDelay(3*time.Second))

	ctx.Line(fmt.Sprintf("Pending tasks: %d", q.Count()))

	done := make(chan struct{})
	go func() {
		q.Work()
		close(done)
	}()
	time.Sleep(500 * time.Millisecond)
	q.Stop()
	<-done

	completed, _ := q.Completed()
	ctx.Line(fmt.Sprintf("Completed: %d, Failed: %d",
		len(completed),
		lenOrZero(q.Failed())))

	return nil
}

func lenOrZero(tasks []*console.Task, err error) int {
	if err != nil {
		return 0
	}
	return len(tasks)
}
