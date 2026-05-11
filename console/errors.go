package console

import "fmt"

type CommandNotFoundError struct {
	Name string
}

func (e *CommandNotFoundError) Error() string {
	return fmt.Sprintf("command not found: %s", e.Name)
}

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error: %s", e.Message)
}

type MissingArgumentError struct {
	ArgumentName string
}

func (e *MissingArgumentError) Error() string {
	return fmt.Sprintf("missing required argument: %s", e.ArgumentName)
}

type OptionError struct {
	OptionName string
	Message    string
}

func (e *OptionError) Error() string {
	return fmt.Sprintf("option --%s: %s", e.OptionName, e.Message)
}
