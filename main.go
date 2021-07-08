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
	Enabled bool
}

func printServices(services []*Service) {
	max := 0

	for _, i := range services {
		if len(i.Name) > max {
			max = len(i.Name)
		}
	}

	for _, i := range services {
		status := "\033[91;1mSTOPPED\033[0m"
		if i.Enabled {
			status = "\033[92;1mRUNNING\033[0m"
		}

		fmt.Printf("\033[1m%s\033[0m%s%s\n", i.Name, strings.Repeat(" ", max-len(i.Name)+5), status)
	}
}

func main() {
	svFiles, err := ioutil.ReadDir(SV_PATH)
	if err != nil {
		panic(err)
	}

	services := make([]*Service, 0)
	for _, i := range svFiles {
		svEnabled := false

		f, err := ioutil.ReadFile(path.Join(SV_PATH, i.Name(), "supervise", "pid"))

		if err != nil && !os.IsNotExist(err) {
			panic(err)
		}

		if err == nil && len(f) > 0 {
			svEnabled = true
		}

		services = append(services, &Service{
			Name:    i.Name(),
			Enabled: svEnabled,
		})
	}

	printServices(services)
}
