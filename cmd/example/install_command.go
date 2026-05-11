package main

import "go-consolekit/console"

type InstallCommand struct{}

func (c *InstallCommand) Name() string {
	return "install"
}

func (c *InstallCommand) Description() string {
	return "Install a new project"
}

func (c *InstallCommand) Configure(config *console.CommandConfig) {
	config.Argument("name").Required().Description("Project name")
	config.Option("db").Shortcut("d").Default("sqlite").Description("Database driver (sqlite, mysql, pgsql)")
}

func (c *InstallCommand) Handle(ctx *console.Context) error {
	v := console.NewValidator().
		Data(map[string]any{"name": ctx.Arg("name"), "db": ctx.Option("db")}).
		Rule("name", console.Required(), console.MinLen(2), console.MaxLen(100)).
		Rule("db", console.Required(), console.In("sqlite", "mysql", "pgsql"))

	if v.Fails() {
		for _, errs := range v.Errors() {
			for _, err := range errs {
				ctx.Error(err)
			}
		}
		return nil
	}

	ctx.Title("ConsoleKit Installer")
	ctx.Success("Project: " + ctx.Arg("name"))
	ctx.Info("Database: " + ctx.Option("db"))
	return nil
}
