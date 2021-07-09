package errs

import (
	"errors"
	"fmt"
	"os"
)

var ErrArguments = errors.New("too few arguments")
var ErrNoSuchService = errors.New("no such service")
var ErrUnknownSubcommand = errors.New("unknown subcommand")
var ErrIsStopped = errors.New("service is already stopped")
var ErrIsStarted = errors.New("service is already running")
var ErrPermission = errors.New("svm is unable to list services; are you sure you are running this as root?")
var ErrNoRunlevel = errors.New("service has no runlevel, starting/stopping it is not yet supported")

func HandleError(err error) {
	if err == nil {
		return
	}

	fmt.Printf("\033[1merror:\033[0m \033[91m%s\033[0m\n", err.Error())
	os.Exit(1)
}
