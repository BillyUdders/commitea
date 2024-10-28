package main

import (
	"commitea/actions"
	"fmt"
	"os"
)

func main() {
	allowedVerbs := map[string]bool{
		"commit": true,
		"sync":   true,
		"ls":     true,
	}

	if len(os.Args) < 2 {
		fmt.Println("Error: No command provided. Use one of: commit, sync, ls")
		os.Exit(1)
	}

	command := os.Args[1]
	if _, ok := allowedVerbs[command]; !ok {
		fmt.Printf("Error: Invalid command '%s'. Use one of: commit, sync, ls\n", command)
		os.Exit(1)
	}

	switch command {
	case "commit":
		actions.RunCommitForm()
	case "sync":
		fmt.Println("Executing 'sync' command...")
	case "ls":
		fmt.Println("Executing 'ls' command...")
	}
}
