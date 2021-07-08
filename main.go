package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

const SV_PATH string = "/etc/runit/sv"

type Service struct {
	Name    string
	Running bool
	Path    string
}

func printServices(services []*Service) {
	max := 0

	svRunning := make([]*Service, 0)
	svStopped := make([]*Service, 0)

	// Calculate max service name length
	// Also separate running services from stopped ones
	for _, i := range services {
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
		len(services),
	)

	svOrdered := make([]*Service, 0, len(svRunning)+len(svStopped))
	svOrdered = append(svOrdered, svRunning...)
	svOrdered = append(svOrdered, svStopped...)

	for _, i := range svOrdered {
		status := "\033[91;1mSTOPPED\033[0m"
		if i.Running {
			status = "\033[92;1mRUNNING\033[0m"
		}

		// Print service
		fmt.Printf("\033[1m%s\033[0m%s%s\n", i.Name, strings.Repeat(" ", max-len(i.Name)+5), status)
	}
}

func listServices() []*Service {
	// List service dirs
	svDirs, err := ioutil.ReadDir(SV_PATH)
	if err != nil {
		panic(err)
	}

	services := make([]*Service, 0)
	for _, i := range svDirs {
		svEnabled := false
		f, err := ioutil.ReadFile(path.Join(SV_PATH, i.Name(), "supervise", "stat"))

		// User has insufficient permissions
		if os.IsPermission(err) {
			fmt.Println("svm is unable to list services; are you sure you are running this as root?")
			os.Exit(1)
		}

		// An unknown error occurred
		if err != nil && !os.IsNotExist(err) {
			panic(err)
		}

		// PID file exists, service is running
		if err == nil && string(f) == "run\n" {
			svEnabled = true
		}

		services = append(services, &Service{
			Name:    i.Name(),
			Running: svEnabled,
			Path:    path.Join(SV_PATH, i.Name()),
		})
	}

	return services
}

func stopService(service *Service) error {
	if !service.Running {
		fmt.Println("service is already stopped")
		os.Exit(1)
	}

	// Open control file for writing
	f, err := os.Create(path.Join(service.Path, "supervise", "control"))
	if err != nil {
		panic(err)
	}

	// Write stop command
	f.Write([]byte("d"))

	return f.Close()
}

func findService(services []*Service, name string) *Service {
	var sv *Service

	for _, i := range services {
		if i.Name == name {
			sv = i
			break
		}
	}

	return sv
}

func main() {
	// List services
	services := listServices()

	if len(os.Args) < 2 {
		fmt.Println("too few arguments")
		os.Exit(1)
	}

	args := os.Args[1:]

	switch args[0] {
	case "list":
	case "l":
		printServices(services)
	case "stop":
	case "d":
		sv := findService(services, args[1])
		if sv == nil {
			fmt.Println("no such service")
			os.Exit(1)
		}

		stopService(sv)
	default:
		fmt.Println("unknown subcommand")
		os.Exit(1)
	}
}
