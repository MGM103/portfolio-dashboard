package tui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	menu uint = iota
	portfolio
	watchlist
	addPosition
	removePosition
	addAsset
	removeAsset
)

type model struct {
	// store *Store
	list      list.Model
	state     uint
	table     table.Model
	textInput textinput.Model
}

func NewModel() model {
	return model{
		list:      NewList(),
		state:     menu,
		table:     table.New(),
		textInput: textinput.New(),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmds []tea.Cmd
		cmd  tea.Cmd
	)

	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)

	m.table, cmd = m.table.Update(msg)
	cmds = append(cmds, cmd)

	m.textInput, cmd = m.textInput.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()
		switch m.state {
		case menu:
			switch key {
			case "enter":
				if item, ok := m.list.SelectedItem().(menuItem); ok {
					item.Action(&m)
				}
			}
		case portfolio:
			switch key {
			case "esc":
				m.state = menu
			}
		case watchlist:
			switch key {
			case "esc":
				m.state = menu
			}
		case addPosition:
			switch key {
			case "esc":
				m.state = menu
			}
		}
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	var s string

	if m.state == menu {
		s += m.list.View()
		s += "\n\n"
	}

	if m.state == portfolio {
		s += m.table.View()
		s += "\n\n"
	}

	if m.state == addPosition {
		s += m.textInput.View()
		s += "\n\n"
	}

	return s
}
