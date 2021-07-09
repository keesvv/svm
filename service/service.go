package service

import (
	"os"
	"path"

	"github.com/keesvv/svm/errs"
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
		return err
	}

	// Write command
	f.Write([]byte(cmd))

	return f.Close()
}

func (service *Service) Stop() error {
	if !service.Running {
		return errs.ErrIsStopped
	}

	return service.WriteCommand("d")
}

func (service *Service) Start() error {
	if service.Running {
		return errs.ErrIsStarted
	}

	return service.WriteCommand("u")
}
