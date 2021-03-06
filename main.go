package main

import (
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	"github.com/keesvv/svm/errs"
	"github.com/keesvv/svm/service"
)

func printServices(serviceList service.ServiceList) {
	const columnMargin int = 8

	max := 0

	svRunning := make(service.ServiceList, 0)
	svStopped := make(service.ServiceList, 0)

	// Calculate max service name length
	// Also separate running services from stopped ones
	for _, i := range serviceList {
		if len(i.Name) > max {
			max = len(i.Name)
		}

		if i.Running {
			svRunning = append(svRunning, i)
		} else {
			svStopped = append(svStopped, i)
		}
	}

	// Print a summary
	fmt.Printf(
		"%d running services, %d stopped, %d total\n\n",
		len(svRunning),
		len(svStopped),
		len(serviceList),
	)

	// Print column headers
	fmt.Printf(
		"NAME%sSTATE%sRUNLEVEL%sMODIFIED\n",
		strings.Repeat(" ", max+columnMargin-4),
		strings.Repeat(" ", columnMargin+2),
		strings.Repeat(" ", columnMargin),
	)

	svOrdered := make(service.ServiceList, 0, len(svRunning)+len(svStopped))
	svOrdered = append(svOrdered, svRunning...)
	svOrdered = append(svOrdered, svStopped...)

	for _, i := range svOrdered {
		status := "\033[91;1mSTOPPED\033[0m"
		if i.Running {
			status = "\033[92;1mRUNNING\033[0m"
		}

		lastModified, err := i.LastModified()
		modifiedStr := fmt.Sprintf(
			"%d mins ago",
			int(math.Floor(time.Since(lastModified).Minutes())),
		)

		if err != nil || lastModified.Before(time.Now().Add(time.Hour*-1)) {
			modifiedStr = ""
		}

		// Print service
		fmt.Printf(
			"\033[1m%s\033[0m%s%s%s%s%s%s\n",
			i.Name,
			strings.Repeat(" ", max-len(i.Name)+columnMargin),
			status,
			strings.Repeat(" ", columnMargin),
			i.Runlevel,
			strings.Repeat(" ", columnMargin*2-len(i.Runlevel)),
			modifiedStr,
		)
	}
}

func main() {
	args := os.Args[1:]

	if len(args) < 1 {
		errs.HandleError(errs.ErrArguments)
	}

	// List services
	services, err := service.ListServices()
	errs.HandleError(err)

	// TODO: refactor
	switch args[0] {
	case "list", "l":
		printServices(services)
	case "stop", "d", "down":
		if len(args) < 2 {
			errs.HandleError(errs.ErrArguments)
		}

		for _, i := range args[1:] {
			sv, err := services.FindByName(i)
			errs.HandleError(err)

			errs.HandleError(sv.Stop())
			fmt.Printf("\033[1m??? Stopped service \033[96m%s\033[0;1m.\033[0m\n", sv.Name)
		}
	case "start", "u", "up":
		if len(args) < 2 {
			errs.HandleError(errs.ErrArguments)
		}

		for _, i := range args[1:] {
			sv, err := services.FindByName(i)
			errs.HandleError(err)

			errs.HandleError(sv.Start())
			fmt.Printf("\033[1m??? Started service \033[96m%s\033[0;1m.\033[0m\n", sv.Name)
		}
	case "enable", "e":
		if len(args) < 2 {
			errs.HandleError(errs.ErrArguments)
		}

		for _, i := range args[1:] {
			sv, err := services.FindByName(i)
			errs.HandleError(err)

			errs.HandleError(sv.SetRunlevel(service.LEVEL_DEFAULT))
			fmt.Printf("\033[1m??? Enabled service \033[96m%s\033[0;1m.\033[0m\n", sv.Name)
		}
	case "disable", "dis":
		if len(args) < 2 {
			errs.HandleError(errs.ErrArguments)
		}

		for _, i := range args[1:] {
			sv, err := services.FindByName(i)
			errs.HandleError(err)

			errs.HandleError(sv.SetRunlevel(service.LEVEL_NONE))
			fmt.Printf("\033[1m??? Disabled service \033[96m%s\033[0;1m.\033[0m\n", sv.Name)
		}
	default:
		errs.HandleError(errs.ErrUnknownSubcommand)
	}
}
