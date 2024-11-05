package common

import (
	"fmt"
	"os"
)

func HandleError(err error) {
	fmt.Println(ErrorText.Render("Error: " + err.Error()))
	os.Exit(1)
}
