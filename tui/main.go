package main

import (
	"log"
	"os"
	"t1me-tui/api"
	"t1me-tui/ui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/joho/godotenv"
)

const envPath = "../.env"

func main() {
	if err := godotenv.Load(envPath); err != nil {
		log.Fatalf("Error loading .env file from %s", envPath)
	}
	clientBaseURL := os.Getenv("CLIENT_BASE_URL")
	client := api.NewClient(clientBaseURL)
	model := ui.New(client)

	p := tea.NewProgram(
		model,
		tea.WithAltScreen(),
	)

	_, err := p.Run()
	if err != nil {
		panic(err)
	}
}
