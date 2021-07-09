package service

import (
	"os"
	"path"

	"github.com/keesvv/svm/errs"
)

type Runlevel string

type Service struct {
	Name     string
	Running  bool
	Path     string
	Runlevel Runlevel
}

const (
	LEVEL_DEFAULT Runlevel = "default"
	LEVEL_SINGLE  Runlevel = "single"
	LEVEL_NONE    Runlevel = "none"
)

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

	if service.Runlevel == LEVEL_NONE {
		return errs.ErrNoRunlevel
	}

	return service.WriteCommand("d")
}

func (service *Service) Start() error {
	if service.Running {
		return errs.ErrIsStarted
	}

	if service.Runlevel == LEVEL_NONE {
		return errs.ErrNoRunlevel
	}

	return service.WriteCommand("u")
}
