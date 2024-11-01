package actions

import (
	"commitea/common"
	"fmt"
	"github.com/charmbracelet/lipgloss/list"
)

func RunStatus(numOfCommits int) {
	status, err := common.NewGitObserver("").Status(numOfCommits)
	if err != nil {
		common.HandleError(err)
	}
	fmt.Print(
		list.New(
			"Files", list.New(status.Files),
			"Branches", list.New(status.Branches),
			"Commits", list.New(status.Commits),
		).ItemStyle(common.SuccessText),
	)
}
