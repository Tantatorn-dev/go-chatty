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

	roomCode *string
	username string

	err error
}

func InitialModel() Model {
	vp := viewport.New(30, 5)
	vp.SetContent(`Welcome to the chat room!
Please enter room code to enter the chatroom.`)

	ti := textinput.New()
	ti.Placeholder = "Enter room code"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 25

	// random username in format user-1234
	username := fmt.Sprintf("user-%d", time.Now().Unix()%10000)

	return Model{
		viewport:  vp,
		textInput: ti,
		username:  username,
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
			if m.roomCode != nil {
				msg := m.textInput.Value()

				c.SendMessage(ctx, &proto.SendMessageRequest{Code: *m.roomCode, Msg: &proto.Msg{
					Message:  msg,
					Username: m.username,
				}})

				res, _ := c.GetMessages(ctx, &proto.GetMessagesRequest{Code: *m.roomCode})

				var messages string
				for _, m := range res.Messages {
					messages += fmt.Sprintf("%s: %s\n", m.Username, m.Message)
				}

				m.viewport.SetContent(fmt.Sprintf("Room %s\n\n%s", *m.roomCode, messages))
				m.textInput.Reset()
			} else {
				code := m.textInput.Value()

				m.viewport.SetContent(fmt.Sprintf("%s is joining room %s...", m.username, code))

				m.roomCode = &code

				res, _ := c.JoinRoom(ctx, &proto.JoinRoomRequest{Code: code, Username: m.username})

				var messages string
				for _, m := range res.Messages {
					messages += fmt.Sprintf("%s: %s\n", m.Username, m.Message)
				}

				m.viewport.SetContent(fmt.Sprintf("Room %s\n\n%s", code, messages))

				m.textInput.Reset()
				m.textInput.Placeholder = "Type your message here"
			}

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
