package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type tableContent struct {
	header string
	footer string
	help   string
}

type overview struct {
	table   table.Model
	content tableContent
}

func NewOverview(cols []table.Column, rows []table.Row, content tableContent) overview {
	styles := table.DefaultStyles()
	styles.Selected = lipgloss.NewStyle().
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(true)

	o := overview{
		table: table.New(
			table.WithColumns(cols),
			table.WithRows(rows),
			table.WithFocused(true),
			table.WithStyles(styles),
		),
		content: content,
	}
	o.table.SetHeight(len(rows) + 1)

	return o
}

func (o overview) Init() tea.Cmd {
	return nil
}

func (o overview) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			return o, cmd
		case "j":
			o.table, cmd = o.table.Update(tea.KeyMsg{Type: tea.KeyDown})
			return o, cmd
		case "k":
			o.table, cmd = o.table.Update(tea.KeyMsg{Type: tea.KeyUp})
			return o, cmd
		}
	}

	o.table, cmd = o.table.Update(msg)
	return o, cmd
}

func (o overview) View() string {
	var b strings.Builder

	b.WriteString(o.content.header)
	b.WriteString("\n\n")
	b.WriteString(o.table.View())
	b.WriteString("\n\n")
	b.WriteString(o.content.footer)
	b.WriteRune('\n')
	b.WriteString(o.content.help)

	return b.String()
}

func (o *overview) SetContent(c tableContent) {
	o.content = c
}

func (o *overview) Focus() {
	o.table.Focus()
}
