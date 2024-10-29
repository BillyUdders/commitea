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
	_, _, commitIter := common.GetGitObjects()

	logCount := 0
	err := commitIter.ForEach(func(c *object.Commit) error {
		if logCount >= 10 {
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
	a := common.InfoText.Render("\uEAFC Commit: ") +
		c.Hash.String()[0:6] +
		common.InfoText.Render(" Author: ") +
		c.Author.Name +
		common.InfoText.Render(" Date: ") +
		formatTime(c.Author.When)

	fmt.Println(a)

	idx := strings.Index(c.Message, ":")
	if idx == -1 {
		fmt.Println(common.SuccessText.PaddingLeft(3).Render(c.Message))
	} else {
		fmt.Println(common.SuccessText.PaddingLeft(3).Render(c.Message[:idx]) + c.Message[idx+1:])
	}
	fmt.Println()
}

func formatTime(t time.Time) string {
	return t.Format("Mon Jan 2 15:04:05 2006 -0700")
}
