package internal

import (
	"log"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/ssh"
	"github.com/yorukot/ssh.yorukot.me/content"
	"github.com/yorukot/ssh.yorukot.me/internal/constants"
	"github.com/yorukot/ssh.yorukot.me/internal/keymap"
	"github.com/yorukot/ssh.yorukot.me/internal/styles"
	"github.com/yorukot/ssh.yorukot.me/pkg/pathutil"
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
		
		blogs: blogPosts,
	}

	return &m, []tea.ProgramOption{}
}

func (m *Model) Init() tea.Cmd {
	// default values
	return tea.Batch(
		tea.RequestBackgroundColor,
	)
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.ColorProfileMsg:
		m.colorProfile = msg.String()
	case tea.BackgroundColorMsg:
		if msg.IsDark() {
			m.bg = "dark"
		}
	case tea.WindowSizeMsg:
		m.windowsSizeChange(msg)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		}
	case tea.MouseWheelMsg:
		switch msg.Button {
		}
	}
	return m, nil
}

func (m *Model) View() tea.View {
	headerContent := m.header.Render()

	content, err := m.renderContent()
	if err != nil {
		content = "Error: " + err.Error()
	}
	
	innerContent := lipgloss.JoinVertical(lipgloss.Left, headerContent, content)
	inner := styles.InnerBox(m.innerWidth, m.innerHeight).Render(innerContent)

	final := lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Top, inner)
	v := tea.NewView(final)
	v.AltScreen = true
	return v
}
