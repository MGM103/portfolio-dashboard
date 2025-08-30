package tui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
)

type menuItem struct {
	title, desc string
	targetPage  uint
}

func (i menuItem) Title() string       { return i.title }
func (i menuItem) Description() string { return i.desc }
func (i menuItem) FilterValue() string { return i.title }
func (i menuItem) TargetPage() uint    { return i.targetPage }
func (i menuItem) Action(m *model) {
	m.state = i.targetPage
}

func NewList() list.Model {
	items := []list.Item{
		menuItem{"View Portfolio", "See all current positions", portfolio},
		menuItem{"Add Position", "Add a new coin to your portfolio", addPosition},
		menuItem{"Remove Position", "Remove a previous held position from your portfolio", removePosition},
		menuItem{"View Watch List", "See the prices of assets you have an eye on", watchlist},
		menuItem{"Add to Watch List", "Add a new asset to you watch list", addAsset},
		menuItem{"Remove from Watch List", "Remove an asset to you watch list", removeAsset},
	}
	defaultHeight := 50
	defaultWidth := 25

	l := list.New(items, list.NewDefaultDelegate(), defaultHeight, defaultWidth)

	l.SetShowTitle(false)
	l.SetShowStatusBar(false)
	l.KeyMap.Quit = key.NewBinding(
		key.WithKeys("q"),
		key.WithHelp("q", "quit"),
	)

	return l
}
