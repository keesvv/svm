package service

import (
	"fmt"
	"os"
	"path"
)

type Service struct {
	Name    string
	Running bool
	Path    string
}

func (service *Service) WriteCommand(cmd string) error {
	// Open control file for writing
	f, err := os.Create(path.Join(service.Path, "supervise", "control"))

	if err != nil {
		panic(err)
	}

	// Write command
	f.Write([]byte(cmd))

	return f.Close()
}

func (service *Service) Stop() error {
	if !service.Running {
		fmt.Println("service is already stopped")
		os.Exit(1)
	}

	return service.WriteCommand("d")
}

func (service *Service) Start() error {
	if service.Running {
		fmt.Println("service is already running")
		os.Exit(1)
	}

	return service.WriteCommand("u")
}
