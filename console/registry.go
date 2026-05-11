package console

import "fmt"

type Registry struct {
	commands map[string]Command
}

func NewRegistry() *Registry {
	return &Registry{
		commands: make(map[string]Command),
	}
}

func (r *Registry) Add(command Command) {
	name := command.Name()
	if _, exists := r.commands[name]; exists {
		panic(fmt.Sprintf("command already registered: %s", name))
	}
	r.commands[name] = command
}

func (r *Registry) Find(name string) (Command, error) {
	cmd, exists := r.commands[name]
	if !exists {
		return nil, &CommandNotFoundError{Name: name}
	}
	return cmd, nil
}

func (r *Registry) All() []Command {
	cmds := make([]Command, 0, len(r.commands))
	for _, cmd := range r.commands {
		cmds = append(cmds, cmd)
	}
	return cmds
}
