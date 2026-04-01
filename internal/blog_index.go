package internal

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
	contentpkg "github.com/yorukot/ssh.yorukot.me/content"
	"github.com/yorukot/ssh.yorukot.me/internal/constants"
	"github.com/yorukot/ssh.yorukot.me/internal/styles"
)

func (m *Model) isBlogIndex() bool {
	return m.path == "/blog"
}

func (m *Model) ensureBlogPosts() error {
	posts, err := contentpkg.BlogPosts()
	if err != nil {
		return err
	}

	m.blogPosts = posts
	m.blogLineOffsets = nil
	if len(m.blogPosts) == 0 {
		m.blogIndex = 0
		return nil
	}

	m.blogIndex = min(max(m.blogIndex, 0), len(m.blogPosts)-1)
	return nil
}

func (m *Model) moveBlogSelection(delta int) {
	if !m.isBlogIndex() {
		m.scrollBy(delta * constants.LineScrollStep)
		return
	}

	if err := m.ensureBlogPosts(); err != nil || len(m.blogPosts) == 0 {
		return
	}

	m.blogIndex = min(max(m.blogIndex+delta, 0), len(m.blogPosts)-1)
	m.syncBlogScrollOffset()
	m.wrappedLines = nil
}

func (m *Model) openSelectedBlog() {
	if !m.isBlogIndex() {
		m.navigateTo("/blog")
		return
	}

	if err := m.ensureBlogPosts(); err != nil || len(m.blogPosts) == 0 {
		return
	}

	m.navigateTo(m.blogPosts[m.blogIndex].Path)
}

func (m *Model) renderBlogIndex(width int) string {
	if err := m.ensureBlogPosts(); err != nil {
		return fmt.Sprintf("failed to load blog index\n\n%s", err)
	}

	m.blogLineOffsets = make([]int, 0, len(m.blogPosts))
	sections := []string{
		styles.BlogIndexTitle(m.bg).Render("Blog"),
		styles.BlogIndexIntro(m.bg).Render("Use up/down to choose a post, then press enter to open it."),
		"",
	}
	lineCount := lipgloss.Height(strings.Join(sections, "\n"))

	for i, post := range m.blogPosts {
		m.blogLineOffsets = append(m.blogLineOffsets, lineCount)

		titleStyle := styles.BlogIndexCard(m.bg)
		if i == m.blogIndex {
			titleStyle = styles.BlogIndexCardActive(m.bg)
		}

		metaParts := make([]string, 0, 1)
		if post.PublishDate != "" {
			metaParts = append(metaParts, post.PublishDate)
		}

		cardParts := []string{titleStyle.Render(post.Title)}
		if len(metaParts) > 0 {
			cardParts = append(cardParts, styles.BlogIndexMeta(m.bg).Render("  "+strings.Join(metaParts, "  |  ")))
		}
		cardParts = append(cardParts, styles.BlogIndexDescription(m.bg).Render("  "+post.Description))

		card := lipgloss.JoinVertical(lipgloss.Left, cardParts...)
		sections = append(sections, card, "")
		lineCount += lipgloss.Height(card) + 1
	}

	if len(m.blogPosts) == 0 {
		sections = append(sections, styles.BlogIndexDescription(m.bg).Render("No posts yet."))
	}

	return lipgloss.NewStyle().Width(width).Render(strings.Join(sections, "\n"))
}

func (m *Model) syncBlogScrollOffset() {
	if !m.isBlogIndex() || len(m.blogPosts) == 0 || len(m.blogLineOffsets) != len(m.blogPosts) {
		return
	}

	selectedTop := m.blogLineOffsets[m.blogIndex]
	selectedBottom := selectedTop + 3
	viewportHeight := m.contentViewportHeight()
	if selectedTop < m.scrollOffset {
		m.scrollOffset = selectedTop
	}
	if selectedBottom >= m.scrollOffset+viewportHeight {
		m.scrollOffset = selectedBottom - viewportHeight + 1
	}
	m.scrollOffset = max(m.scrollOffset, constants.MinScrollOffset)
}
