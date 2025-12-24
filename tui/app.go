package tui

import (
	"fmt"
	"log"
	"sort"
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
	updateAsset
)

type model struct {
	inputs       inputFields
	list         list.Model
	notification string
	state        uint
	store        *data.Store
	table        overview
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
		tempModel, cmd := m.table.Update(msg)
		m.table = tempModel.(overview)
		cmds = append(cmds, cmd)
	case addAsset, addPosition, removeAsset, removePosition, updateAsset:
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
				m.notification = ""
				if item, ok := m.list.SelectedItem().(menuItem); ok {
					switch item.TargetPage() {
					case portfolio:
						colHeaders := []table.Column{
							{Title: "Asset", Width: 10},
							{Title: "ID", Width: 10},
							{Title: "Amount", Width: 10},
							{Title: "Price", Width: 10},
							{Title: "Value", Width: 10},
						}

						positions, _ := m.store.GetPositions()

						var assetIds []string
						for _, p := range positions {
							assetIds = append(assetIds, p.ID)
						}

						assetDetails, _ := api.GetAssetData(assetIds, "AUD")

						idToPice := make(map[string]float64, len(positions))
						for _, d := range assetDetails {
							id := strconv.Itoa(d.Id)
							idToPice[id] = d.Price
						}

						var rows []table.Row
						var portfolioValue float64
						for _, p := range positions {
							price := idToPice[p.ID]
							value := idToPice[p.ID] * p.Amount
							rows = append(rows, table.Row{p.Ticker, p.ID, fmt.Sprintf("%.4f", p.Amount), fmt.Sprintf("%.2f", price), fmt.Sprintf("%.2f", value)})
							portfolioValue += idToPice[p.ID] * p.Amount
						}

						sort.Slice(rows, func(i, j int) bool {
							price1, _ := strconv.ParseFloat(rows[i][4], 64)
							price2, _ := strconv.ParseFloat(rows[j][4], 64)

							return price1 > price2
						})

						m.table = NewOverview(colHeaders, rows, tableContent{footer: fmt.Sprintf("Portfolio value: %f", portfolioValue)})
						m.table.Focus()

					case watchlist:
						colHeaders := []table.Column{
							{Title: "Asset", Width: 10},
							{Title: "ID", Width: 10},
						}

						assets, _ := m.store.GetWatchlist()

						var rows []table.Row
						for _, a := range assets {
							rows = append(rows, table.Row{a.Ticker, a.ID})
						}

						m.table = NewOverview(colHeaders, rows, tableContent{})
						m.table.Focus()

					case addPosition:
						m.inputs = NewInputFields(2, []string{"Asset id...", "Position amount..."})

					case addAsset:
						m.inputs = NewInputFields(1, []string{"Asset id..."})

					case removeAsset:
						assets, _ := m.store.GetWatchlist()

						watchlistDesc := ""
						for _, a := range assets {
							watchlistDesc += fmt.Sprintf("%s(%s)\n", a.Ticker, a.ID)
						}

						m.inputs = NewInputFields(1, []string{"Asset id..."})
						m.inputs.description += watchlistDesc

					case removePosition:
						positions, _ := m.store.GetPositions()

						positionDesc := ""
						for _, p := range positions {
							positionDesc += fmt.Sprintf("%s(%s)  %f\n", p.Ticker, p.ID, p.Amount)
						}

						m.inputs = NewInputFields(1, []string{"Asset id..."})
						m.inputs.description = positionDesc
					}

					item.Action(&m)
				}
			}

		case portfolio:
			switch key {
			case "esc":
				m.state = menu

			case "enter":
				selectedRow := m.table.table.SelectedRow()
				if len(selectedRow) > 0 {
					assetId := selectedRow[1]
					amount := selectedRow[2]
					m.inputs = NewInputFields(2, []string{"Asset id...", "Position amount..."})
					m.inputs.SetValues([]string{assetId, amount})
					m.state = updateAsset
				}

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
				inputValues := m.inputs.GetValues()
				assetIdField := strings.Fields(inputValues[0])
				assetAmountField := strings.Fields(inputValues[1])
				assetData, _ := api.GetAssetData(assetIdField, "AUD")
				assetTicker := assetData[0].Ticker

				if len(assetIdField) != 1 || len(assetAmountField) != 1 {
					m.inputs.description = "Please enter a single asset id and an amount.\n"
					m.inputs.ClearValues()
					break
				}

				positionAmount, _ := strconv.ParseFloat(assetAmountField[0], 64)
				positionDetails := data.Asset{ID: assetIdField[0], Ticker: assetTicker, Amount: positionAmount}

				m.store.SaveToPositions(positionDetails)

				m.notification = fmt.Sprintf("Position added: %s\t%f\n", positionDetails.Ticker, positionAmount)
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
				inputValues := m.inputs.GetValues()
				assetIds := strings.Fields(inputValues[0])
				positions, _ := m.store.GetPositions()

				idToTicker := make(map[string]string, len(positions))
				for _, asset := range positions {
					idToTicker[asset.ID] = asset.Ticker
				}

				var removedPositions []string
				for _, id := range assetIds {
					removedPositions = append(removedPositions, idToTicker[id])
				}

				err := m.store.RemoveFromPositions(assetIds)
				if err != nil {
					m.notification = fmt.Sprintf("Failed to remove assets from positions: %s", err)
				}

				m.notification = "The following were removed from the watchlist:\n" + strings.Join(removedPositions, "\n")
				m.state = menu

			case "ctrl+c":
				return m, tea.Quit
			}

		case updateAsset:
			switch key {
			case "esc":
				m.state = menu

			case "enter":
				inputValues := m.inputs.GetValues()
				assetIdField := strings.Fields(inputValues[0])
				assetAmountField := strings.Fields(inputValues[1])
				assetData, _ := api.GetAssetData(assetIdField, "AUD")
				assetTicker := assetData[0].Ticker

				if len(assetIdField) != 1 || len(assetAmountField) != 1 {
					m.inputs.description = "Please enter a single asset id and an amount.\n"
					m.inputs.ClearValues()
					break
				}

				positionAmount, _ := strconv.ParseFloat(assetAmountField[0], 64)
				positionDetails := data.Asset{ID: assetIdField[0], Ticker: assetTicker, Amount: positionAmount}

				m.store.SaveToPositions(positionDetails)

				m.notification = fmt.Sprintf("Position updated: %s\t%f\n", positionDetails.Ticker, positionAmount)
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

	case addAsset, addPosition, removeAsset, removePosition, updateAsset:
		s += m.inputs.View()
		s += "\n\n"
	}

	return s
}
