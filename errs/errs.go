package errs

import (
	"errors"
	"fmt"
	"os"
)

var ErrNoSuchService = errors.New("no such service")
var ErrUnknownSubcommand = errors.New("unknown subcommand")

func HandleError(err error) {
	fmt.Printf("\033[91m%s\033[0m\n", err.Error())
	os.Exit(1)
}
