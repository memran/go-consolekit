# ConsoleKit

A Go framework for building beautiful console applications with minimal boilerplate. Symfony Console-style command registry, Cobra adapter, fluent DX API, and renderer abstraction.

## Installation

```bash
go get github.com/memran/go-consolekit
```

## Quick Start

```go
package main

import (
    "go-consolekit/console"
)

func main() {
    app := console.New("myapp").
        Version("1.0.0").
        Description("My console application")

    app.Command("hello").
        Description("Say hello").
        Argument("name").
            Required().
            Description("Who to greet").
        Handle(func(ctx *console.Context) error {
            ctx.Success("Hello, " + ctx.Arg("name") + "!")
            return nil
        })

    app.Run()
}
```

## Class-Based Commands

```go
type InstallCommand struct{}

func (c *InstallCommand) Name() string           { return "install" }
func (c *InstallCommand) Description() string    { return "Install a new project" }

func (c *InstallCommand) Configure(config *console.CommandConfig) {
    config.Argument("name").Required().Description("Project name")
    config.Option("db").Shortcut("d").Default("sqlite").Description("Database driver")
}

func (c *InstallCommand) Handle(ctx *console.Context) error {
    ctx.Title("Installer")
    ctx.Success("Project: " + ctx.Arg("name"))
    ctx.Info("Database: " + ctx.Option("db"))
    return nil
}

func main() {
    app := console.New("app")
    app.Register(&InstallCommand{})
    app.Run()
}
```

## Fluent Command Builder

```go
app.Command("greet").
    Description("Greet someone").
    Argument("name").
        Required().
        Description("Name to greet").
    Option("lang").
        Shortcut("l").
        Default("en").
        Description("Language").
    Handle(func(ctx *console.Context) error {
        ctx.Success("Hello, " + ctx.Arg("name"))
        return nil
    })
```

## Interactive Input

```go
name := ctx.Input().
    Ask("Project name").
    Required().
    Default("my-app").
    Run()

ok := ctx.Input().
    Confirm("Continue?").
    Default(true).
    Run()

db := ctx.Input().
    Select("Database").
    Options("SQLite", "PostgreSQL", "MySQL").
    Default("SQLite").
    Run()

password := ctx.Input().
    Secret("Password").
    Required().
    Run()
```

## Styled Output

```go
ctx.Output().Line("Text")
ctx.Output().Info("Info message")
ctx.Output().Success("Success message")
ctx.Output().Warning("Warning message")
ctx.Output().Error("Error message")
ctx.Output().Title("Section Title")

ctx.Output().
    Text("Deployment complete").
    Green().
    Bold().
    Prefix("OK").
    Print()
```

## Progress Bar

```go
ctx.Output().
    Progress("Installing", 100).
    Run(func(p *console.Progress) {
        for i := 0; i < 100; i++ {
            p.Advance()
        }
    })
```

## Table Rendering

```go
ctx.Output().
    Table().
    Headers("Name", "Role").
    Row("Emran", "Admin").
    Row("Marwa", "User").
    Render()
```

## Running Examples

```bash
go run ./cmd/example hello Emran
go run ./cmd/example install blog --db postgres
go run ./cmd/example make model User --table users
go run ./cmd/example table demo
go run ./cmd/example progress demo
go run ./cmd/example input demo
```

## Architecture

```
consolekit/
├── cmd/example/main.go   # Example commands
├── console/
│   ├── app.go            # App entry point, Cobra adapter
│   ├── command.go        # Command interface, CommandBuilder
│   ├── registry.go       # Command registry
│   ├── config.go         # CommandConfig, ArgumentConfig, OptionConfig
│   ├── context.go        # Execution context
│   ├── input.go          # Fluent input API (Ask, Confirm, Select, Secret)
│   ├── output.go         # Output rendering, text builder
│   ├── renderer.go       # Renderer interface, CLIRenderer, TUIRenderer
│   ├── progress.go       # Progress bar
│   ├── spinner.go        # Spinner
│   ├── table.go          # Table builder
│   ├── errors.go         # Error types
│   └── app_test.go       # Unit tests
├── go.mod
└── README.md
```

## Roadmap

- [x] CLI renderer
- [ ] Bubble Tea-based TUI renderer
- [ ] Command autocompletion
- [ ] Event system
- [ ] Plugin support
