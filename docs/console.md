# Console Framework

Core types for building CLI applications: `App`, `CommandBuilder`, `Registry`, `Context`, `Renderer`, and errors.

## App

`App` is the main application entry point. Created via `New(name)`.

### Methods

```go
func New(name string) *App
func (a *App) Version(version string) *App
func (a *App) Description(description string) *App
func (a *App) Register(commands ...Command) *App
func (a *App) Command(name string) *CommandBuilder
func (a *App) EnableDaemon() *App
func (a *App) PIDFile(path string) *App
func (a *App) LogFile(path string) *App
func (a *App) Status() (int, error)
func (a *App) Stop() error
func (a *App) Restart() error
func (a *App) Run() error
```

### Example

```go
app := console.New("tool").
    Version("1.0.0").
    Description("A CLI tool")

// Using CommandBuilder fluent API
app.Command("user:create").
    Description("Create a user").
    Argument("email").Required().
    Option("role").Default("user").
    Option("admin").Shortcut("a").
    Handle(func(ctx *console.Context) error {
        email := ctx.Arg("email")
        role := ctx.Option("role")
        ctx.Success("Created user " + email)
        return nil
    })

// Using the Command interface
type GreetCommand struct{}

func (c *GreetCommand) Name() string          { return "greet" }
func (c *GreetCommand) Description() string    { return "Say hello" }
func (c *GreetCommand) Configure(cfg *console.CommandConfig) {
    cfg.Argument("name").Required()
    cfg.Option("lang").Default("en")
}
func (c *GreetCommand) Handle(ctx *console.Context) error {
    ctx.Info("Hello, " + ctx.Arg("name"))
    return nil
}

app.Register(&GreetCommand{})
app.Run()
```

## CommandBuilder

Fluent builder for defining commands inline.

```go
func (b *CommandBuilder) Description(desc string) *CommandBuilder
func (b *CommandBuilder) Argument(name string) *ArgumentConfig
func (b *CommandBuilder) Option(name string) *OptionConfig
func (b *CommandBuilder) Handle(handler func(*Context) error) *CommandBuilder
```

## Command Interface

```go
type Command interface {
    Name() string
    Description() string
    Configure(*CommandConfig)
    Handle(*Context) error
}
```

## CommandConfig

```go
func NewCommandConfig(name string) *CommandConfig
func (c *CommandConfig) Argument(name string) *ArgumentConfig
func (c *CommandConfig) Option(name string) *OptionConfig
```

## ArgumentConfig

```go
func (a *ArgumentConfig) Required() *ArgumentConfig
func (a *ArgumentConfig) Default(value string) *ArgumentConfig
func (a *ArgumentConfig) Description(text string) *CommandBuilder
```

## OptionConfig

```go
func (o *OptionConfig) Shortcut(short string) *OptionConfig
func (o *OptionConfig) Required() *OptionConfig
func (o *OptionConfig) Default(value string) *OptionConfig
func (o *OptionConfig) Description(text string) *CommandBuilder
```

## Registry

Internal command registry. Commands are namespaced with `:` separators.

```go
func NewRegistry() *Registry
func (r *Registry) Add(command Command)
func (r *Registry) Find(name string) (Command, error)
func (r *Registry) All() []Command
```

`Add` panics on duplicate command names.

## Context

`Context` provides access to arguments, options, input/output, and lifecycle hooks.

```go
func (c *Context) Arg(name string) string
func (c *Context) Option(name string) string
func (c *Context) Input() *Input
func (c *Context) Output() *Output
func (c *Context) Line(text string)
func (c *Context) Info(text string)
func (c *Context) Success(text string)
func (c *Context) Warning(text string)
func (c *Context) Error(text string)
func (c *Context) Title(text string)
func (c *Context) Done() <-chan struct{}
func (c *Context) IsCancelled() bool
func (c *Context) OnShutdown(fn func() error)
func (c *Context) Shutdown(message string)
func (c *Context) IsDaemon() bool
func (c *Context) PID() int
func (c *Context) WritePID(path string) error
func (c *Context) StopDaemon() error
func (c *Context) RunWorker(fn func() error)
```

### Example

```go
app.Command("process").
    Argument("file").Required().
    Handle(func(ctx *console.Context) error {
        file := ctx.Arg("file")
        ctx.Info("Processing " + file)

        content, err := console.File(file).Read()
        if err != nil {
            ctx.Error("Cannot read file")
            return err
        }

        ctx.OnShutdown(func() error {
            ctx.Line("Cleaning up...")
            return nil
        })

        ctx.Success("Done")
        return nil
    })
```

## Renderer

```go
type Renderer interface {
    Line(text string)
    Info(text string)
    Success(text string)
    Warning(text string)
    Error(text string)
    Title(text string)
}
```

- `CLIRenderer` (default) - Color output via `fatih/color`
- `TUIRenderer` - Plain text prefix output

## TextBuilder

```go
func (o *Output) Text(text string) *TextBuilder
func (t *TextBuilder) Green() *TextBuilder
func (t *TextBuilder) Red() *TextBuilder
func (t *TextBuilder) Yellow() *TextBuilder
func (t *TextBuilder) Cyan() *TextBuilder
func (t *TextBuilder) Bold() *TextBuilder
func (t *TextBuilder) Prefix(prefix string) *TextBuilder
func (t *TextBuilder) Print()
```

## Errors

```go
type CommandNotFoundError struct{ Name string }
type ValidationError struct{ Message string }
type MissingArgumentError struct{ ArgumentName string }
type OptionError struct{ OptionName, Message string }
```

## Output

```go
func NewOutput(renderer Renderer) *Output
func (o *Output) Line(text string)
func (o *Output) Info(text string)
func (o *Output) Success(text string)
func (o *Output) Warning(text string)
func (o *Output) Error(text string)
func (o *Output) Title(text string)
func (o *Output) Text(text string) *TextBuilder
func (o *Output) Progress(description string, total int) *Progress
func (o *Output) Spinner(description string) *Spinner
func (o *Output) Table() *TableBuilder
```

## Input

```go
func NewInput() *Input
func (i *Input) Ask(prompt string) *AskBuilder
func (i *Input) Confirm(prompt string) *ConfirmBuilder
func (i *Input) Select(prompt string) *SelectBuilder
func (i *Input) Secret(prompt string) *SecretBuilder
```

### Input Builders

```go
// AskBuilder
func (a *AskBuilder) Required() *AskBuilder
func (a *AskBuilder) Default(value string) *AskBuilder
func (a *AskBuilder) Run() string

// ConfirmBuilder
func (c *ConfirmBuilder) Default(value bool) *ConfirmBuilder
func (c *ConfirmBuilder) Run() bool

// SelectBuilder
func (s *SelectBuilder) Options(opts ...string) *SelectBuilder
func (s *SelectBuilder) Default(value string) *SelectBuilder
func (s *SelectBuilder) Run() string

// SecretBuilder
func (s *SecretBuilder) Required() *SecretBuilder
func (s *SecretBuilder) Run() string
```

### Input Example

```go
ctx.Input().Ask("Enter name").
    Required().
    Default("guest").
    Run()

confirmed := ctx.Input().Confirm("Continue?").
    Default(true).
    Run()

role := ctx.Input().Select("Choose role").
    Options("admin", "user", "viewer").
    Run()

password := ctx.Input().Secret("Password").
    Required().
    Run()
```
