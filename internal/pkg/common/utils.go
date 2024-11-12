package common

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

func HandleError(err error) {
	fmt.Println(ErrorText.Render("Error: " + err.Error()))
	os.Exit(1)
}

func TrimAll(str string) string {
	processed := strings.ReplaceAll(str, "\n", " ")
	processed = regexp.MustCompile(`\s{2,}`).ReplaceAllString(processed, " ")
	return strings.TrimSpace(processed)
}
