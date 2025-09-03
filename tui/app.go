package tui

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mgm103/portfolio-dashboard/api"
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
	inputs       inputFields
	list         list.Model
	notification string
	state        uint
	store        *data.Store
	table        table.Model
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

					case addAsset:
						m.inputs = NewInputFields(1, []string{"Asset id..."})

					case removeAsset:
						assets, _ := m.store.GetWatchlist()

						watchlistTickers := ""
						for _, asset := range assets {
							watchlistTickers += asset.Ticker
							watchlistTickers += "\n"
						}

						m.inputs = NewInputFields(1, []string{"Asset id..."})
						m.inputs.description += "Current watchlist:\n" + watchlistTickers

					case removePosition:
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
				inputValues := m.inputs.GetValues()
				assetIds := strings.Fields(inputValues[0])
				assetData, _ := api.GetAssetData(assetIds, "AUD")

				var assetTickers []string
				var assets []data.Asset
				for _, asset := range assetData {
					assetTickers = append(assetTickers, asset.Ticker)
					assets = append(assets, data.Asset{ID: strconv.Itoa(asset.Id), Ticker: asset.Ticker})
				}

				m.store.SaveToWatchlist(assets)

				m.notification = "The following were added to the watchlist:\n" + strings.Join(assetTickers, "\n")
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
				inputValues := m.inputs.GetValues()
				assetIds := strings.Fields(inputValues[0])
				watchlist, _ := m.store.GetWatchlist()

				idToTicker := make(map[string]string, len(watchlist))
				for _, asset := range watchlist {
					idToTicker[asset.ID] = asset.Ticker
				}

				var removedAssets []string
				for _, id := range assetIds {
					removedAssets = append(removedAssets, idToTicker[id])
				}

				m.notification = "The following were removed from the watchlist:\n" + strings.Join(removedAssets, "\n")
				err := m.store.RemoveFromWatchlist(assetIds)
				if err != nil {
					m.notification = fmt.Sprintf("Failed to remove assets from watchlist: %s", err)
				}

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
		s += m.notification
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
