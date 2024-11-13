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
			"Files", formatList(status.Files),
			"Branches", formatList(status.Branches),
			"Commits", formatList(status.Commits),
		).ItemStyle(common.InfoText),
	)
	return status
}

func formatList(items []string) *list.List {
	return list.New(items).Enumerator(common.CommiteaEnumerator).EnumeratorStyle(common.WarningText)
}
