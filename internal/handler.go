package internal

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/key"
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
		keys:   newKeyMap(),
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
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Up):
			m.scrollOffset = max(m.scrollOffset-1, 0)
		case key.Matches(msg, m.keys.Down):
			m.scrollOffset++
		case key.Matches(msg, m.keys.PageUp):
			m.scrollOffset = max(m.scrollOffset-max(m.contentViewportHeight()-1, 1), 0)
		case key.Matches(msg, m.keys.PageDown):
			m.scrollOffset += max(m.contentViewportHeight()-1, 1)
		case key.Matches(msg, m.keys.Home):
			m.scrollOffset = 0
		case key.Matches(msg, m.keys.End):
			m.scrollOffset = 1 << 30
		}
	}
	return m, nil
}

func (m Model) View() tea.View {
	h := header.New(m.innerWidth, m.innerHeight)
	headerContent := h.Render()

	markdownContent, err := contentpkg.MarkdownContent(m.path)
	if err != nil {
		markdownContent = fmt.Sprintf("\n\nfailed to load markdown for %s\n\n%s", m.path, err)
	}

	contentWidth := max(m.innerWidth-5, 20)
	markdownContent = markdownpkg.New(contentWidth, m.bg).Render(markdownContent)
	contentView := m.renderScrollableContent(strings.TrimSpace(markdownContent), contentWidth, headerContent)
	innerContent := lipgloss.JoinVertical(lipgloss.Left, headerContent, "", contentView)
	inner := styles.InnerBox(m.innerWidth, m.innerHeight).Render(innerContent)

	final := lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Top, inner)
	v := tea.NewView(final)
	v.AltScreen = true
	return v
}
