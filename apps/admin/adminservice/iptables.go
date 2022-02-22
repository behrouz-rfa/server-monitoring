package adminservice

import (
	"fmt"
	"os/exec"
	"runtime"
)

var (
	IpTableService IpTableServiceInterface = &ipTableService{}
)

type ipTableService struct {
}

type IpTableServiceInterface interface {
	BlockIP(string) (string, error)
	UnblockIp(string) error
}

// run os command for blokgin ip
//base on OS
func (i ipTableService) BlockIP(s string) (string, error) {

	if runtime.GOOS == "windows" {
		command := fmt.Sprintf("New-NetFirewallRule -DisplayName 'Block host' -Direction Outbound –LocalPort Any -Protocol TCP -Action Block -RemoteAddress  %s", s)
		out, err := exec.Command("powershell", command).Output()
		if err != nil {
			return "", err
		} else {
			fmt.Printf("%s", out)
		}
		return fmt.Sprintf("%s", out), err
	}
	command := fmt.Sprintf("sudo ufw deny from %s to any", s)

	out, err := exec.Command("sh", "-c", command).Output()
	if err != nil {
		return "", err
	}
	fmt.Printf("%s", out)
	return fmt.Sprintf("%s", out), err
}

// unlbloc ip base on OS
func (i ipTableService) UnblockIp(s string) error {
	panic("implement me")
}
