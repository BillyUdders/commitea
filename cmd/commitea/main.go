package main

import (
	"commitea/internal/pkg/actions"
	"commitea/internal/pkg/common"
	"errors"
	"fmt"
	"maps"
	"os"
	"slices"
)

func main() {
	verbs := map[string]bool{
		"commit": true,
		"log":    true,
		"status": true,
		"watch":  true,
		"sync":   true,
	}
	if len(os.Args) < 2 {
		common.Exit(inputErr("No command provided", verbs))
	}
	command := os.Args[1]
	if _, ok := verbs[command]; !ok {
		common.Exit(inputErr("Unknown command: "+command, verbs))
	}
	switch command {
	case "commit":
		actions.RunCommitForm()
	case "log":
		actions.RunLog()
	case "status":
		actions.RunStatus(20)
	case "sync":
		fmt.Println("Executing 'sync' command...")
	}
}

func inputErr(command string, m map[string]bool) error {
	return errors.New(fmt.Sprintf("%s. Use one of: %s \n", command, slices.Collect(maps.Keys(m))))
}
