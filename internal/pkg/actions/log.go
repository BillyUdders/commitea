package actions

import (
	"commitea/internal/pkg/common"
	"fmt"
	"github.com/charmbracelet/lipgloss/list"
)

func RunLog() {
	obs, err := common.NewGitObserver("")
	if err != nil {
		common.Exit(err)
	}
	status, err := obs.Status(25)
	if err != nil {
		common.Exit(err)
	}
	fmt.Println(list.New("Commits", common.SubList(status.Commits)).ItemStyle(common.InfoText))
}
