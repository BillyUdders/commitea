package actions

import (
	"commitea/internal/pkg/common"
	"fmt"
	"github.com/charmbracelet/lipgloss/list"
)

func RunStatus(numOfCommits int) common.GitStatus {
	obs, err := common.NewGitObserver("")
	if err != nil {
		common.HandleError(err)
	}
	status, err := obs.Status(numOfCommits)
	if err != nil {
		common.HandleError(err)
	}
	fmt.Println(
		list.New(
			"Files", list.New(status.Files),
			"Branches", list.New(status.Branches),
			"Commits", list.New(status.Commits),
		).ItemStyle(common.InfoText),
	)
	println()
	return status
}
