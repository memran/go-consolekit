package main

import (
	"go-consolekit/console"
	"time"
)

type ProgressDemoCommand struct{}

func (c *ProgressDemoCommand) Name() string {
	return "progress:demo"
}

func (c *ProgressDemoCommand) Description() string {
	return "Demonstrate progress bar"
}

func (c *ProgressDemoCommand) Configure(config *console.CommandConfig) {}

func (c *ProgressDemoCommand) Handle(ctx *console.Context) error {
	ctx.Output().Progress("Working", 100).Run(func(p *console.Progress) {
		for i := 0; i < 100; i++ {
			p.Advance()
			time.Sleep(20 * time.Millisecond)
		}
	})
	ctx.Success("Done!")
	return nil
}

type TableDemoCommand struct{}

func (c *TableDemoCommand) Name() string {
	return "table:demo"
}

func (c *TableDemoCommand) Description() string {
	return "Demonstrate table rendering"
}

func (c *TableDemoCommand) Configure(config *console.CommandConfig) {}

func (c *TableDemoCommand) Handle(ctx *console.Context) error {
	ctx.Output().
		Table().
		Headers("Name", "Role").
		Row("Emran", "Admin").
		Row("Marwa", "User").
		Row("Alex", "Editor").
		Render()
	return nil
}

type InputDemoCommand struct{}

func (c *InputDemoCommand) Name() string {
	return "input:demo"
}

func (c *InputDemoCommand) Description() string {
	return "Demonstrate interactive input"
}

func (c *InputDemoCommand) Configure(config *console.CommandConfig) {}

func (c *InputDemoCommand) Handle(ctx *console.Context) error {
	name := ctx.Input().
		Ask("Project name").
		Default("my-app").
		Run()

	db := ctx.Input().
		Select("Database").
		Options("SQLite", "PostgreSQL", "MySQL").
		Default("SQLite").
		Run()

	confirm := ctx.Input().
		Confirm("Continue?").
		Default(true).
		Run()

	v := console.NewValidator().
		Data(map[string]any{"name": name}).
		Rule("name", console.Required(), console.MinLen(2))

	if v.Fails() {
		for _, errs := range v.Errors() {
			for _, err := range errs {
				ctx.Error(err)
			}
		}
		return nil
	}

	ctx.Success("Name: " + name)
	ctx.Info("Database: " + db)
	if confirm {
		ctx.Line("You chose to continue")
	} else {
		ctx.Line("You declined")
	}
	return nil
}
