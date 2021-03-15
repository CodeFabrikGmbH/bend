package application

import (
	"code-fabrik.com/bend/domain/config"
	"code-fabrik.com/bend/domain/environment"
	"sort"
)

type ConfigService struct {
	Env environment.Environment
}

func (cs ConfigService) GetConfigData(path string) config.ConfigData {
	configs := cs.Env.ConfigRepository.FindAll()
	sort.SliceStable(configs, func(i, j int) bool {
		return configs[i].Path > configs[j].Path
	})

	return config.ConfigData{
		Configs:       configs,
		CurrentConfig: cs.getCurrentConfig(path),
	}
}

func (cs ConfigService) getCurrentConfig(path string) config.Config {
	configData := cs.Env.ConfigRepository.Find(path)
	if configData != nil {
		return *configData
	}
	return config.Config{
		Path:   "",
		Target: "",
		Response: config.Response{
			StatusCode: 200,
			Body:       "ok",
		},
	}
}

func (cs ConfigService) Delete(path string) error {
	return cs.Env.ConfigRepository.Delete(path)
}

func (cs ConfigService) Save(config config.Config) error {
	return cs.Env.ConfigRepository.Save(config)
}
