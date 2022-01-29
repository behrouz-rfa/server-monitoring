package adminservice

import (
	"os/exec"
	"runtime"
)

var (
	HomeServices homeServicesInterface = &homeService{}
)

type homeServicesInterface interface {
	GetSsh() (string, error)
	GetUserList() (string, error)
}
type homeService struct {
}

func (h homeService) GetUserList() (string, error) {
	if runtime.GOOS == "windows" {
		out, err := exec.Command("powershell","net","user").Output()
		if err != nil {
			return "",err
		} else {
			return string(out[:]), nil
		}
	}else {
		out, err := exec.Command("awk -F: '{ print $1}' /etc/passwd").Output()
		if err != nil {
			return "", nil
		}

		return string(out[:]), nil
	}
}

func (h homeService) GetSsh() (string, error) {
	out, err := exec.Command("ss | grep -i ssh").Output()
	if err != nil {
		return "", nil
	}

	return string(out[:]), nil
}
