package main

import (
	"fmt"
	"go-consolekit/console"
	"time"
)

type FileDemoCommand struct{}

func (c *FileDemoCommand) Name() string {
	return "file:demo"
}

func (c *FileDemoCommand) Description() string {
	return "Demonstrate file operations"
}

func (c *FileDemoCommand) Configure(config *console.CommandConfig) {}

func (c *FileDemoCommand) Handle(ctx *console.Context) error {
	fo := console.File("storage/demo.txt")
	fo.Write("hello from ConsoleKit\n")
	fo.Append("line 2\n")
	fo.Prepend("line 0\n")

	content, _ := fo.Read()
	ctx.Output().Info("File content:")
	ctx.Output().Line(content)

	lines, _ := fo.Lines()
	ctx.Output().Info("Lines: " + itoa(len(lines)))
	ctx.Output().Info("Exists: " + boolStr(fo.Exists()))
	ctx.Output().Info("Basename: " + fo.Basename())
	ctx.Output().Info("Extension: " + fo.Extension())

	fo.Delete()
	ctx.Output().Success("File cleaned up")
	return nil
}

type HttpDemoCommand struct{}

func (c *HttpDemoCommand) Name() string {
	return "http:demo"
}

func (c *HttpDemoCommand) Description() string {
	return "Demonstrate HTTP client"
}

func (c *HttpDemoCommand) Configure(config *console.CommandConfig) {}

func (c *HttpDemoCommand) Handle(ctx *console.Context) error {
	resp, err := console.Http().
		WithHeader("Accept", "application/json").
		Timeout(5 * time.Second).
		Get("https://httpbin.org/get")
	if err != nil {
		ctx.Warning("HTTP request failed: " + err.Error())
		ctx.Info("(httpbin.org may be unavailable)")
		return nil
	}
	ctx.Success("HTTP " + itoa(resp.StatusCode()))
	ctx.Info("Content-Type: " + resp.Headers()["Content-Type"])
	return nil
}

type ConfigDemoCommand struct{}

func (c *ConfigDemoCommand) Name() string {
	return "config:demo"
}

func (c *ConfigDemoCommand) Description() string {
	return "Demonstrate config management"
}

func (c *ConfigDemoCommand) Configure(config *console.CommandConfig) {}

func (c *ConfigDemoCommand) Handle(ctx *console.Context) error {
	cfg := console.NewConfig().
		Set("app.name", "ConsoleKit").
		Set("app.version", "1.0").
		Set("app.debug", true).
		Set("database.default", "mysql").
		Set("database.connections.mysql.host", "localhost").
		Set("database.connections.mysql.port", 3306)

	ctx.Success("app.name: " + cfg.GetString("app.name"))
	ctx.Success("app.debug: " + boolStr(cfg.GetBool("app.debug")))
	ctx.Info("database.default: " + cfg.GetString("database.default"))
	ctx.Info("mysql.host: " + cfg.GetString("database.connections.mysql.host"))
	ctx.Info("mysql.port: " + itoa(cfg.GetInt("database.connections.mysql.port")))
	ctx.Info("missing.key: " + cfg.GetString("missing.key", "fallback"))
	ctx.Line("Has app.name: " + boolStr(cfg.Has("app.name")))
	return nil
}

type EnvDemoCommand struct{}

func (c *EnvDemoCommand) Name() string {
	return "env:demo"
}

func (c *EnvDemoCommand) Description() string {
	return "Demonstrate env management"
}

func (c *EnvDemoCommand) Configure(config *console.CommandConfig) {}

func (c *EnvDemoCommand) Handle(ctx *console.Context) error {
	env := console.NewEnv().
		Set("APP_NAME", "ConsoleKit").
		Set("APP_DEBUG", "true").
		Set("APP_PORT", "8080").
		Set("APP_RATE", "3.5")

	ctx.Success("APP_NAME: " + env.GetString("APP_NAME"))
	ctx.Success("APP_DEBUG: " + boolStr(env.GetBool("APP_DEBUG")))
	ctx.Info("APP_PORT: " + itoa(env.GetInt("APP_PORT")))
	ctx.Info("APP_RATE: " + fmt.Sprintf("%.1f", env.GetFloat("APP_RATE")))
	ctx.Info("MISSING: " + env.Get("MISSING", "fallback"))
	ctx.Line("Has APP_NAME: " + boolStr(env.Has("APP_NAME")))
	ctx.Line("Has MISSING: " + boolStr(env.Has("MISSING")))
	return nil
}
