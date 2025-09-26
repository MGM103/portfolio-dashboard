package main

import (
	"log"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/joho/godotenv"
	data "github.com/mgm103/portfolio-dashboard/data"
	"github.com/mgm103/portfolio-dashboard/tui"
)

func loadEnvFile() {
	if _, err := os.Stat(".env"); err == nil {
		_ = godotenv.Load(".env")
		return
	}

	configDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatal(err)
	}

	configEnvPath := filepath.Join(configDir, "portfolio-dashboard", ".env")
	if _, err := os.Stat(configEnvPath); err == nil {
		_ = godotenv.Load(configEnvPath)
	}
}

func main() {
	loadEnvFile()

	store := &data.Store{}
	m := tui.NewModel(store)
	p := tea.NewProgram(m)

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
