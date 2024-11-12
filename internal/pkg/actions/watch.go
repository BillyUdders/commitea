package actions

import (
	"commitea/internal/pkg/common"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss/list"
	"net"
	"os"
	"runtime"
	"strings"
)

var (
	WindowsSocket = socketInfo{address: "127.0.0.1", network: "tcp"}
	UnixSocket    = socketInfo{address: "/tmp/commitea.sock", network: "unix"}
)

type socketMsg string

type model struct {
	socketInfo socketInfo
	messages   []string
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case socketMsg:
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
	if len(m.messages) == 0 {
		return common.SuccessText.Render("Welcome to watch!")
	}

	path := strings.Trim(m.messages[len(m.messages)-1], " ")
	path = strings.ReplaceAll(path, "\n", "")
	obs, err := common.NewGitObserver(path)
	if err != nil {
		return common.WarningText.Render(path, " is not a Git Repository!")
	}

	status, err := obs.Status(20)
	if err != nil {
		common.HandleError(err)
	}
	s := list.New(
		"Files", list.New(status.Files),
		"Branches", list.New(status.Branches),
		"Commits", list.New(status.Commits),
	).ItemStyle(common.InfoText)

	return s.String()
}

func socketListener(info socketInfo, ch chan<- tea.Msg) {
	// Remove the socket if it exists
	if _, err := os.Stat(info.address); err == nil {
		err = os.Remove(info.address)
		if err != nil {
			common.HandleError(err)
		}
	}

	listener, err := net.Listen(info.network, info.address)
	if err != nil {
		common.HandleError(err)
	}
	defer listener.Close()

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
		info = WindowsSocket
	} else {
		info = UnixSocket
	}

	msgChannel := make(chan tea.Msg)
	go socketListener(info, msgChannel)
	p := tea.NewProgram(model{info, nil})
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
