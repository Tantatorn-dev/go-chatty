package main

import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(initialModel())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

type model struct {
	viewport viewport.Model
	roomCode string
}

func initialModel() model {
	vp := viewport.New(30, 5)
	vp.SetContent(`Welcome to the chat room!
Please enter room code to enter the chatroom.`)

	return model{
		viewport: vp,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var vpCmd tea.Cmd

	m.viewport, vpCmd = m.viewport.Update(msg)

	return m, tea.Batch(vpCmd)
}

func (m model) View() string {
	return fmt.Sprintf(
		"%s\n\n",
		m.viewport.View(),
	) + "\n\n"
}
