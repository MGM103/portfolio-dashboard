package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/joho/godotenv"
	data "github.com/mgm103/portfolio-dashboard/data"
	"github.com/mgm103/portfolio-dashboard/tui"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	store := &data.Store{}
	m := tui.NewModel(store)
	p := tea.NewProgram(m)

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
