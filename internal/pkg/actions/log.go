package actions

import (
	"commitea/internal/pkg/common"
	"fmt"
	"github.com/charmbracelet/lipgloss/list"
)

func RunLog() {
	obs, err := common.NewGitObserver("")
	if err != nil {
		common.HandleError(err)
	}
	status, err := obs.Status(25)
	if err != nil {
		common.HandleError(err)
	}
	fmt.Println(list.New("Commits", status.Commits).ItemStyle(common.InfoText))
}
