package actions

import (
	"commitea/internal/pkg/common"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"net"
	"os"
	"runtime"
)

// TODO
// store the message from the directory
// dict {dir: message}

type socketMsg string

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

func socketListener(info socketInfo, ch chan<- tea.Msg) {
	// Remove the socket if it exists
	if _, err := os.Stat(info.address); err == nil {
		err = os.Remove(info.address)
		if err != nil {
			common.HandleError(err)
		}
	}

	// Listen on the Unix socket
	listener, err := net.Listen(info.network, info.address)
	if err != nil {
		common.HandleError(err)
	}
	defer listener.Close()
	fmt.Printf("Listening socket type %s address: %s", info.network, info.address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			common.HandleError(err)
		}
		go func(conn net.Conn) {
			defer conn.Close()
			buf := make([]byte, 1024)
			n, err := conn.Read(buf)
			if err != nil {
				common.HandleError(err)
			}
			ch <- socketMsg(buf[:n])
		}(conn)
	}
}

type socketInfo struct {
	address, network string
}

func Watch() {
	var info socketInfo
	if runtime.GOOS == "windows" {
		info = socketInfo{"127.0.0.1", "tcp"}
	} else {
		info = socketInfo{"/tmp/commitea.sock", "unix"}
	}

	msgChannel := make(chan tea.Msg)
	go socketListener(info, msgChannel)
	p := tea.NewProgram(model{})
	go func() {
		for msg := range msgChannel {
			p.Send(msg)
		}
	}()

	_, err := p.Run()
	if err != nil {
		fmt.Println("Error starting Bubble Tea program:", err)
		os.Exit(1)
	}
}
