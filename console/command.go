package console

type Command interface {
	Name() string
	Description() string
	Configure(*CommandConfig)
	Handle(*Context) error
}

type CommandBuilder struct {
	app     *App
	name    string
	desc    string
	config  *CommandConfig
	handler func(*Context) error
}

func newCommandBuilder(app *App, name string) *CommandBuilder {
	return &CommandBuilder{
		app:    app,
		name:   name,
		config: NewCommandConfig(name),
	}
}

func (b *CommandBuilder) Description(desc string) *CommandBuilder {
	b.desc = desc
	return b
}

func (b *CommandBuilder) Argument(name string) *ArgumentConfig {
	arg := b.config.Argument(name)
	arg.builder = b
	return arg
}

func (b *CommandBuilder) Option(name string) *OptionConfig {
	opt := b.config.Option(name)
	opt.builder = b
	return opt
}

func (b *CommandBuilder) Handle(handler func(*Context) error) *CommandBuilder {
	b.handler = handler
	b.app.registerBuilder(b)
	return b
}

type builderCommand struct {
	name        string
	description string
	config      *CommandConfig
	handler     func(*Context) error
}

func (c *builderCommand) Name() string             { return c.name }
func (c *builderCommand) Description() string       { return c.description }
func (c *builderCommand) Configure(cfg *CommandConfig) {
	cfg.Name = c.config.Name
	cfg.Description = c.config.Description
	for _, arg := range c.config.Arguments {
		a := &ArgumentConfig{Name: arg.Name, description: arg.description, required: arg.required, defaultVal: arg.defaultVal, parent: cfg}
		cfg.Arguments = append(cfg.Arguments, a)
		cfg.argsMap[arg.Name] = a
	}
	for _, opt := range c.config.Options {
		o := &OptionConfig{Name: opt.Name, description: opt.description, shortcut: opt.shortcut, required: opt.required, defaultVal: opt.defaultVal, parent: cfg}
		cfg.Options = append(cfg.Options, o)
		cfg.optsMap[opt.Name] = o
	}
}
func (c *builderCommand) Handle(ctx *Context) error { return c.handler(ctx) }
