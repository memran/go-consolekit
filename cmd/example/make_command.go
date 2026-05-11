package main

import (
	"go-consolekit/console"
	"regexp"
)

type MakeModelCommand struct{}

func (c *MakeModelCommand) Name() string {
	return "make:model"
}

func (c *MakeModelCommand) Description() string {
	return "Create a new model class"
}

func (c *MakeModelCommand) Configure(config *console.CommandConfig) {
	config.Argument("name").Required().Description("Model name (PascalCase)")
	config.Option("table").Shortcut("t").Description("Table name")
}

func (c *MakeModelCommand) Handle(ctx *console.Context) error {
	name := ctx.Arg("name")

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

	if !regexp.MustCompile(`^[A-Z][a-zA-Z0-9]+$`).MatchString(name) {
		ctx.Error("The name must be in PascalCase (e.g. UserProfile).")
		return nil
	}

	ctx.Success("Model created: " + name)
	if table := ctx.Option("table"); table != "" {
		ctx.Info("Table: " + table)
	}
	return nil
}
