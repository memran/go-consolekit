package main

import (
	"go-consolekit/console"
	"time"
)

type WorkerDemoCommand struct{}

func (c *WorkerDemoCommand) Name() string {
	return "worker:demo"
}

func (c *WorkerDemoCommand) Description() string {
	return "Demonstrate worker loop with daemon support"
}

func (c *WorkerDemoCommand) Configure(config *console.CommandConfig) {}

func (c *WorkerDemoCommand) Handle(ctx *console.Context) error {
	ctx.Output().Info("Worker started. PID: " + itoa(ctx.PID()))
	ctx.Output().Info("Daemon: " + boolStr(ctx.IsDaemon()))

	ctx.OnShutdown(func() error {
		ctx.Output().Warning("Closing resources...")
		return nil
	})

	ctx.RunWorker(func() error {
		ctx.Output().Info("Processing job...")
		time.Sleep(time.Second)
		return nil
	})

	ctx.Output().Warning("Worker stopped gracefully.")
	return nil
}

type QueueWorkCommand struct{}

func (c *QueueWorkCommand) Name() string {
	return "queue:work"
}

func (c *QueueWorkCommand) Description() string {
	return "Process queue jobs (daemon-friendly)"
}

func (c *QueueWorkCommand) Configure(config *console.CommandConfig) {}

func (c *QueueWorkCommand) Handle(ctx *console.Context) error {
	ctx.Output().Info("Queue worker started. PID: " + itoa(ctx.PID()))
	ctx.Output().Info("Daemon: " + boolStr(ctx.IsDaemon()))

	ctx.OnShutdown(func() error {
		ctx.Output().Warning("Shutting down queue worker...")
		return nil
	})

	for {
		select {
		case <-ctx.Done():
			ctx.Output().Warning("Queue worker stopped.")
			return nil
		default:
			ctx.Output().Success("Processing queue job...")
			time.Sleep(2 * time.Second)
		}
	}
}

type WorkerStopCommand struct{}

func (c *WorkerStopCommand) Name() string {
	return "worker:stop"
}

func (c *WorkerStopCommand) Description() string {
	return "Stop the worker daemon"
}

func (c *WorkerStopCommand) Configure(config *console.CommandConfig) {}

func (c *WorkerStopCommand) Handle(ctx *console.Context) error {
	pid, err := app.Status()
	if err != nil {
		ctx.Warning("No worker daemon is running.")
		return nil
	}
	ctx.Success("Stopping worker (PID " + itoa(pid) + ")...")
	app.Stop()
	ctx.Success("Worker stopped.")
	return nil
}

type WorkerStatusCommand struct{}

func (c *WorkerStatusCommand) Name() string {
	return "worker:status"
}

func (c *WorkerStatusCommand) Description() string {
	return "Check worker daemon status"
}

func (c *WorkerStatusCommand) Configure(config *console.CommandConfig) {}

func (c *WorkerStatusCommand) Handle(ctx *console.Context) error {
	pid, err := app.Status()
	if err != nil {
		ctx.Warning("Worker daemon is not running.")
		return nil
	}
	ctx.Success("Worker daemon is running (PID " + itoa(pid) + ")")
	return nil
}

type WorkerRestartCommand struct{}

func (c *WorkerRestartCommand) Name() string {
	return "worker:restart"
}

func (c *WorkerRestartCommand) Description() string {
	return "Restart the worker daemon"
}

func (c *WorkerRestartCommand) Configure(config *console.CommandConfig) {}

func (c *WorkerRestartCommand) Handle(ctx *console.Context) error {
	pid, err := app.Status()
	if err == nil {
		ctx.Success("Stopping worker (PID " + itoa(pid) + ")...")
		app.Stop()
	}
	if err := app.Restart(); err != nil {
		ctx.Error("Failed to restart: " + err.Error())
		return nil
	}
	newPid, _ := app.Status()
	if newPid != 0 {
		ctx.Success("Worker restarted (PID " + itoa(newPid) + ")")
	} else {
		ctx.Success("Worker restarted")
	}
	return nil
}
