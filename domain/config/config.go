package config

import (
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

func (c *Config) HasTargetPath() bool {
	return len(c.Target) != 0
}

func (c *Config) GenerateFinalTargetPath(path string) string {
	if c.Path == path {
		return c.Target
	}

	var sb strings.Builder
	sb.WriteString(c.Target)
	domains := strings.Split(path, `/`)
	for _, part := range domains {
		if len(part) > 0 {
			sb.WriteString(`/`)
			sb.WriteString(part)
		}
	}
	return sb.String()
}

type ConfigData struct {
	Configs       []Config
	CurrentConfig Config
}

func FindMatchingConfig(configs []Config, path string) *Config {
	for _, conf := range configs {
		if conf.Path == path {
			return &conf
		}
		matched, _ := regexp.MatchString(`^`+conf.Path+`$`, path)
		if matched {
			return &conf
		}
	}
	return nil
}
