package console

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestTuiModelInit(t *testing.T) {
	m := newTuiModel()
	cmd := m.Init()
	if cmd != nil {
		t.Fatal("expected nil cmd from Init")
	}
}

func TestTuiModelWindowSizeMsgMakesReady(t *testing.T) {
	m := newTuiModel()
	if m.ready {
		t.Fatal("expected not ready initially")
	}
	m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	if !m.ready {
		t.Fatal("expected ready after WindowSizeMsg")
	}
	if m.width != 80 || m.height != 24 {
		t.Fatalf("expected 80x24, got %dx%d", m.width, m.height)
	}
}

func TestTuiModelAddsMessages(t *testing.T) {
	m := newTuiModel()
	m.Update(tuiMsg{msgType: "success", text: "hello"})
	m.Update(tuiMsg{msgType: "info", text: "world"})
	if len(m.messages) != 2 {
		t.Fatalf("expected 2 messages, got %d", len(m.messages))
	}
	if m.messages[0].text != "hello" || m.messages[0].msgType != "success" {
		t.Fatal("first message mismatch")
	}
	if m.messages[1].text != "world" || m.messages[1].msgType != "info" {
		t.Fatal("second message mismatch")
	}
}

func TestTuiModelQuitOnQ(t *testing.T) {
	m := newTuiModel()
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	if cmd == nil {
		t.Fatal("expected non-nil cmd for q key")
	}
}

func TestTuiModelQuitOnCtrlC(t *testing.T) {
	m := newTuiModel()
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	if cmd == nil {
		t.Fatal("expected non-nil cmd for Ctrl+C")
	}
}

func TestTuiModelScrollUpDown(t *testing.T) {
	m := newTuiModel()
	m.Update(tea.WindowSizeMsg{Width: 80, Height: 10})

	for i := 0; i < 20; i++ {
		m.Update(tuiMsg{msgType: "line", text: "msg"})
	}

	if m.scroll != 0 {
		t.Fatalf("expected scroll 0 initially, got %d", m.scroll)
	}

	m.Update(tea.KeyMsg{Type: tea.KeyUp})
	if m.scroll != 1 {
		t.Fatalf("expected scroll 1 after up, got %d", m.scroll)
	}

	m.Update(tea.KeyMsg{Type: tea.KeyDown})
	if m.scroll != 0 {
		t.Fatalf("expected scroll 0 after down, got %d", m.scroll)
	}

	m.Update(tea.KeyMsg{Type: tea.KeyDown})
	if m.scroll != 0 {
		t.Fatalf("expected scroll 0 when at bottom, got %d", m.scroll)
	}
}

func TestTuiModelScrollHomeEnd(t *testing.T) {
	m := newTuiModel()
	m.Update(tea.WindowSizeMsg{Width: 80, Height: 10})

	for i := 0; i < 20; i++ {
		m.Update(tuiMsg{msgType: "line", text: "msg"})
	}

	m.Update(tea.KeyMsg{Type: tea.KeyHome})
	maxScroll := len(m.messages) - (m.height - 1)
	if m.scroll != maxScroll {
		t.Fatalf("expected scroll %d after home, got %d", maxScroll, m.scroll)
	}

	m.Update(tea.KeyMsg{Type: tea.KeyEnd})
	if m.scroll != 0 {
		t.Fatalf("expected scroll 0 after end, got %d", m.scroll)
	}
}

func TestTuiModelViewNotReady(t *testing.T) {
	m := newTuiModel()
	v := m.View()
	if !strings.Contains(v, "initializing") {
		t.Fatalf("expected 'initializing' in view when not ready, got: %s", v)
	}
}

func TestTuiModelViewEmpty(t *testing.T) {
	m := newTuiModel()
	m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	v := m.View()
	if !strings.Contains(v, "No output yet") {
		t.Fatalf("expected 'No output yet' in view when empty, got: %s", v)
	}
}

func TestTuiModelViewShowsMessages(t *testing.T) {
	m := newTuiModel()
	m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	m.Update(tuiMsg{msgType: "line", text: "hello world"})
	v := m.View()
	if !strings.Contains(v, "hello world") {
		t.Fatalf("expected message in view, got: %s", v)
	}
	if !strings.Contains(v, "q quit") {
		t.Fatalf("expected quit hint in view, got: %s", v)
	}
}

func TestRenderTuiLineTypes(t *testing.T) {
	tests := []struct {
		msgType string
		text    string
	}{
		{"line", "plain"},
		{"success", "done"},
		{"info", "note"},
		{"warning", "caution"},
		{"error", "fail"},
		{"title", "Section"},
	}
	for _, tt := range tests {
		l := styledLine{msgType: tt.msgType, text: tt.text}
		result := renderTuiLine(l)
		if !strings.Contains(result, tt.text) {
			t.Errorf("renderTuiLine(%s, %q) missing text: %s", tt.msgType, tt.text, result)
		}
	}
}

func TestNewTUIRendererImplementsRenderer(t *testing.T) {
	var r Renderer = NewTUIRenderer()
	if r == nil {
		t.Fatal("expected non-nil TUIRenderer")
	}
}

func TestTUIRendererStopBeforeStart(t *testing.T) {
	r := NewTUIRenderer()
	defer func() {
		if err := recover(); err != nil {
			t.Fatalf("Stop() panicked: %v", err)
		}
	}()
	r.Stop()
}

func TestTUIRendererWaitBeforeStart(t *testing.T) {
	r := NewTUIRenderer()
	err := r.Wait()
	if err != nil {
		t.Fatalf("Wait() before start returned error: %v", err)
	}
}

func TestTUIRendererMethods(t *testing.T) {
	r := NewTUIRenderer()
	defer r.Stop()
	defer func() {
		if err := recover(); err != nil {
			t.Fatalf("method panicked: %v", err)
		}
	}()
	r.Line("test")
	r.Info("test")
	r.Success("test")
	r.Warning("test")
	r.Error("test")
	r.Title("test")
}
