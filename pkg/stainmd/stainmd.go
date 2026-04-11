package stainmd

import (
	"regexp"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	extast "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/text"
)

var ansiPattern = regexp.MustCompile(`\x1b\[[0-9;]*[A-Za-z]`)

type Renderer struct {
	Document DocumentStyle
	Header   HeaderStyle
	Content  ContentStyle
}

func New() Renderer {
	return DefaultStyles()
}

func (r Renderer) Render(markdown string, width int) (string, error) {
	if width < 1 {
		width = 1
	}

	source := []byte(markdown)
	doc := goldmark.New(
		goldmark.WithExtensions(
			extension.Table,
			extension.Strikethrough,
		),
	).Parser().Parse(text.NewReader(source))
	body := r.renderBlocks(doc, source, width)
	return r.Document.Container.Render(body), nil
}

func (r Renderer) renderBlocks(parent ast.Node, source []byte, width int) string {
	return r.renderBlocksWithStyle(parent, source, width, lipgloss.Style{})
}

func (r Renderer) renderBlocksWithStyle(parent ast.Node, source []byte, width int, base lipgloss.Style) string {
	blocks := make([]string, 0, parent.ChildCount())
	for child := parent.FirstChild(); child != nil; child = child.NextSibling() {
		block := r.renderBlockWithStyle(child, source, width, base)
		if strings.TrimSpace(block) == "" {
			continue
		}
		blocks = append(blocks, block)
	}

	separator := strings.Repeat("\n", max(r.Document.Spacing, 1)+1)
	return strings.Join(blocks, separator)
}

func (r Renderer) renderBlock(node ast.Node, source []byte, width int) string {
	return r.renderBlockWithStyle(node, source, width, lipgloss.Style{})
}

func (r Renderer) renderBlockWithStyle(node ast.Node, source []byte, width int, base lipgloss.Style) string {
	switch n := node.(type) {
	case *ast.Heading:
		headingStyle := mergeStyle(base, r.headingStyle(n.Level))
		text := r.renderInlineChildrenWithStyle(n, source, headingStyle)
		prefix := strings.Repeat("#", n.Level) + " "
		return wrapText(headingStyle.Render(prefix)+text, width)
	case *ast.Paragraph:
		return wrapText(r.renderInlineChildrenWithStyle(n, source, mergeStyle(base, r.Content.Paragraph)), width)
	case *ast.TextBlock:
		if n.HasChildren() {
			return wrapText(r.renderInlineChildrenWithStyle(n, source, mergeStyle(base, r.Content.Text)), width)
		}
		return mergeStyle(base, r.Content.Text).Render(wrapText(string(n.Lines().Value(source)), width))
	case *ast.Blockquote:
		return r.renderBlockQuote(n, source, width, mergeStyle(base, r.Content.BlockQuote.Body))
	case *ast.List:
		return r.renderList(n, source, width, mergeStyle(base, r.Content.List.Item))
	case *extast.Table:
		return r.renderTable(n, source, width)
	case *ast.FencedCodeBlock:
		return r.renderFencedCodeBlock(n, source, width)
	case *ast.CodeBlock:
		return r.renderCodeBlock(string(n.Lines().Value(source)), width, "")
	case *ast.ThematicBreak:
		return r.Content.Rule.Render(strings.Repeat("-", width))
	default:
		if node.HasChildren() {
			return r.renderBlocksWithStyle(node, source, width, base)
		}
		return strings.TrimSpace(nodeTextValue(node, source))
	}
}

func mergeStyle(base, own lipgloss.Style) lipgloss.Style {
	return base.Inherit(own)
}

func nodeTextValue(node ast.Node, source []byte) string {
	switch n := node.(type) {
	case *ast.Text:
		return string(n.Value(source))
	case *ast.String:
		return string(n.Value)
	}

	type lineNode interface {
		Lines() *text.Segments
	}

	if n, ok := node.(lineNode); ok {
		return string(n.Lines().Value(source))
	}

	if !node.HasChildren() {
		return ""
	}

	var out strings.Builder
	for child := node.FirstChild(); child != nil; child = child.NextSibling() {
		out.WriteString(nodeTextValue(child, source))
	}
	return out.String()
}

func wrapText(value string, width int) string {
	if width < 1 {
		return value
	}
	return lipgloss.Wrap(value, width, "")
}

func stripANSI(value string) string {
	return ansiPattern.ReplaceAllString(value, "")
}
