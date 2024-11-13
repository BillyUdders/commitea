package actions

import (
	"commitea/internal/pkg/common"
	"fmt"
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
	fmt.Println(status.AsList().String())
	return status
}
