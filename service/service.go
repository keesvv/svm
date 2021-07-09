package service

import (
	"os"
	"path"
	"time"

	"github.com/keesvv/svm/consts"
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

func (service *Service) LastModified() (time.Time, error) {
	f, err := os.Stat(path.Join(service.Path, "supervise", "control"))

	if err != nil {
		return time.Unix(0, 0), err
	}

	return f.ModTime(), nil
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

func (service *Service) SetRunlevel(rl Runlevel) error {
	targetRl := path.Join(consts.RUNSVDIR_PATH, string(rl), service.Name)
	defaultRl := path.Join(consts.RUNSVDIR_PATH, "default", service.Name)

	if rl == LEVEL_NONE {
		if _, err := os.Stat(defaultRl); os.IsNotExist(err) {
			return errs.ErrAlreadyDisabled
		}

		return os.Remove(defaultRl)
	}

	if _, err := os.Stat(targetRl); !os.IsNotExist(err) {
		return errs.ErrRunlevelExists
	}

	return os.Symlink(
		service.Path,
		targetRl,
	)
}
