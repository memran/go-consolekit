package console

type ArgumentConfig struct {
	Name        string
	description string
	required    bool
	defaultVal  string
	parent      *CommandConfig
	builder     *CommandBuilder
}

func (a *ArgumentConfig) Required() *ArgumentConfig {
	a.required = true
	return a
}

func (a *ArgumentConfig) Default(value string) *ArgumentConfig {
	a.defaultVal = value
	return a
}

func (a *ArgumentConfig) Description(text string) *CommandBuilder {
	a.description = text
	return a.builder
}

type OptionConfig struct {
	Name        string
	description string
	shortcut    string
	required    bool
	defaultVal  string
	parent      *CommandConfig
	builder     *CommandBuilder
}

func (o *OptionConfig) Shortcut(short string) *OptionConfig {
	o.shortcut = short
	return o
}

func (o *OptionConfig) Required() *OptionConfig {
	o.required = true
	return o
}

func (o *OptionConfig) Default(value string) *OptionConfig {
	o.defaultVal = value
	return o
}

func (o *OptionConfig) Description(text string) *CommandBuilder {
	o.description = text
	return o.builder
}

type CommandConfig struct {
	Name        string
	Description string
	Arguments   []*ArgumentConfig
	Options     []*OptionConfig
	argsMap     map[string]*ArgumentConfig
	optsMap     map[string]*OptionConfig
}

func NewCommandConfig(name string) *CommandConfig {
	return &CommandConfig{
		Name:    name,
		argsMap: make(map[string]*ArgumentConfig),
		optsMap: make(map[string]*OptionConfig),
	}
}

func (c *CommandConfig) Argument(name string) *ArgumentConfig {
	if _, exists := c.argsMap[name]; exists {
		panic("duplicate argument: " + name)
	}
	arg := &ArgumentConfig{Name: name, parent: c}
	c.Arguments = append(c.Arguments, arg)
	c.argsMap[name] = arg
	return arg
}

func (c *CommandConfig) Option(name string) *OptionConfig {
	if _, exists := c.optsMap[name]; exists {
		panic("duplicate option: " + name)
	}
	opt := &OptionConfig{Name: name, parent: c}
	c.Options = append(c.Options, opt)
	c.optsMap[name] = opt
	return opt
}
