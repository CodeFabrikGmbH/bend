package application

import (
	"code-fabrik.com/bend/domain/config"
	"code-fabrik.com/bend/domain/environment"
)

type ConfigService struct {
	Env environment.Environment
}

func (cs ConfigService) GetConfigData(path string) config.ConfigData {
	return config.ConfigData{}
}

func (cs ConfigService) Delete(path string) error {
	return nil
}
