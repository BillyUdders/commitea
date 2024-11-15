package actions

import (
	"commitea/internal/pkg/common"
	tea "github.com/charmbracelet/bubbletea"
	"net"
	"os"
	"runtime"
	"time"
)

var (
	WindowsSocket = socketInfo{address: "127.0.0.1", network: "tcp"}
	UnixSocket    = socketInfo{address: "/tmp/commitea.sock", network: "unix"}
)

type socketMsg string

type refreshMsg string

type socketInfo struct {
	address, network string
}

type model struct {
	socketInfo socketInfo
	msg        socketMsg
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case refreshMsg:
		return m, nil
	case socketMsg:
		m.msg = msg
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	if msg := m.msg; msg == "" {
		return common.SuccessText.Render("Welcome to watch!")
	}
	path := common.TrimAll(string(m.msg))
	obs, err := common.NewGitObserver(path)
	if err != nil {
		return common.WarningText.Render(path, " is not a Git Repository!")
	}
	status, err := obs.Status(20)
	if err != nil {
		common.Exit(err)
	}
	return status.AsList().String()
}

func socketListener(info socketInfo, ch chan<- tea.Msg) {
	if _, err := os.Stat(info.address); err == nil {
		err = os.Remove(info.address)
		if err != nil {
			common.Exit(err)
		}
	}

	listener, err := net.Listen(info.network, info.address)
	if err != nil {
		common.Exit(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			common.Exit(err)
		}
		go func(conn net.Conn) {
			defer conn.Close()
			buf := make([]byte, 1024)
			n, err := conn.Read(buf)
			if err != nil {
				common.Exit(err)
			}
			ch <- socketMsg(buf[:n])
		}(conn)
	}
}

func refresh(interval time.Duration, ch chan tea.Msg) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				ch <- refreshMsg("")
			}
		}
	}()
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
	go refresh(100*time.Millisecond, msgChannel)
	p := tea.NewProgram(model{info, ""})
	go func() {
		for msg := range msgChannel {
			p.Send(msg)
		}
	}()

	_, err := p.Run()
	if err != nil {
		common.Exit(err)
	}
}
