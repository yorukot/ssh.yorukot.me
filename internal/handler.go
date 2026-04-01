package internal

import (
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/ssh"
	"github.com/yorukot/ssh.yorukot.me/internal/components/header"
	"github.com/yorukot/ssh.yorukot.me/internal/constants"
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
		width:       pty.Window.Width,
		height:      pty.Window.Height,
		innerWidth:  min(pty.Window.Width, constants.MaxContentWidth),
		innerHeight: pty.Window.Height,
		bg:          "light",
		help:        help.New(),
		keys:        newKeyMap(),
		path:        requestPath,
	}
	m.contentWidth = max(m.innerWidth-constants.ContentWidthInset, constants.MinContentWidth)
	m.help.SetWidth(max(m.innerWidth-constants.HelpWidthInset, constants.MinScrollOffset))
	return &m, []tea.ProgramOption{}
}

func (m *Model) Init() tea.Cmd {
	// default values
	return tea.Batch(
		tea.RequestBackgroundColor,
		footerQuoteTickCmd(),
	)
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.ColorProfileMsg:
		m.profile = msg.String()
	case tea.BackgroundColorMsg:
		if msg.IsDark() {
			m.bg = "dark"
		}
		m.setHelpStyle()
	case footerQuoteTickMsg:
		m.updateFooterQuote()
		return m, footerQuoteTickCmd()
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
		m.innerHeight = m.height
		m.innerWidth = min(m.width, constants.MaxContentWidth)
		m.contentWidth = max(m.innerWidth-constants.ContentWidthInset, constants.MinContentWidth)
		m.help.SetWidth(max(m.innerWidth-constants.HelpWidthInset, constants.MinScrollOffset))
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Up):
			m.moveBlogSelection(-1)
		case key.Matches(msg, m.keys.Down):
			m.moveBlogSelection(1)
		case key.Matches(msg, m.keys.Back):
			m.navigateBack()
		case key.Matches(msg, m.keys.Enter):
			m.openSelectedBlog()
		}
	case tea.MouseWheelMsg:
		switch msg.Button {
		case tea.MouseWheelUp:
			if m.isBlogIndex() {
				m.moveBlogSelection(-1)
			} else {
				m.scrollBy(-constants.MouseWheelStep)
			}
		case tea.MouseWheelDown:
			if m.isBlogIndex() {
				m.moveBlogSelection(1)
			} else {
				m.scrollBy(constants.MouseWheelStep)
			}
		}
	}
	return m, nil
}

func (m *Model) View() tea.View {
	h := header.New(m.innerWidth, m.bg, m.path)
	headerContent := h.Render()
	contentWidth := max(m.innerWidth-constants.ContentWidthInset, constants.MinContentWidth)
	m.contentWidth = contentWidth

	m.cached(contentWidth)
	helpView := m.renderHelpFooter(contentWidth)
	contentView := m.renderScrollableContent(contentWidth, headerContent, helpView)
	innerContent := lipgloss.JoinVertical(lipgloss.Left, headerContent, "", contentView, "", helpView)
	inner := styles.InnerBox(m.innerWidth, m.innerHeight).Render(innerContent)

	final := lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Top, inner)
	v := tea.NewView(final)
	v.AltScreen = true
	return v
}

func (m *Model) scrollBy(delta int) {
	m.scrollOffset = max(m.scrollOffset+delta, constants.MinScrollOffset)
}

func (m *Model) renderHelpFooter(width int) string {
	m.setHelpStyle()
	m.help.SetWidth(width)

	return lipgloss.NewStyle().Width(width).Render(m.help.View(m.keys))
}

func (m *Model) setHelpStyle() {
	m.help.Styles = help.DefaultStyles(m.bg == "dark")
}
