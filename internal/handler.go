package internal

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/ssh"
	"github.com/yorukot/ssh.yorukot.me/content"
	"github.com/yorukot/ssh.yorukot.me/internal/components/header"
	"github.com/yorukot/ssh.yorukot.me/internal/components/mkrender"
	"github.com/yorukot/ssh.yorukot.me/internal/constants"
	"github.com/yorukot/ssh.yorukot.me/internal/keymap"
	"github.com/yorukot/ssh.yorukot.me/internal/styles"
	"github.com/yorukot/ssh.yorukot.me/pkg/pathutil"
	"log"
)

// You can wire any Bubble Tea model up to the middleware with a function that
// handles the incoming ssh.Session. Here we just grab the terminal info and
// pass it to the new model. You can also return tea.ProgramOptions (such as
// tea.WithAltScreen) on a session by session basis.
func TeaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	pty, _, _ := s.Pty()
	// Get the request path and parse it
	requestPath := "/"
	if command := s.Command(); len(command) > 0 {
		requestPath = pathutil.NormalizePath(command[0])
	}

	// load blog posts
	blogPosts, err := content.BlogPosts()
	if err != nil {
		log.Fatalf("Error to load the blog posts: %v", err)
	}

	m := Model{
		width:       pty.Window.Width,
		height:      pty.Window.Height,
		innerWidth:  min(pty.Window.Width, constants.MaxContentWidth),
		innerHeight: pty.Window.Height,

		keys: keymap.New(),

		path: requestPath,

		blogs:    blogPosts,
		markdown: mkrender.New(),
	}

	m.header = header.New(m.innerWidth, m.bg, m.path)
	m.syncViewport()

	return &m, []tea.ProgramOption{}
}

func (m *Model) Init() tea.Cmd {
	return tea.Batch(
		tea.RequestBackgroundColor,
	)
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.ColorProfileMsg:
		m.colorProfile = msg.String()
	case tea.BackgroundColorMsg:
		if msg.IsDark() {
			m.bg = "dark"
		}
		m.header = header.New(m.innerWidth, m.bg, m.path)
		m.syncViewport()
	case tea.WindowSizeMsg:
		m.windowsSizeChange(msg)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		}
	}

	m.main, cmd = m.main.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m *Model) View() tea.View {
	headerContent := m.header.Render()
	mainView := m.main.View()
	scrollbarView := m.scrollbarView()

	var body string
	gap := lipgloss.NewStyle().Width(constants.ScrollbarGap).Render("")
	body = lipgloss.JoinHorizontal(lipgloss.Top, mainView, gap, scrollbarView)

	innerContent := lipgloss.JoinVertical(lipgloss.Left, headerContent, body)
	inner := styles.InnerBox(m.innerWidth, m.innerHeight).Render(innerContent)

	final := lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Top, inner)

	v := tea.NewView(final)
	v.AltScreen = true
	return v
}
