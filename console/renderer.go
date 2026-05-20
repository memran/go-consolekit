package console

import (
	"fmt"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
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

type TUIRenderer struct {
	program *tea.Program
	model   *tuiModel
	mu      sync.Mutex
	started bool
	done    chan struct{}
	err     error
}

func NewTUIRenderer() *TUIRenderer {
	model := newTuiModel()
	return &TUIRenderer{
		model: model,
		done:  make(chan struct{}),
	}
}

func (r *TUIRenderer) start() {
	r.mu.Lock()
	if r.started {
		r.mu.Unlock()
		return
	}
	r.started = true
	r.program = tea.NewProgram(r.model,
		tea.WithAltScreen(),
		tea.WithoutSignalHandler(),
	)
	r.mu.Unlock()

	go func() {
		defer close(r.done)
		if _, err := r.program.Run(); err != nil {
			r.err = err
		}
	}()
}

func (r *TUIRenderer) send(msg tuiMsg) {
	r.mu.Lock()
	if !r.started {
		r.mu.Unlock()
		return
	}
	prog := r.program
	r.mu.Unlock()
	prog.Send(msg)
}

func (r *TUIRenderer) Line(text string) {
	r.start()
	r.send(tuiMsg{msgType: "line", text: text})
}

func (r *TUIRenderer) Info(text string) {
	r.start()
	r.send(tuiMsg{msgType: "info", text: text})
}

func (r *TUIRenderer) Success(text string) {
	r.start()
	r.send(tuiMsg{msgType: "success", text: text})
}

func (r *TUIRenderer) Warning(text string) {
	r.start()
	r.send(tuiMsg{msgType: "warning", text: text})
}

func (r *TUIRenderer) Error(text string) {
	r.start()
	r.send(tuiMsg{msgType: "error", text: text})
}

func (r *TUIRenderer) Title(text string) {
	r.start()
	r.send(tuiMsg{msgType: "title", text: text})
}

func (r *TUIRenderer) Stop() {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.started && r.program != nil {
		r.program.Quit()
	}
}

func (r *TUIRenderer) Wait() error {
	r.mu.Lock()
	if !r.started {
		r.mu.Unlock()
		return nil
	}
	r.mu.Unlock()
	<-r.done
	return r.err
}
