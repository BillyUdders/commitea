package actions

import (
	"commitea/internal/pkg/common"
	"fmt"
	"github.com/charmbracelet/lipgloss/list"
)

func RunStatus(dirpath string, numOfCommits int) common.GitStatus {
	obs, err := common.NewGitObserver(dirpath)
	if err != nil {
		common.HandleError(err)
	}
	status, err := obs.Status(numOfCommits)
	if err != nil {
		common.HandleError(err)
	}
	fmt.Println(
		list.New(
			"Files", common.SubList(status.Files),
			"Branches", common.SubList(status.Branches),
			"Commits", common.SubList(status.Commits),
		).ItemStyle(common.InfoText),
	)
	return status
}
