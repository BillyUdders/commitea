package actions

import (
	"commitea/common"
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

func RunStatus() {
	actor := common.NewGitActor("")

	stats, err := actor.RepoStats()
	if err != nil {
		common.HandleError(err)
	}
	infoTable := table.New().
		Border(lipgloss.RoundedBorder()).
		BorderColumn(true).
		BorderRow(true).
		BorderStyle(common.SuccessText).
		Rows(stats...)
	fmt.Println(infoTable.Render())
}
