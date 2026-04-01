package internal

import (
	"time"

	tea "charm.land/bubbletea/v2"
	contentpkg "github.com/yorukot/ssh.yorukot.me/content"
	"github.com/yorukot/ssh.yorukot.me/internal/constants"
)

type footerQuoteTickMsg struct{}

func footerQuoteTickCmd() tea.Cmd {
	return tea.Tick(constants.FooterQuoteTickInterval, func(_ time.Time) tea.Msg {
		return footerQuoteTickMsg{}
	})
}

func (m *Model) updateFooterQuote() {
	quoteRunes := []rune(contentpkg.GetContent().FooterQuote)
	if len(quoteRunes) == 0 {
		return
	}

	m.footerCursorVisible = !m.footerCursorVisible

	if m.footerQuotePause > 0 {
		m.footerQuotePause--
		m.wrappedLines = nil
		return
	}

	if m.footerQuoteDeleting {
		if m.footerQuoteIndex > 0 {
			m.footerQuoteIndex--
		} else {
			m.footerQuoteDeleting = false
			m.footerQuotePause = constants.FooterQuotePauseTicks
		}
		m.wrappedLines = nil
		return
	}

	if m.footerQuoteIndex < len(quoteRunes) {
		m.footerQuoteIndex++
	} else {
		m.footerQuoteDeleting = true
		m.footerQuotePause = constants.FooterQuotePauseTicks
	}

	m.wrappedLines = nil
}

func (m Model) footerQuoteText() string {
	quoteRunes := []rune(contentpkg.GetContent().FooterQuote)
	if len(quoteRunes) == 0 {
		return ""
	}

	index := m.footerQuoteIndex
	if index < 0 {
		index = 0
	}
	if index > len(quoteRunes) {
		index = len(quoteRunes)
	}

	return string(quoteRunes[:index])
}
