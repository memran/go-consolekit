package console

import (
	"fmt"

	"github.com/fatih/color"
)

type Renderer interface {
	Line(text string)
	Info(text string)
	Success(text string)
	Warning(text string)
	Error(text string)
	Title(text string)
}

type CLIRenderer struct{}

func NewCLIRenderer() *CLIRenderer {
	return &CLIRenderer{}
}

func (r *CLIRenderer) Line(text string) {
	fmt.Println(text)
}

func (r *CLIRenderer) Info(text string) {
	color.Cyan("ℹ %s", text)
}

func (r *CLIRenderer) Success(text string) {
	color.Green("✓ %s", text)
}

func (r *CLIRenderer) Warning(text string) {
	color.Yellow("⚠ %s", text)
}

func (r *CLIRenderer) Error(text string) {
	color.Red("✗ %s", text)
}

func (r *CLIRenderer) Title(text string) {
	color.New(color.FgWhite, color.Bold).Printf("\n── %s ──────────────────────────\n\n", text)
}

type TUIRenderer struct{}

func NewTUIRenderer() *TUIRenderer {
	return &TUIRenderer{}
}

func (r *TUIRenderer) Line(text string) {
	fmt.Println(text)
}

func (r *TUIRenderer) Info(text string) {
	fmt.Println("Info:", text)
}

func (r *TUIRenderer) Success(text string) {
	fmt.Println("Success:", text)
}

func (r *TUIRenderer) Warning(text string) {
	fmt.Println("Warning:", text)
}

func (r *TUIRenderer) Error(text string) {
	fmt.Println("Error:", text)
}

func (r *TUIRenderer) Title(text string) {
	fmt.Println("Title:", text)
}
