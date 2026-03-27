package internal

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/ssh"
	contentpkg "github.com/yorukot/ssh.yorukot.me/content"
	"github.com/yorukot/ssh.yorukot.me/internal/components/header"
	"github.com/yorukot/ssh.yorukot.me/internal/styles"
	"github.com/yorukot/ssh.yorukot.me/pkg/pathutil"
)

// You can wire any Bubble Tea model up to the middleware with a function that
// handles the incoming ssh.Session. Here we just grab the terminal info and
// pass it to the new model. You can also return tea.ProgramOptions (such as
// tea.WithAltScreen) on a session by session basis.
func TeaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	pty, _, _ := s.Pty()
	requestPath := "/"
	if command := s.Command(); len(command) > 0 {
		requestPath = pathutil.NormalizePath(command[0])
	}

	m := Model{
		width:  pty.Window.Width,
		height: pty.Window.Height,
		bg:     "light",
		path:   requestPath,
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
		m.screenHeight = msg.Height
		m.screenWidth = msg.Width
		m.width = min(m.screenWidth - 2, 100)
		m.height = m.screenHeight - 2
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
	var content string
	h := header.New(m.width)
	headerContent := h.Render()

	content += headerContent

	markdownContent, err := contentpkg.MarkdownContent(m.path)
	if err != nil {
		markdownContent = fmt.Sprintf("\n\nfailed to load markdown for %s\n\n%s", m.path, err)
	}

	content += "\n\n" + strings.TrimSpace(markdownContent)

	inner := styles.FullScreenBox(m.screenHeight, m.screenHeight).Render(content)
	view := lipgloss.Place(
		m.screenWidth,
		m.screenHeight,
		lipgloss.Center,
		lipgloss.Top,
		inner,
	)

	v := tea.NewView(view)
	v.AltScreen = true
	return v
}
