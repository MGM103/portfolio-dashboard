package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
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
	o := overview{table: table.New(table.WithColumns(cols), table.WithRows(rows)), content: content}
	o.table.SetHeight(len(rows) + 1)

	return o
}

func (o overview) Init() tea.Cmd {
	return nil
}

func (o overview) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
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
