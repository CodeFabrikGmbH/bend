package env

import (
	"os"
)

type KeyCloakConfig struct {
	Host     string
	ClientId string
	Realm    string
}

var (
	KeycloakConfig = KeyCloakConfig{
		Host:     getEnvironmentStringValue("KEYCLOAK_HOST", ""),
		ClientId: getEnvironmentStringValue("KEYCLOAK_CLIENT_ID", ""),
		Realm:    getEnvironmentStringValue("KEYCLOAK_REALM", ""),
	}
)

func getEnvironmentStringValue(name string, defaultValue string) string {
	value := os.Getenv(name)
	if value == "" {
		value = defaultValue
	}
	return value
}
