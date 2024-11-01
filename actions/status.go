package actions

import (
	"commitea/common"
	"fmt"
	"github.com/charmbracelet/lipgloss/list"
)

func RunStatus(numOfCommits int) common.GitStatus {
	status, err := common.NewGitObserver("").Status(numOfCommits)
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
