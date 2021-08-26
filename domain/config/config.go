package config

import (
	"errors"
	"regexp"
	"strings"
)

type Response struct {
	StatusCode int    `json:"statusCode"`
	Body       string `json:"body"`
}

type Config struct {
	Path     string   `json:"path"`
	Target   string   `json:"target"`
	Response Response `json:"response"`
}

type ConfigData struct {
	Configs       []Config
	CurrentConfig Config
}

func TestIfRegexConfigMatches(configs []Config, path string) (string, error) {
	for _, config := range configs {
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
			return sb.String(), nil
		}
	}
	return "", errors.New("no match")
}
