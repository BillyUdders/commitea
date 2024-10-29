package actions

import (
	"commitea/common"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/go-git/go-git/v5/plumbing/object"
)

func RunLog() {
	actor := common.NewGitActor("")

	logCount := 0
	err := actor.Commits.ForEach(func(c *object.Commit) error {
		if logCount >= 20 {
			return nil
		}
		prettyPrintCommit(c)
		logCount++
		return nil
	})
	if err != nil {
		log.Fatalf("Error iterating through commits: %s", err)
	}
}

func prettyPrintCommit(c *object.Commit) {
	idx := strings.Index(c.Message, ":")
	var msg string
	if idx == -1 {
		msg = common.SuccessText.Render(c.Message)
	} else {
		msg = fmt.Sprintf(common.SuccessText.Render(c.Message[:idx]) + c.Message[idx+1:])
	}

	fmt.Println(
		common.InfoText.Foreground(common.Purple).Underline(true).Render(c.Hash.String()[0:6]),
		"-",
		msg,
		common.InfoText.Foreground(common.Gray).Render(fmt.Sprintf("(%s)", formatTime(c.Author.When))),
		common.InfoText.Foreground(common.Cyan).Render(fmt.Sprintf("[%s]", c.Author.Name)),
	)
}

func formatTime(t time.Time) string {
	duration := time.Since(t)
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	seconds := int(duration.Seconds()) % 60

	if hours > 0 {
		return fmt.Sprintf("%d hours", hours)
	} else if minutes > 0 {
		return fmt.Sprintf("%d minutes and %d seconds ago", minutes, seconds)
	} else {
		return fmt.Sprintf("%d seconds ago", seconds)
	}
}
