package main

import (
	"fmt"
	"go-consolekit/console"
	"os"
	"path/filepath"
	"time"
)

type SchedulerDemoCommand struct{}

func (c *SchedulerDemoCommand) Name() string {
	return "scheduler:demo"
}

func (c *SchedulerDemoCommand) Description() string {
	return "Demonstrate scheduler with file-based storage"
}

func (c *SchedulerDemoCommand) Configure(config *console.CommandConfig) {}

func (c *SchedulerDemoCommand) Handle(ctx *console.Context) error {
	qdir := filepath.Join(os.TempDir(), "consolekit-sched-queue")
	sdir := filepath.Join(os.TempDir(), "consolekit-sched-data")
	defer os.RemoveAll(qdir)
	defer os.RemoveAll(sdir)

	q := console.NewQueue(qdir).
		Name("scheduler-demo").
		Concurrency(3).
		Logger(func(format string, args ...any) {
			ctx.Line(fmt.Sprintf("[queue] "+format, args...))
		})

	q.Register("send-digest", func(tc *console.TaskContext) error {
		ctx.Line("Sending daily digest email...")
		return nil
	})

	q.Register("cleanup-logs", func(tc *console.TaskContext) error {
		ctx.Line("Cleaning up old log files...")
		return nil
	})

	s := console.NewScheduler(sdir).
		Queue(q).
		Logger(func(format string, args ...any) {
			ctx.Line(fmt.Sprintf("[scheduler] "+format, args...))
		})

	s.Every("10s").Call("send-digest", nil)
	s.Every("30s").Call("cleanup-logs", nil)

	ctx.Output().Success("Scheduler started (runs for 5 seconds)")

	go func() {
		q.Run()
	}()
	s.Start()

	time.Sleep(5 * time.Second)

	s.Stop()
	q.Stop()

	entries := s.Entries()
	ctx.Line(fmt.Sprintf("\nSchedule entries: %d", len(entries)))
	for _, e := range entries {
		lastRun := "never"
		if e.LastRun != nil {
			lastRun = e.LastRun.Format("15:04:05")
		}
		ctx.Line(fmt.Sprintf("  %s: %s -> %s (last: %s)", e.ID[:8], e.Expression, e.TaskName, lastRun))
	}

	completed, _ := q.Completed()
	ctx.Line(fmt.Sprintf("Tasks completed: %d", len(completed)))

	return nil
}
