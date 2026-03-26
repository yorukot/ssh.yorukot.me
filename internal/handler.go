package internal

import (
	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/ssh"
	"github.com/yorukot/ssh.yorukot.me/internal/components/header"
	"github.com/yorukot/ssh.yorukot.me/internal/styles"
)

// You can wire any Bubble Tea model up to the middleware with a function that
// handles the incoming ssh.Session. Here we just grab the terminal info and
// pass it to the new model. You can also return tea.ProgramOptions (such as
// tea.WithAltScreen) on a session by session basis.
func TeaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	pty, _, _ := s.Pty()
	m := Model{
		width:     pty.Window.Width,
		height:    pty.Window.Height,
		bg:        "light",
	}
	return m, []tea.ProgramOption{}
}

func (m Model) Init() tea.Cmd {
	// default values
	return tea.Batch(
		tea.RequestBackgroundColor,
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.ColorProfileMsg:
		m.profile = msg.String()
	case tea.BackgroundColorMsg:
		if msg.IsDark() {
			m.bg = "dark"
		}
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m Model) View() tea.View {
	// Header Render
	h := header.New(m.width)
	headerContent := h.Render()
	v := tea.NewView(headerContent + "\n\n" + styles.QuitText().Render("Press 'q' to quit\n"))
	v.AltScreen = true
	return v
}
