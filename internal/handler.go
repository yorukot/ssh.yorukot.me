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
	markdownpkg "github.com/yorukot/ssh.yorukot.me/pkg/markdown"
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
		m.height = msg.Height
		m.width = msg.Width
		m.innerHeight = m.height
		m.innerWidth = min(m.width, 82)
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
	h := header.New(m.innerWidth, m.innerHeight)
	headerContent := h.Render()

	content += headerContent

	markdownContent, err := contentpkg.MarkdownContent(m.path)
	if err != nil {
		markdownContent = fmt.Sprintf("\n\nfailed to load markdown for %s\n\n%s", m.path, err)
	}

	markdownContent = markdownpkg.New(max(m.innerWidth-2, 40), m.bg).Render(markdownContent)

	content += "\n\n" + strings.TrimSpace(markdownContent)

	inner := styles.InnerBox(m.innerWidth, m.innerHeight).Render(content)

	final := lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Top, inner)
	v := tea.NewView(final)
	v.AltScreen = true
	return v
}
