package main

import "go-consolekit/console"

type HelloCommand struct{}

func (c *HelloCommand) Name() string {
	return "hello"
}

func (c *HelloCommand) Description() string {
	return "Say hello to someone"
}

func (c *HelloCommand) Configure(config *console.CommandConfig) {
	config.Argument("name").Required().Description("Who to greet")
}

func (c *HelloCommand) Handle(ctx *console.Context) error {
	v := console.NewValidator().
		Data(map[string]any{"name": ctx.Arg("name")}).
		Rule("name", console.Required(), console.MinLen(2), console.MaxLen(50))

	if v.Fails() {
		for _, errs := range v.Errors() {
			for _, err := range errs {
				ctx.Error(err)
			}
		}
		return nil
	}

	ctx.Success("Hello, " + ctx.Arg("name") + "!")
	return nil
}
