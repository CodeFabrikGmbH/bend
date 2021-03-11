package env

import (
	"os"
)

var KEYCLOAK_HOST = getEnvironmentStringValue("KEYCLOAK_HOST", "")
var KEYCLOAK_CLIENT_ID = getEnvironmentStringValue("KEYCLOAK_CLIENT_ID", "")
var KEYCLOAK_REALM = getEnvironmentStringValue("KEYCLOAK_REALM", "")

func getEnvironmentStringValue(name string, defaultValue string) string {
	value := os.Getenv(name)
	if value == "" {
		value = defaultValue
	}
	return value
}
