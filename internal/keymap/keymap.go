package keymap

import "charm.land/bubbles/v2/key"

type Bindings struct {
	Up    key.Binding
	Down  key.Binding
	Back  key.Binding
	Enter key.Binding
	Quit  key.Binding
}

func New() Bindings {	
	return Bindings{
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("up/k", "scroll up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("down/j", "scroll down"),
		),
		Back: key.NewBinding(
			key.WithKeys("backspace", "left"),
			key.WithHelp("left/backspace", "go home"),
		),
		Enter: key.NewBinding(
			key.WithKeys("enter", "right"),
			key.WithHelp("enter/right", "enter blog"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
	}
}

func (k Bindings) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Back, k.Enter, k.Quit}
}

func (k Bindings) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Back, k.Enter},
		{k.Quit},
	}
}
