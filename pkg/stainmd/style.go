package stainmd

import "charm.land/lipgloss/v2"

func DefaultStyles() Renderer {
	return MochaStyles()
}

func DefaultDocumentStyle() DocumentStyle {
	return DefaultStyles().Document
}

func DefaultHeaderStyle() HeaderStyle {
	return DefaultStyles().Header
}

func DefaultContentStyle() ContentStyle {
	return DefaultStyles().Content
}

func LatteStyles() Renderer {
	return themedRenderer(themePalette{
		text:        "5",
		subtext:     "8",
		blue:        "12",
		lavender:    "13",
		teal:        "14",
		green:       "10",
		yellow:      "11",
		peach:       "209",
		red:         "9",
		mauve:       "13",
		surface:     "252",
		overlay:     "250",
		base:        "230",
		chromaTheme: "github",
	})
}

func MochaStyles() Renderer {
	return themedRenderer(themePalette{
		text:        "252",
		subtext:     "246",
		blue:        "117",
		lavender:    "183",
		teal:        "116",
		green:       "120",
		yellow:      "222",
		peach:       "216",
		red:         "210",
		mauve:       "183",
		surface:     "236",
		overlay:     "240",
		base:        "#1e1e2e",
		chromaTheme: "catppuccin-mocha",
	})
}

type themePalette struct {
	text        string
	subtext     string
	blue        string
	lavender    string
	teal        string
	green       string
	yellow      string
	peach       string
	red         string
	mauve       string
	surface     string
	overlay     string
	base        string
	chromaTheme string
}

func themedRenderer(p themePalette) Renderer {
	return Renderer{
		Document: DocumentStyle{
			Container: lipgloss.NewStyle().Foreground(lipgloss.Color(p.text)),
			Spacing:   1,
		},
		Header: HeaderStyle{
			HeadingOne:   lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(p.blue)),
			HeadingTwo:   lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(p.lavender)),
			HeadingThree: lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(p.teal)),
			HeadingFour:  lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(p.green)),
			HeadingFive:  lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(p.yellow)),
			HeadingSix:   lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(p.peach)),
		},
		Content: ContentStyle{
			Paragraph:     lipgloss.NewStyle().Foreground(lipgloss.Color(p.text)),
			Text:          lipgloss.NewStyle().Foreground(lipgloss.Color(p.text)),
			Strong:        lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(p.text)),
			Emphasis:      lipgloss.NewStyle().Italic(true).Foreground(lipgloss.Color(p.lavender)),
			InlineCode:    lipgloss.NewStyle().Foreground(lipgloss.Color(p.yellow)).Background(lipgloss.Color(p.surface)).Padding(0, 1),
			Strikethrough: lipgloss.NewStyle().Strikethrough(true).Foreground(lipgloss.Color(p.overlay)),
			Rule:          lipgloss.NewStyle().Foreground(lipgloss.Color(p.overlay)),
			Link: LinkStyle{
				Text: lipgloss.NewStyle().Underline(true).Foreground(lipgloss.Color(p.teal)),
				URL:  lipgloss.NewStyle().Foreground(lipgloss.Color(p.teal)),
			},
			Image: ImageStyle{
				Alt:  lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(p.mauve)),
				Path: lipgloss.NewStyle().Foreground(lipgloss.Color(p.subtext)),
			},
			BlockQuote: BlockQuoteStyle{
				Container: lipgloss.NewStyle(),
				Prefix:    lipgloss.NewStyle().Foreground(lipgloss.Color(p.overlay)),
				Body:      lipgloss.NewStyle().Foreground(lipgloss.Color(p.subtext)).Italic(true),
			},
			List: ListStyle{
				Container:   lipgloss.NewStyle(),
				Item:        lipgloss.NewStyle().Foreground(lipgloss.Color(p.text)),
				Bullet:      lipgloss.NewStyle().Foreground(lipgloss.Color(p.blue)).Bold(true),
				Enumeration: lipgloss.NewStyle().Foreground(lipgloss.Color(p.blue)).Bold(true),
				Indent:      2,
				LevelIndent: 2,
			},
			TaskList: TaskListStyle{
				Container: lipgloss.NewStyle(),
				Item:      lipgloss.NewStyle().Foreground(lipgloss.Color(p.text)),
				Ticked:    lipgloss.NewStyle().Foreground(lipgloss.Color(p.green)).Bold(true),
				Unticked:  lipgloss.NewStyle().Foreground(lipgloss.Color(p.overlay)),
			},
			CodeBlock: CodeBlockStyle{
				Container: lipgloss.NewStyle().Background(lipgloss.Color(p.base)).Padding(0, 1),
				Code:      lipgloss.NewStyle().Foreground(lipgloss.Color(p.text)),
				Language:  lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(p.mauve)),
				Theme:     p.chromaTheme,
			},
			Table: TableStyle{
				Container: lipgloss.NewStyle(),
				Header:    lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(p.blue)),
				Cell:      lipgloss.NewStyle().Foreground(lipgloss.Color(p.text)),
				Border:    lipgloss.NewStyle().Foreground(lipgloss.Color(p.overlay)),
			},
		},
	}
}
