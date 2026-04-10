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
	blocks := make([]string, 0, parent.ChildCount())
	for child := parent.FirstChild(); child != nil; child = child.NextSibling() {
		block := r.renderBlock(child, source, width)
		if strings.TrimSpace(block) == "" {
			continue
		}
		blocks = append(blocks, block)
	}

	separator := strings.Repeat("\n", max(r.Document.Spacing, 1)+1)
	return strings.Join(blocks, separator)
}

func (r Renderer) renderBlock(node ast.Node, source []byte, width int) string {
	switch n := node.(type) {
	case *ast.Heading:
		text := r.renderInlineChildren(n, source)
		prefix := strings.Repeat("#", n.Level) + " "
		return r.headingStyle(n.Level).Render(wrapText(prefix+text, width))
	case *ast.Paragraph:
		return r.Content.Paragraph.Render(wrapText(r.renderInlineChildren(n, source), width))
	case *ast.TextBlock:
		if n.HasChildren() {
			return r.Content.Text.Render(wrapText(r.renderInlineChildren(n, source), width))
		}
		return r.Content.Text.Render(wrapText(string(n.Lines().Value(source)), width))
	case *ast.Blockquote:
		return r.renderBlockQuote(n, source, width)
	case *ast.List:
		return r.renderList(n, source, width)
	case *extast.Table:
		return r.renderTable(n, source, width)
	case *ast.FencedCodeBlock:
		return r.renderFencedCodeBlock(n, source, width)
	case *ast.CodeBlock:
		return r.Content.CodeBlock.Container.Render(r.Content.CodeBlock.Code.Render(string(n.Lines().Value(source))))
	case *ast.ThematicBreak:
		return r.Content.Rule.Render(strings.Repeat("-", width))
	default:
		if node.HasChildren() {
			return r.renderBlocks(node, source, width)
		}
		return strings.TrimSpace(nodeTextValue(node, source))
	}
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
