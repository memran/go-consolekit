package main

import (
	"fmt"
	"go-consolekit/console"
	"time"
)

type ProcessDemoCommand struct{}

func (c *ProcessDemoCommand) Name() string {
	return "process:demo"
}

func (c *ProcessDemoCommand) Description() string {
	return "Demonstrate process management"
}

func (c *ProcessDemoCommand) Configure(config *console.CommandConfig) {}

func (c *ProcessDemoCommand) Handle(ctx *console.Context) error {
	ctx.Output().Success("--- Run and capture ---")
	result := console.Run("go", "version").Run()
	ctx.Line("Output: " + result.Output())
	ctx.Line("Exit code: " + itoa(result.ExitCode()))
	ctx.Line("Success: " + boolStr(result.IsSuccessful()))

	ctx.Output().Success("--- Start and Wait ---")
	proc := console.NewProcess("go", "version")
	if err := proc.Start(); err != nil {
		ctx.Error("Start failed: " + err.Error())
		return nil
	}
	ctx.Line("PID: " + itoa(proc.PID()))
	r := proc.Wait()
	ctx.Line("Result: " + r.Output())

	ctx.Output().Success("--- With env ---")
	result = console.Run("go", "env", "GOPATH").
		WithEnv("GOPATH", "/custom/gopath").
		Run()
	ctx.Line("GOPATH: " + result.Output())

	ctx.Output().Success("--- With timeout ---")
	result = console.NewProcess("go", "version").
		Timeout(5 * time.Second).
		Run()
	ctx.Line("Timed run: " + result.Output())

	ctx.Output().Success("--- Lines ---")
	result = console.Run("go", "version").Run()
	for i, line := range result.Lines() {
		ctx.Line(fmt.Sprintf("line[%d]: %s", i, line))
	}

	return nil
}
