package actions

import (
	"commitea/internal/pkg/common"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fsnotify/fsnotify"
	"net"
	"os"
	"runtime"
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
	watcher    *fsnotify.Watcher
	gitObs     *common.GitObserver
	msg        socketMsg
	msgHistory []socketMsg
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case refreshMsg:
		return m, nil
	case socketMsg:
		lastPath := m.msgHistory[len(m.msgHistory)-1]
		path := socketMsg(common.TrimAll(string(msg)))
		if lastPath != path {
			err := m.watcher.Add(string(path))
			if err != nil {
				return nil, tea.Quit
			}
			err = m.watcher.Remove(string(path))
			if err != nil {
				return nil, tea.Quit
			}
			observer, err := common.NewGitObserver(string(path))
			if err != nil {
				return nil, tea.Quit
			}
			m.gitObs = observer
			m.msgHistory = append(m.msgHistory, path)
			m.msg = path
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	if m.msg == "" {
		return common.SuccessText.Render("Welcome to watch!")
	}
	if m.gitObs != nil {
		return common.WarningText.Render(string(m.msg), " is not a Git Repository!")
	}
	status, err := m.gitObs.Status(20)
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

func refresh(watcher *fsnotify.Watcher, ch chan tea.Msg) {
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					ch <- refreshMsg("")
				}
			case _, ok := <-watcher.Errors:
				if !ok {
					return
				}
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

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		common.Exit(err)
	}
	defer watcher.Close()

	msgChannel := make(chan tea.Msg)
	go socketListener(info, msgChannel)
	go refresh(watcher, msgChannel)

	p := tea.NewProgram(
		model{
			socketInfo: info,
			watcher:    watcher,
			gitObs:     nil,
		},
	)

	go func() {
		for msg := range msgChannel {
			p.Send(msg)
		}
	}()

	_, err = p.Run()
	if err != nil {
		common.Exit(err)
	}
}
