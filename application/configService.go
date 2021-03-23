package application

import (
	"code-fabrik.com/bend/domain/config"
	"sort"
)

type ConfigService struct {
	ConfigRepository config.Repository
}

func (cs ConfigService) GetConfigData(path string) config.ConfigData {
	configs := cs.ConfigRepository.FindAll()
	sort.SliceStable(configs, func(i, j int) bool {
		return configs[i].Path > configs[j].Path
	})

	return config.ConfigData{
		Configs:       configs,
		CurrentConfig: cs.getCurrentConfig(path),
	}
}

func (cs ConfigService) getCurrentConfig(path string) config.Config {
	configData := cs.ConfigRepository.Find(path)
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
	return cs.ConfigRepository.Delete(path)
}

func (cs ConfigService) Save(config config.Config) error {
	return cs.ConfigRepository.Save(config)
}
