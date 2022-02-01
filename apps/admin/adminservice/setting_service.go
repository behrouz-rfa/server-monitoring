package adminservice

import (
	"server-monitoring/domain/settings"
)

var (
	SettingService settingServiceInterface = &settingService{}
)

type settingServiceInterface interface {
	Get(*settings.Setting) error
	Login(*settings.Setting) error
	Create(*settings.Setting) error
}
type settingService struct {
}

func (s settingService) Login(setting *settings.Setting) error {
	return setting.FindByUserName()
}

func (s settingService) Get(setting *settings.Setting) error {
	return setting.FindFist()
}

func (s settingService) Create(setting *settings.Setting) error {
	return setting.Update()
}
