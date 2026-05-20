package console

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type tuiMsg struct {
	msgType string
	text    string
}

type tuiModel struct {
	messages []styledLine
	scroll   int
	width    int
	height   int
	ready    bool
}

type styledLine struct {
	msgType string
	text    string
}

var (
	tuiStyleTitle   = lipgloss.NewStyle().Bold(true)
	tuiStyleSuccess = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	tuiStyleInfo    = lipgloss.NewStyle().Foreground(lipgloss.Color("6"))
	tuiStyleWarning = lipgloss.NewStyle().Foreground(lipgloss.Color("3"))
	tuiStyleError   = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
	tuiStyleFooter  = lipgloss.NewStyle().Faint(true)
)

func newTuiModel() *tuiModel {
	return &tuiModel{}
}

func (m *tuiModel) Init() tea.Cmd {
	return nil
}

func (m *tuiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			m.scroll++
		case "down", "j":
			if m.scroll > 0 {
				m.scroll--
			}
		case "home", "g":
			maxScroll := len(m.messages) - (m.height - 1)
			if maxScroll < 0 {
				maxScroll = 0
			}
			m.scroll = maxScroll
		case "end", "G":
			m.scroll = 0
		}
	case tuiMsg:
		m.messages = append(m.messages, styledLine{msgType: msg.msgType, text: msg.text})
	}
	return m, nil
}

func (m *tuiModel) View() string {
	if !m.ready {
		return "ConsoleKit TUI initializing..."
	}
	if len(m.messages) == 0 {
		return "\n  No output yet.\n\n  Press q to quit.\n"
	}

	footerLines := 1
	contentHeight := m.height - footerLines
	if contentHeight < 1 {
		contentHeight = 1
	}

	maxScroll := len(m.messages) - contentHeight
	if maxScroll < 0 {
		maxScroll = 0
	}
	if m.scroll > maxScroll {
		m.scroll = maxScroll
	}
	if m.scroll < 0 {
		m.scroll = 0
	}

	start := len(m.messages) - contentHeight - m.scroll
	if start < 0 {
		start = 0
	}
	end := start + contentHeight
	if end > len(m.messages) {
		end = len(m.messages)
	}

	var s strings.Builder
	for _, line := range m.messages[start:end] {
		s.WriteString(renderTuiLine(line))
		s.WriteByte('\n')
	}

	for i := end - start; i < contentHeight; i++ {
		s.WriteByte('\n')
	}

	pct := 0
	if maxScroll > 0 {
		pct = int(float64(m.scroll) / float64(maxScroll) * 100)
	}
	footer := tuiStyleFooter.Render(fmt.Sprintf("↑/↓ scroll • %d%% • q quit", pct))
	s.WriteString(footer)

	return s.String()
}

func renderTuiLine(l styledLine) string {
	switch l.msgType {
	case "title":
		return tuiStyleTitle.Render("── " + l.text + " ──")
	case "success":
		return tuiStyleSuccess.Render("✓ " + l.text)
	case "info":
		return tuiStyleInfo.Render("ℹ " + l.text)
	case "warning":
		return tuiStyleWarning.Render("⚠ " + l.text)
	case "error":
		return tuiStyleError.Render("✗ " + l.text)
	default:
		return l.text
	}
}
