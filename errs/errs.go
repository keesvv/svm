package errs

import (
	"fmt"
	"os"
)

func HandleError(err error) {
	fmt.Printf("\033[91m%s\033[0m\n", err.Error())
	os.Exit(1)
}
