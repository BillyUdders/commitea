package actions

import (
	"commitea/internal/pkg/common"
	"fmt"
)

func RunStatus(numOfCommits int) common.GitStatus {
	obs, err := common.NewGitObserver()
	if err != nil {
		common.Exit(err)
	}
	status, err := obs.Status(numOfCommits)
	if err != nil {
		common.Exit(err)
	}
	fmt.Println(status.AsList().String())
	return status
}
