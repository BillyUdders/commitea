package main

import (
	"commitea/internal/pkg/actions"
	"commitea/internal/pkg/common"
	"fmt"
	"maps"
	"os"
	"slices"
	"sort"
	"strings"
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
	verbs := slices.Collect(maps.Keys(m))
	sort.Strings(verbs)
	return fmt.Errorf("%s. Use one of: %s", command, strings.Join(verbs, ", "))
}
