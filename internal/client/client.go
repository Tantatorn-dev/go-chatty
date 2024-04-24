package client

import (
	"context"
	"fmt"
	"go-chatty/proto"
	"log"
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type errMsg error

type Model struct {
	viewport  viewport.Model
	textInput textinput.Model
	roomCode  *string
	messages  []string
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

	conn, err := grpc.Dial("localhost:8090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := proto.NewChatProtoClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			code := m.textInput.Value()

			m.roomCode = &code
			m.viewport.SetContent(fmt.Sprintf("Joining room %s...", code))

			c.JoinRoom(ctx, &proto.JoinRoomRequest{Code: code})

			m.viewport.SetContent(fmt.Sprintf("Joined room %s", code))

			res, _ := c.GetMessages(ctx, &proto.GetMessagesRequest{Code: code})

			m.messages = res.Messages

			m.textInput.Reset()
			m.textInput.Placeholder = "Type your message here"
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
