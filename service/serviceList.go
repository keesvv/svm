package service

import (
	"fmt"
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

func ListServices() ServiceList {
	// List service dirs
	svDirs, err := ioutil.ReadDir(consts.SV_PATH)
	if err != nil {
		panic(err)
	}

	services := make([]*Service, 0)
	for _, i := range svDirs {
		svEnabled := false
		stat, err := ioutil.ReadFile(path.Join(consts.SV_PATH, i.Name(), "supervise", "stat"))

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
		if err == nil && string(stat) == "run\n" {
			svEnabled = true
		}

		services = append(services, &Service{
			Name:    i.Name(),
			Running: svEnabled,
			Path:    path.Join(consts.SV_PATH, i.Name()),
		})
	}

	return services
}
