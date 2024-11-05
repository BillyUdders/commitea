package actions

import (
	"commitea/internal/app/commitea/common"
	"fmt"
	"github.com/charmbracelet/lipgloss/list"
)

func RunLog() {
	status, err := common.NewGitObserver("").Status(25)
	if err != nil {
		common.HandleError(err)
	}
	fmt.Println(list.New("Commits", status.Commits).ItemStyle(common.InfoText))
}
