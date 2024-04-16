package main

import (
	"log"

	"go-chatty/internal/client"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(client.InitialModel())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
