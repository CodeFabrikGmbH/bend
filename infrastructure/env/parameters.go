package env

import (
	"os"
)

var KEYCLOAK_HOST = getEnvironmentStringValue("KEYCLOAK_HOST", "")
var KEYCLOAK_CLIENT_ID = getEnvironmentStringValue("KEYCLOAK_CLIENT_ID", "")
var KEYCLOAK_REALM = getEnvironmentStringValue("KEYCLOAK_REALM", "")

// COOKIE_SECURE controls the Secure flag on session cookies. Defaults to true so
// tokens are only sent over HTTPS; set COOKIE_SECURE=false for local HTTP setups.
var COOKIE_SECURE = getEnvironmentBoolValue("COOKIE_SECURE", true)

// LISTEN_ADDR is the address the HTTP server binds to. Override with BEND_ADDR.
var LISTEN_ADDR = getEnvironmentStringValue("BEND_ADDR", ":8080")

func getEnvironmentStringValue(name string, defaultValue string) string {
	value := os.Getenv(name)
	if value == "" {
		value = defaultValue
	}
	return value
}

func getEnvironmentBoolValue(name string, defaultValue bool) bool {
	value := os.Getenv(name)
	if value == "" {
		return defaultValue
	}
	return value == "true" || value == "1"
}
