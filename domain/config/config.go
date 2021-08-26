package config

import (
	"errors"
	"github.com/google/uuid"
	"regexp"
	"strings"
)

type Response struct {
	StatusCode int    `json:"statusCode"`
	Body       string `json:"body"`
}

type Config struct {
	Path     string    `json:"path"`
	Target   string    `json:"target"`
	Response Response  `json:"response"`
	Id       uuid.UUID `json:"id"`
}

type ConfigData struct {
	Configs       []Config
	CurrentConfig Config
}

func FindFirstMatchingConfig(configs []Config, path string) (Config, error) {
	for _, config := range configs {
		if config.Path == path {
			return config, nil
		}
		matched, _ := regexp.MatchString(`^`+config.Path+`$`, path)
		if matched && len(config.Target) != 0 {
			domains := strings.Split(path, `/`)
			var sb strings.Builder
			sb.WriteString(config.Target)
			for _, part := range domains {
				if len(part) > 0 {
					sb.WriteString(`/`)
					sb.WriteString(part)
				}
			}
			config.Target = sb.String()
			return config, nil
		}
	}
	return Config{}, errors.New("no match")
}
