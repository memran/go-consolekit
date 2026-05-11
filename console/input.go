package console

import (
	"github.com/AlecAivazis/survey/v2"
)

type Input struct{}

func NewInput() *Input {
	return &Input{}
}

func (i *Input) Ask(prompt string) *AskBuilder {
	return &AskBuilder{
		prompt: prompt,
	}
}

func (i *Input) Confirm(prompt string) *ConfirmBuilder {
	return &ConfirmBuilder{
		prompt: prompt,
	}
}

func (i *Input) Select(prompt string) *SelectBuilder {
	return &SelectBuilder{
		prompt: prompt,
	}
}

func (i *Input) Secret(prompt string) *SecretBuilder {
	return &SecretBuilder{
		prompt: prompt,
	}
}

type AskBuilder struct {
	prompt   string
	required bool
	default_ string
}

func (a *AskBuilder) Required() *AskBuilder {
	a.required = true
	return a
}

func (a *AskBuilder) Default(value string) *AskBuilder {
	a.default_ = value
	return a
}

func (a *AskBuilder) Run() string {
	var result string
	prompt := &survey.Input{
		Message: a.prompt,
		Default: a.default_,
	}
	survey.AskOne(prompt, &result, survey.WithValidator(survey.Required))
	return result
}

type ConfirmBuilder struct {
	prompt   string
	default_ bool
}

func (c *ConfirmBuilder) Default(value bool) *ConfirmBuilder {
	c.default_ = value
	return c
}

func (c *ConfirmBuilder) Run() bool {
	var result bool
	prompt := &survey.Confirm{
		Message: c.prompt,
		Default: c.default_,
	}
	survey.AskOne(prompt, &result)
	return result
}

type SelectBuilder struct {
	prompt   string
	options  []string
	default_ string
}

func (s *SelectBuilder) Options(opts ...string) *SelectBuilder {
	s.options = opts
	return s
}

func (s *SelectBuilder) Default(value string) *SelectBuilder {
	s.default_ = value
	return s
}

func (s *SelectBuilder) Run() string {
	var result string
	prompt := &survey.Select{
		Message: s.prompt,
		Options: s.options,
		Default: s.default_,
	}
	survey.AskOne(prompt, &result)
	return result
}

type SecretBuilder struct {
	prompt   string
	required bool
}

func (s *SecretBuilder) Required() *SecretBuilder {
	s.required = true
	return s
}

func (s *SecretBuilder) Run() string {
	var result string
	prompt := &survey.Password{
		Message: s.prompt,
	}
	askOpts := make([]survey.AskOpt, 0)
	if s.required {
		askOpts = append(askOpts, survey.WithValidator(survey.Required))
	}
	survey.AskOne(prompt, &result, askOpts...)
	return result
}
