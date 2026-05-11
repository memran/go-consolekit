package console

import (
	"github.com/fatih/color"
)

type Output struct {
	renderer Renderer
}

func NewOutput(renderer Renderer) *Output {
	return &Output{renderer: renderer}
}

func (o *Output) Line(text string) {
	o.renderer.Line(text)
}

func (o *Output) Info(text string) {
	o.renderer.Info(text)
}

func (o *Output) Success(text string) {
	o.renderer.Success(text)
}

func (o *Output) Warning(text string) {
	o.renderer.Warning(text)
}

func (o *Output) Error(text string) {
	o.renderer.Error(text)
}

func (o *Output) Title(text string) {
	o.renderer.Title(text)
}

func (o *Output) Text(text string) *TextBuilder {
	return &TextBuilder{
		text: text,
	}
}

func (o *Output) Progress(description string, total int) *Progress {
	return NewProgress(description, total)
}

func (o *Output) Spinner(description string) *Spinner {
	return NewSpinner(description)
}

func (o *Output) Table() *TableBuilder {
	return NewTableBuilder()
}

type TextBuilder struct {
	text   string
	green  bool
	red    bool
	yellow bool
	cyan   bool
	bold   bool
	prefix string
}

func (t *TextBuilder) Green() *TextBuilder {
	t.green = true
	return t
}

func (t *TextBuilder) Red() *TextBuilder {
	t.red = true
	return t
}

func (t *TextBuilder) Yellow() *TextBuilder {
	t.yellow = true
	return t
}

func (t *TextBuilder) Cyan() *TextBuilder {
	t.cyan = true
	return t
}

func (t *TextBuilder) Bold() *TextBuilder {
	t.bold = true
	return t
}

func (t *TextBuilder) Prefix(prefix string) *TextBuilder {
	t.prefix = prefix
	return t
}

func (t *TextBuilder) Print() {
	c := color.New()
	if t.green {
		c.Add(color.FgGreen)
	}
	if t.red {
		c.Add(color.FgRed)
	}
	if t.yellow {
		c.Add(color.FgYellow)
	}
	if t.cyan {
		c.Add(color.FgCyan)
	}
	if t.bold {
		c.Add(color.Bold)
	}
	msg := t.text
	if t.prefix != "" {
		msg = t.prefix + " " + msg
	}
	c.Println(msg)
}
