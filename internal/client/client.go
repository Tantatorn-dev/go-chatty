package client

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type errMsg error

type Model struct {
	viewport  viewport.Model
	textInput textinput.Model
	roomCode  string
	err       error
}

func InitialModel() Model {
	vp := viewport.New(30, 5)
	vp.SetContent(`Welcome to the chat room!
Please enter room code to enter the chatroom.`)

	ti := textinput.New()
	ti.Placeholder = "Enter room code"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return Model{
		viewport:  vp,
		textInput: ti,
		err:       nil,
	}
}

func (m Model) Init() tea.Cmd {
	return textarea.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var vpCmd tea.Cmd
	var textInputCmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			m.roomCode = m.textInput.Value()
			m.viewport.SetContent(fmt.Sprintf("Joining room %s...", m.roomCode))
			m.textInput.Reset()
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.viewport, vpCmd = m.viewport.Update(msg)
	m.textInput, textInputCmd = m.textInput.Update(msg)

	return m, tea.Batch(vpCmd, textInputCmd)
}

func (m Model) View() string {
	return fmt.Sprintf(
		"%s\n%s\n",
		m.viewport.View(),
		m.textInput.View(),
	) + "\n\n"
}
