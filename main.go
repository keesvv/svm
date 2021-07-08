package main

import (
	"fmt"
	"io/ioutil"
	"path"
)

const SV_PATH string = "/etc/runit/sv"

type Service struct {
	Name    string
	Enabled bool
}

func printServices(services []*Service) {
	for _, i := range services {
		status := "\033[91mDISABLED\033[0m"
		if i.Enabled {
			status = "\033[32mENABLED\033[0m"
		}

		fmt.Printf("\033[1m%s\033[0m\t%s\n", i.Name, status)
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
