package application

import (
	"code-fabrik.com/bend/domain/config"
	"github.com/google/uuid"
	"sort"
)

type ConfigService struct {
	ConfigRepository config.Repository
}

func (cs ConfigService) GetConfigData(id uuid.UUID) config.ConfigData {
	configs := cs.ConfigRepository.FindAll()
	sort.SliceStable(configs, func(i, j int) bool {
		return configs[i].Path > configs[j].Path
	})

	return config.ConfigData{
		Configs:       configs,
		CurrentConfig: cs.getCurrentConfig(id),
	}
}

func (cs ConfigService) getCurrentConfig(id uuid.UUID) config.Config {
	configData := cs.ConfigRepository.Find(id)
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

func (cs ConfigService) Delete(id uuid.UUID) error {
	return cs.ConfigRepository.Delete(id)
}

func (cs ConfigService) Save(config config.Config) error {
	return cs.ConfigRepository.Save(config)
}
