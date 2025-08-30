package tui

import (
	"log"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	data "github.com/mgm103/portfolio-dashboard/data"
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
	store  *data.Store
	list   list.Model
	state  uint
	table  table.Model
	inputs inputFields
}

func NewModel(store *data.Store) model {
	err := store.Init()
	if err != nil {
		log.Fatalf("Failed to set up db: %s", err)
	}

	return model{
		list:  NewList(),
		state: menu,
		store: store,
		table: table.New(),
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

	switch m.state {
	case menu:
		m.list, cmd = m.list.Update(msg)
		cmds = append(cmds, cmd)
	case portfolio, watchlist:
		m.table, cmd = m.table.Update(msg)
		cmds = append(cmds, cmd)
	case addAsset, addPosition, removeAsset, removePosition:
		tempModel, cmd := m.inputs.Update(msg)
		m.inputs = tempModel.(inputFields)
		cmds = append(cmds, cmd)

	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()
		switch m.state {
		case menu:
			switch key {
			case "enter":
				if item, ok := m.list.SelectedItem().(menuItem); ok {
					switch item.TargetPage() {
					case addPosition:
						m.inputs = NewInputFields(2, []string{"Asset id...", "Position amount..."})

					case addAsset, removeAsset, removePosition:
						m.inputs = NewInputFields(1, []string{"Asset id..."})

					}

					item.Action(&m)
				}
			}

		case portfolio:
			switch key {
			case "esc":
				m.state = menu

			case "ctrl+c":
				return m, tea.Quit
			}

		case watchlist:
			switch key {
			case "esc":
				m.state = menu

			case "ctrl+c":
				return m, tea.Quit

			}

		case addAsset:
			switch key {
			case "esc":
				m.state = menu

			case "enter":
				m.state = menu

			case "ctrl+c":
				return m, tea.Quit
			}

		case addPosition:

			switch key {
			case "esc":
				m.state = menu

			case "enter":
				m.state = menu

			case "ctrl+c":
				return m, tea.Quit
			}

		case removeAsset:
			switch key {
			case "esc":
				m.state = menu

			case "enter":
				m.state = menu

			case "ctrl+c":
				return m, tea.Quit
			}

		case removePosition:
			switch key {
			case "esc":
				m.state = menu

			case "enter":
				m.state = menu

			case "ctrl+c":
				return m, tea.Quit
			}

		}
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	var s string

	switch m.state {
	case menu:
		s += m.list.View()
		s += "\n\n"

	case portfolio, watchlist:
		s += m.table.View()
		s += "\n\n"

	case addAsset, addPosition, removeAsset, removePosition:
		s += m.inputs.View()
		s += "\n\n"
	}

	return s
}
