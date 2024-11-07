package actions

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"net"
	"os"
)

const (
	ADDRESS = "/tmp/commitea.sock"
)

type socketMsg string

// TODO
// store the message from the directory
// dict {dir: message}

type model struct {
	messages []string
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case socketMsg:
		// TODO
		// cd into dir
		// call commitea log, status
		m.messages = append(m.messages, string(msg))

	case tea.KeyMsg:

		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	s := "Messages:\n"
	for _, msg := range m.messages {
		s += msg + "\n"
	}
	return s
}

func listenUnixSocket(ch chan<- tea.Msg) {
	socketPath := ADDRESS

	// Remove the socket if it exists
	if _, err := os.Stat(socketPath); err == nil {
		os.Remove(socketPath)
	}

	// Listen on the Unix socket
	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		fmt.Println("Error listening on socket:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Listening on Unix socket:", socketPath)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go func(conn net.Conn) {
			defer conn.Close()
			buf := make([]byte, 1024)
			n, err := conn.Read(buf)
			if err != nil {
				fmt.Println("Error reading from connection:", err)
				return
			}
			ch <- socketMsg(string(buf[:n]))
		}(conn)
	}
}

func Watch() {

	// TODO
	// initialize the GUI

	msgChannel := make(chan tea.Msg)
	go listenUnixSocket(msgChannel)

	p := tea.NewProgram(model{})

	go func() {
		for msg := range msgChannel {
			p.Send(msg)
		}
	}()

	if err := p.Start(); err != nil {
		fmt.Println("Error starting Bubble Tea program:", err)
		os.Exit(1)
	}
}
