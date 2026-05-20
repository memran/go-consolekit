package main

import (
	"fmt"
	"go-consolekit/console"
	"strings"
	"time"
)

type TuiDemoCommand struct{}

func (c *TuiDemoCommand) Name() string {
	return "tui:demo"
}

func (c *TuiDemoCommand) Description() string {
	return "Demonstrate TUI renderer output"
}

func (c *TuiDemoCommand) Configure(config *console.CommandConfig) {
	config.Argument("name").Required().Description("Your name")
}

func (c *TuiDemoCommand) Handle(ctx *console.Context) error {
	name := ctx.Arg("name")

	ctx.Title("TUI Demo")
	ctx.Success("Hello, " + name + "!")
	ctx.Info("This is an info message")
	ctx.Line("Just a plain line of text")

	ctx.Title("Status Summary")
	ctx.Success("3 users loaded")
	ctx.Warning("1 user is inactive")
	ctx.Info("Last sync: " + time.Now().Format("15:04:05"))
	ctx.Error("No errors found")

	ctx.Title("Simulated Work")
	for i := 0; i < 10; i++ {
		ctx.Line(fmt.Sprintf("  Processing batch %d/10...", i+1))
		time.Sleep(50 * time.Millisecond)
	}
	ctx.Success("All batches processed")

	ctx.Title("User Info Table")
	renderTable(ctx, []string{"Name", "Role", "Status"},
		[]string{name, "Admin", "Active"},
		[]string{"Marwa", "User", "Active"},
		[]string{"Alex", "Editor", "Inactive"},
	)

	ctx.Title("Done")
	ctx.Line("Press q to quit the TUI.")

	return nil
}

func renderTable(ctx *console.Context, headers []string, rows ...[]string) {
	var buf strings.Builder
	widths := make([]int, len(headers))
	for i, h := range headers {
		widths[i] = len(h)
	}
	for _, row := range rows {
		for i, cell := range row {
			if len(cell) > widths[i] {
				widths[i] = len(cell)
			}
		}
	}

	sep := func(a, bC, c string) {
		buf.WriteString(a)
		for i, w := range widths {
			buf.WriteString(strings.Repeat("─", w+2))
			if i < len(widths)-1 {
				buf.WriteString(bC)
			}
		}
		buf.WriteString(c)
		buf.WriteByte('\n')
	}

	row := func(cells []string) {
		buf.WriteString("│")
		for i, cell := range cells {
			buf.WriteString(" ")
			buf.WriteString(cell)
			buf.WriteString(strings.Repeat(" ", widths[i]-len(cell)))
			buf.WriteString(" │")
		}
		buf.WriteByte('\n')
	}

	sep("┌", "┬", "┐")
	row(headers)
	sep("├", "┼", "┤")
	for _, r := range rows {
		row(r)
	}
	sep("└", "┴", "┘")

	ctx.Line(strings.TrimRight(buf.String(), "\n"))
}
