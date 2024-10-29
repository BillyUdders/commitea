package main

import (
	"commitea/actions"
	"commitea/common"
	"fmt"
	"github.com/pkg/errors"
	"os"
)

func main() {
	allowedVerbs := map[string]bool{
		"commit": true,
		"sync":   true,
		"log":    true,
	}
	if len(os.Args) < 2 {
		common.HandleError(errors.New("No command provided. Use one of: commit, sync, ls"))
		os.Exit(1)
	}
	command := os.Args[1]
	if _, ok := allowedVerbs[command]; !ok {
		common.HandleError(errors.New(fmt.Sprintf("Invalid command '%s'. Use one of: commit, sync, log\n", command)))
		os.Exit(1)
	}

	switch command {
	case "commit":
		actions.RunCommitForm()
	case "log":
		actions.RunLog()
	case "sync":
		fmt.Println("Executing 'sync' command...")
	}
}
