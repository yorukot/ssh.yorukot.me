package constants

import "time"

const (
	MouseWheelStep = 3

	MaxContentWidth = 100
	MinContentWidth = 20

	HelpWidthInset    = 2
	ContentWidthInset = 5
	HeaderFrameInset  = 2

	LayoutVerticalSpacing = 5

	MinViewportHeight = 3

	MinScrollbarHeight = 1
	MinScrollbarOffset = 1
	MinScrollOffset    = 0

	LineScrollStep = 5

	MaxScrollOffsetSentinel = 1 << 30

	FooterMinWrapWidth = 12

	HeaderBoxPaddingTop  = 0
	HeaderBoxPaddingSide = 1

	InnerBoxPaddingTop  = 1
	InnerBoxPaddingSide = 1

	FooterQuotePauseTicks = 6
)

const FooterQuoteTickInterval = 120 * time.Millisecond
