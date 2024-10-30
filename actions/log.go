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
	maxLogCount := 15
	err := actor.Commits.ForEach(func(c *object.Commit) error {
		if logCount >= maxLogCount {
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
		msg = fmt.Sprintf(common.SuccessText.Underline(true).Render(c.Message[:idx]) + c.Message[idx:])
	}
	fmt.Println(
		common.LogText1.Render(c.Hash.String()[0:6]),
		"-",
		msg,
		common.LogText2.Render(fmt.Sprintf("(%s)", formatTime(c.Author.When))),
		common.LogText3.Render(fmt.Sprintf("[%s]", c.Author.Name)),
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
