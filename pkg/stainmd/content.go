package stainmd

import (
	"bytes"
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/table"
	"github.com/alecthomas/chroma/v2/quick"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/yuin/goldmark/ast"
	extast "github.com/yuin/goldmark/extension/ast"
)

type DocumentStyle struct {
	Container lipgloss.Style
	Spacing   int
}

type ContentStyle struct {
	Paragraph     lipgloss.Style
	Text          lipgloss.Style
	Strong        lipgloss.Style
	Emphasis      lipgloss.Style
	InlineCode    lipgloss.Style
	Strikethrough lipgloss.Style
	Rule          lipgloss.Style
	Link          LinkStyle
	Image         ImageStyle
	BlockQuote    BlockQuoteStyle
	List          ListStyle
	TaskList      TaskListStyle
	CodeBlock     CodeBlockStyle
	Table         TableStyle
}

type LinkStyle struct {
	Text lipgloss.Style
	URL  lipgloss.Style
}

type ImageStyle struct {
	Alt  lipgloss.Style
	Path lipgloss.Style
}

type BlockQuoteStyle struct {
	Container lipgloss.Style
	Prefix    lipgloss.Style
	Body      lipgloss.Style
}

type ListStyle struct {
	Container   lipgloss.Style
	Item        lipgloss.Style
	Bullet      lipgloss.Style
	Enumeration lipgloss.Style
	Indent      int
	LevelIndent int
}

type TaskListStyle struct {
	Container lipgloss.Style
	Item      lipgloss.Style
	Ticked    lipgloss.Style
	Unticked  lipgloss.Style
}

type CodeBlockStyle struct {
	Container lipgloss.Style
	Code      lipgloss.Style
	Language  lipgloss.Style
	Theme     string
}

type TableStyle struct {
	Container   lipgloss.Style
	Header      lipgloss.Style
	Cell        lipgloss.Style
	Border      lipgloss.Border
	BorderStyle lipgloss.Style
}

func (r Renderer) renderTable(node *extast.Table, source []byte, width int) string {
	rows := make([][]string, 0, node.ChildCount())

	for rowNode := node.FirstChild(); rowNode != nil; rowNode = rowNode.NextSibling() {
		var cells []string

		switch row := rowNode.(type) {
		case *extast.TableHeader:
			cells = r.renderTableCells(row, source)
		case *extast.TableRow:
			cells = r.renderTableCells(row, source)
		default:
			continue
		}

		rows = append(rows, cells)
	}

	if len(rows) == 0 || len(rows[0]) == 0 {
		return ""
	}

	headers := append([]string(nil), rows[0]...)
	body := make([][]string, 0, len(rows)-1)
	for _, row := range rows[1:] {
		body = append(body, append([]string(nil), row...))
	}

	tbl := table.New().
		Headers(headers...).
		Rows(body...).
		Width(max(width, 1)).
		Wrap(true).
		Border(r.Content.Table.Border).
		BorderStyle(r.Content.Table.BorderStyle).
		BorderTop(true).
		BorderBottom(true).
		BorderLeft(true).
		BorderRight(true).
		BorderColumn(true).
		BorderHeader(true).
		BorderRow(false).
		StyleFunc(func(row, _ int) lipgloss.Style {
			if row == table.HeaderRow {
				return r.Content.Table.Header
			}
			return r.Content.Table.Cell
		})

	return r.Content.Table.Container.Render(tbl.Render())
}

func (r Renderer) renderTableCells(row ast.Node, source []byte) []string {
	cells := make([]string, 0, row.ChildCount())
	for cellNode := row.FirstChild(); cellNode != nil; cellNode = cellNode.NextSibling() {
		cell, ok := cellNode.(*extast.TableCell)
		if !ok {
			continue
		}
		cells = append(cells, strings.TrimSpace(r.renderInlineChildren(cell, source)))
	}
	return cells
}

func (r Renderer) renderBlockQuote(node *ast.Blockquote, source []byte, width int, bodyStyle lipgloss.Style) string {
	prefix := r.Content.BlockQuote.Prefix.Render("│ ")
	innerWidth := max(1, width-lipgloss.Width(prefix))
	body := r.renderBlocksWithStyle(node, source, innerWidth, bodyStyle)
	lines := strings.Split(body, "\n")
	for i, line := range lines {
		if strings.TrimSpace(line) == "" {
			lines[i] = prefix
			continue
		}
		lines[i] = prefix + line
	}

	return r.Content.BlockQuote.Container.Render(strings.Join(lines, "\n"))
}

func (r Renderer) renderList(node *ast.List, source []byte, width int, itemStyle lipgloss.Style) string {
	items := make([]string, 0, node.ChildCount())
	index := node.Start
	if index == 0 {
		index = 1
	}

	for item := node.FirstChild(); item != nil; item = item.NextSibling() {
		marker := "• "
		markerStyle := r.Content.List.Bullet
		if node.IsOrdered() {
			marker = fmt.Sprintf("%d. ", index)
			markerStyle = r.Content.List.Enumeration
			index++
		}

		itemBody := r.renderListItem(item, source, max(1, width-lipgloss.Width(marker)), itemStyle)
		itemLines := strings.Split(itemBody, "\n")
		for i, line := range itemLines {
			prefix := strings.Repeat(" ", lipgloss.Width(marker))
			if i == 0 {
				prefix = markerStyle.Render(marker)
			}
			itemLines[i] = prefix + line
		}

		items = append(items, strings.Join(itemLines, "\n"))
	}

	return r.Content.List.Container.Render(strings.Join(items, "\n"))
}

func (r Renderer) renderListItem(node ast.Node, source []byte, width int, base lipgloss.Style) string {
	parts := make([]string, 0, node.ChildCount())
	for child := node.FirstChild(); child != nil; child = child.NextSibling() {
		parts = append(parts, r.renderBlockWithStyle(child, source, width, base))
	}
	return strings.Join(parts, "\n")
}

func (r Renderer) renderFencedCodeBlock(node *ast.FencedCodeBlock, source []byte, width int) string {
	language := strings.TrimSpace(string(node.Language(source)))
	content := strings.TrimRight(nodeTextValue(node, source), "\n")
	return r.renderCodeBlock(content, width, language)
}

func (r Renderer) renderCodeBlock(content string, width int, language string) string {
	innerWidth := max(1, width-r.Content.CodeBlock.Container.GetHorizontalFrameSize())
	rendered := r.renderHighlightedCodeBlock(content, language)
	wrapped := lipgloss.Wrap(rendered, innerWidth, "")
	return r.Content.CodeBlock.Container.Width(width).Render(wrapped)
}

func (r Renderer) renderHighlightedCodeBlock(content, language string) string {
	if strings.TrimSpace(content) == "" {
		return r.Content.CodeBlock.Code.Render(content)
	}

	var buf bytes.Buffer
	if err := quick.Highlight(&buf, content, language, "terminal256", r.codeBlockStyleName()); err == nil {
		return strings.TrimRight(buf.String(), "\n")
	}

	return r.Content.CodeBlock.Code.Render(content)
}

func (r Renderer) codeBlockStyleName() string {
	style := strings.ToLower(r.Content.CodeBlock.Theme)
	if style == "" {
		return "monokai"
	}
	if styles.Get(style) == nil {
		return "monokai"
	}
	return style
}

func (r Renderer) renderInlineChildren(parent ast.Node, source []byte) string {
	return r.renderInlineChildrenWithStyle(parent, source, lipgloss.Style{})
}

func mergeInlineStyle(base, own lipgloss.Style) lipgloss.Style {
	return own.Inherit(base)
}

func (r Renderer) renderInlineChildrenWithStyle(parent ast.Node, source []byte, base lipgloss.Style) string {
	var out strings.Builder
	for child := parent.FirstChild(); child != nil; child = child.NextSibling() {
		out.WriteString(r.renderInlineNode(child, source, base))
	}
	return out.String()
}

func (r Renderer) renderInlineNode(node ast.Node, source []byte, base lipgloss.Style) string {
	switch n := node.(type) {
	case *ast.Text:
		text := string(n.Value(source))
		if n.HardLineBreak() {
			return base.Render(text) + "\n"
		}
		if n.SoftLineBreak() {
			return base.Render(text) + " "
		}
		return base.Render(text)
	case *ast.String:
		return base.Render(string(n.Value))
	case *ast.CodeSpan:
		style := mergeInlineStyle(base, r.Content.InlineCode)
		return style.Render(r.renderInlineChildrenWithStyle(n, source, lipgloss.Style{}))
	case *ast.Emphasis:
		style := mergeInlineStyle(base, r.Content.Emphasis)
		if n.Level >= 2 {
			style = mergeInlineStyle(base, r.Content.Strong)
		}
		return r.renderInlineChildrenWithStyle(n, source, style)
	case *ast.Link:
		style := mergeInlineStyle(base, r.Content.Link.Text).Hyperlink(string(n.Destination))
		return r.renderInlineChildrenWithStyle(n, source, style)
	case *ast.AutoLink:
		label := string(n.Label(source))
		return mergeInlineStyle(base, r.Content.Link.Text).Hyperlink(string(n.URL(source))).Render(label)
	case *ast.Image:
		label := r.renderInlineChildrenWithStyle(n, source, mergeInlineStyle(base, r.Content.Image.Alt))
		if strings.TrimSpace(label) == "" {
			label = string(n.Destination)
		}
		var parts []string
		parts = append(parts, r.Content.Image.Alt.Render(label))
		if dest := strings.TrimSpace(string(n.Destination)); dest != "" {
			if r.ImagePathResolver != nil {
				dest = r.ImagePathResolver(dest)
			}
			parts = append(parts, r.Content.Image.Path.Render(dest))
		}
		return strings.Join(parts, " ")
	default:
		if node.HasChildren() {
			return r.renderInlineChildrenWithStyle(node, source, base)
		}
		return nodeTextValue(node, source)
	}
}
