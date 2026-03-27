package markdown

import (
	"bytes"
	"regexp"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/alecthomas/chroma/v2/quick"
)

const minWidth = 24

var (
	frontMatterDelimiter = regexp.MustCompile(`(?m)^(-{3}|\.{3})\s*$`)
	fencePattern         = regexp.MustCompile("^\\s*(```|~~~)([^`]*)$")
	headingPattern       = regexp.MustCompile(`^(#{1,6})(\s+)(.*)$`)
	blockquotePattern    = regexp.MustCompile(`^(\s*>+\s?)(.*)$`)
	unorderedListPattern = regexp.MustCompile(`^(\s*[-+*]\s+)(.*)$`)
	orderedListPattern   = regexp.MustCompile(`^(\s*\d+\.\s+)(.*)$`)
	rulePattern          = regexp.MustCompile(`^\s{0,3}((\*\s*){3,}|(-\s*){3,}|(_\s*){3,})\s*$`)
	inlinePattern        = regexp.MustCompile("!?\\[[^\\]]*\\]\\([^\\)]*\\)|`[^`]+`|\\*\\*[^*]+\\*\\*|__[^_]+__|\\*[^*]+\\*|_[^_]+_")
)

type Markdown struct {
	Width         int
	Background    string
	Text          lipgloss.Style
	Muted         lipgloss.Style
	HeadingMarker lipgloss.Style
	HeadingText   lipgloss.Style
	QuoteMarker   lipgloss.Style
	QuoteText     lipgloss.Style
	ListMarker    lipgloss.Style
	Rule          lipgloss.Style
	Fence         lipgloss.Style
	CodeText      lipgloss.Style
	InlineCode    lipgloss.Style
	Link          lipgloss.Style
	Image         lipgloss.Style
	Emphasis      lipgloss.Style
	Strong        lipgloss.Style
}

func New(width int, background string) Markdown {
	baseWidth := max(width, minWidth)
	text := "252"
	muted := "244"
	headingMarker := "213"
	headingText := "229"
	quote := "151"
	list := "221"
	link := "117"
	image := "180"
	inlineCode := "222"
	codeBG := "236"
	strong := "230"
	emphasis := "186"
	if background == "light" {
		text = "236"
		muted = "244"
		headingMarker = "161"
		headingText = "25"
		quote = "31"
		list = "166"
		link = "26"
		image = "95"
		inlineCode = "124"
		codeBG = "255"
		strong = "18"
		emphasis = "55"
	}

	return Markdown{
		Width:         baseWidth,
		Background:    background,
		Text:          lipgloss.NewStyle().Foreground(lipgloss.Color(text)),
		Muted:         lipgloss.NewStyle().Foreground(lipgloss.Color(muted)),
		HeadingMarker: lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(headingMarker)),
		HeadingText:   lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(headingText)),
		QuoteMarker:   lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(quote)),
		QuoteText:     lipgloss.NewStyle().Foreground(lipgloss.Color(quote)),
		ListMarker:    lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(list)),
		Rule:          lipgloss.NewStyle().Foreground(lipgloss.Color(muted)),
		Fence:         lipgloss.NewStyle().Foreground(lipgloss.Color(muted)),
		CodeText:      lipgloss.NewStyle().Foreground(lipgloss.Color(text)).Background(lipgloss.Color(codeBG)),
		InlineCode:    lipgloss.NewStyle().Foreground(lipgloss.Color(inlineCode)).Background(lipgloss.Color(codeBG)).Bold(true),
		Link:          lipgloss.NewStyle().Underline(true).Foreground(lipgloss.Color(link)).Bold(true),
		Image:         lipgloss.NewStyle().Italic(true).Foreground(lipgloss.Color(image)),
		Emphasis:      lipgloss.NewStyle().Italic(true).Foreground(lipgloss.Color(emphasis)),
		Strong:        lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(strong)),
	}
}

func (m Markdown) Render(content string) string {
	content = strings.TrimSpace(stripFrontMatter(content))
	if content == "" {
		return ""
	}

	lines := strings.Split(content, "\n")
	parts := make([]string, 0, len(lines))

	inFence := false
	fence := ""
	language := ""
	codeLines := make([]string, 0)

	flushCode := func() {
		if len(codeLines) == 0 {
			return
		}
		parts = append(parts, m.renderHighlightedCode(strings.Join(codeLines, "\n"), language))
		codeLines = codeLines[:0]
	}

	for _, line := range lines {
		if !inFence {
			if matches := fencePattern.FindStringSubmatch(line); matches != nil {
				inFence = true
				fence = matches[1]
				language = strings.TrimSpace(matches[2])
				parts = append(parts, m.Fence.Render(line))
				continue
			}

			parts = append(parts, m.renderLine(line))
			continue
		}

		if strings.HasPrefix(strings.TrimSpace(line), fence) {
			flushCode()
			parts = append(parts, m.Fence.Render(line))
			inFence = false
			fence = ""
			language = ""
			continue
		}

		codeLines = append(codeLines, line)
	}

	if inFence {
		flushCode()
	}

	return strings.Join(parts, "\n")
}

func (m Markdown) renderLine(line string) string {
	if strings.TrimSpace(line) == "" {
		return ""
	}

	if rulePattern.MatchString(line) {
		return m.Rule.Render(line)
	}

	if matches := headingPattern.FindStringSubmatch(line); matches != nil {
		return m.HeadingMarker.Render(matches[1]) + matches[2] + m.renderInline(matches[3], m.HeadingText)
	}

	if matches := blockquotePattern.FindStringSubmatch(line); matches != nil {
		return m.QuoteMarker.Render(matches[1]) + m.renderInline(matches[2], m.QuoteText)
	}

	if matches := unorderedListPattern.FindStringSubmatch(line); matches != nil {
		return m.ListMarker.Render(matches[1]) + m.renderInline(matches[2], m.Text)
	}

	if matches := orderedListPattern.FindStringSubmatch(line); matches != nil {
		return m.ListMarker.Render(matches[1]) + m.renderInline(matches[2], m.Text)
	}

	return m.renderInline(line, m.Text)
}

func (m Markdown) renderInline(line string, base lipgloss.Style) string {
	locs := inlinePattern.FindAllStringIndex(line, -1)
	if len(locs) == 0 {
		return base.Render(line)
	}

	var b strings.Builder
	last := 0
	for _, loc := range locs {
		if loc[0] > last {
			b.WriteString(base.Render(line[last:loc[0]]))
		}

		token := line[loc[0]:loc[1]]
		b.WriteString(m.renderToken(token))
		last = loc[1]
	}

	if last < len(line) {
		b.WriteString(base.Render(line[last:]))
	}

	return b.String()
}

func (m Markdown) renderToken(token string) string {
	switch {
	case strings.HasPrefix(token, "!["):
		return m.Image.Render(token)
	case strings.HasPrefix(token, "["):
		return m.Link.Render(token)
	case strings.HasPrefix(token, "`"):
		return m.InlineCode.Render(token)
	case strings.HasPrefix(token, "**") || strings.HasPrefix(token, "__"):
		return m.Strong.Render(token)
	case strings.HasPrefix(token, "*") || strings.HasPrefix(token, "_"):
		return m.Emphasis.Render(token)
	default:
		return m.Text.Render(token)
	}
}

func (m Markdown) renderHighlightedCode(code, language string) string {
	code = strings.TrimRight(code, "\n")
	if code == "" {
		return ""
	}

	theme := "monokai"
	if m.Background == "light" {
		theme = "github"
	}

	var buf bytes.Buffer
	if err := quick.Highlight(&buf, code, fallbackLanguage(language), "terminal256", theme); err == nil {
		return strings.TrimRight(buf.String(), "\n")
	}

	return m.CodeText.Render(code)
}

func stripFrontMatter(content string) string {
	content = strings.TrimLeft(content, "\ufeff")
	if !strings.HasPrefix(content, "---\n") && !strings.HasPrefix(content, "---\r\n") {
		return content
	}

	loc := frontMatterDelimiter.FindAllStringIndex(content, 2)
	if len(loc) < 2 {
		return content
	}

	return strings.TrimLeft(content[loc[1][1]:], "\r\n")
}

func fallbackLanguage(language string) string {
	if language == "" {
		return "text"
	}
	return language
}
