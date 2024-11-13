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
	s := strings.ReplaceAll(str, "\n", " ")
	return strings.TrimSpace(regexp.MustCompile(`\s{2,}`).ReplaceAllString(s, " "))
}
