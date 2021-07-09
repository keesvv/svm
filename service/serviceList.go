package service

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/keesvv/svm/consts"
	"github.com/keesvv/svm/errs"
)

type ServiceList []*Service

func (list *ServiceList) FindByName(name string) (*Service, error) {
	var sv *Service

	for _, i := range *list {
		if i.Name == name {
			sv = i
			break
		}
	}

	if sv == nil {
		return nil, errs.ErrNoSuchService
	}

	return sv, nil
}

func ListServices() (ServiceList, error) {
	// List service dirs
	svDirs, err := ioutil.ReadDir(consts.SV_PATH)
	if err != nil {
		panic(err)
	}

	services := make(ServiceList, 0)
	for _, i := range svDirs {
		svEnabled := false
		svRunlevel := LEVEL_NONE

		stat, err := ioutil.ReadFile(path.Join(consts.SV_PATH, i.Name(), "supervise", "stat"))

		// User has insufficient permissions
		if os.IsPermission(err) {
			return nil, errs.ErrPermission
		}

		// An unknown error occurred
		if err != nil && !os.IsNotExist(err) {
			return nil, err
		}

		// PID file exists, service is running
		if err == nil && string(stat) == "run\n" {
			svEnabled = true
		}

		// Default & single user runlevel
		defaultRl := path.Join(consts.RUNSVDIR_PATH, "default", i.Name())
		singleRl := path.Join(consts.RUNSVDIR_PATH, "single", i.Name())

		if _, err := os.Stat(defaultRl); !os.IsNotExist(err) {
			svRunlevel = LEVEL_DEFAULT
		}

		if _, err := os.Stat(singleRl); !os.IsNotExist(err) {
			svRunlevel = LEVEL_SINGLE
		}

		services = append(services, &Service{
			Name:     i.Name(),
			Running:  svEnabled,
			Path:     path.Join(consts.SV_PATH, i.Name()),
			Runlevel: svRunlevel,
		})
	}

	return services, nil
}
