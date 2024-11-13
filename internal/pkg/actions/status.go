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
			"Files", common.TeaList(status.Files),
			"Branches", common.TeaList(status.Branches),
			"Commits", common.TeaList(status.Commits),
		).ItemStyle(common.InfoText),
	)
	return status
}
