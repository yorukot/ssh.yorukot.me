package blogindex

import (
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/yorukot/ssh.yorukot.me/content"
	"github.com/yorukot/ssh.yorukot.me/internal/styles"
)

type Metrics struct {
	LineStarts  []int
	LineHeights []int
}

func Render(posts []content.BlogPost, width, selected int, bg string) (string, Metrics) {
	metrics := Metrics{}

	if len(posts) == 0 {
		title := styles.BlogIndexTitle(bg).Render("Blog")
		intro := styles.BlogIndexIntro(bg).Width(width).Render("No posts yet.")
		return lipgloss.JoinVertical(lipgloss.Left, title, "", intro), metrics
	}

	title := styles.BlogIndexTitle(bg).Render("Blog")
	intro := styles.BlogIndexIntro(bg).Width(width).Render("Use up/down to choose a post, then press enter to open it.")
	line := lipgloss.Height(title) + 1 + lipgloss.Height(intro) + 1

	items := make([]string, 0, len(posts))
	for i, post := range posts {
		item := renderItem(post, width, i == selected, bg)
		items = append(items, item)
		metrics.LineStarts = append(metrics.LineStarts, line)
		metrics.LineHeights = append(metrics.LineHeights, lipgloss.Height(item))
		line += lipgloss.Height(item)
	}

	parts := []string{title, "", intro, ""}
	parts = append(parts, items...)
	return lipgloss.JoinVertical(lipgloss.Left, parts...), metrics
}

func renderItem(post content.BlogPost, width int, selected bool, bg string) string {
	cardStyle := styles.BlogIndexCard(bg)
	titleStyle := styles.BlogIndexTitle(bg)
	metaStyle := styles.BlogIndexMeta(bg)
	pathStyle := styles.BlogIndexPath(bg)

	if selected {
		cardStyle = styles.BlogIndexCardActive(bg)
		titleStyle = titleStyle.Bold(true)
		pathStyle = styles.BlogIndexPathActive(bg)
	}

	meta := post.PublishDate
	if post.Description != "" {
		if meta != "" {
			meta += " | "
		}
		meta += post.Description
	}

	lines := []string{
		titleStyle.Width(width).Render(post.Title),
	}

	if meta != "" {
		lines = append(lines, metaStyle.Width(width).Render(meta))
	}

	lines = append(lines, pathStyle.Width(width).Render(post.Path))

	return cardStyle.Width(width).MarginBottom(1).Render(strings.Join(lines, "\n"))
}
