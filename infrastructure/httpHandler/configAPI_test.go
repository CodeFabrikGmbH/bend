package httpHandler

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"strconv"
	"testing"
)

func Test_ConfigAPI_PUT(t *testing.T) {
	before()
	requestBody, _ := json.Marshal(defaultTestConfigInput)

	req, _ := http.NewRequest("PUT", "/api/configs", bytes.NewBuffer(requestBody))
	statusCode, response := runTestServe(req, configAPI)

	assert.Equal(t, 200, statusCode)
	assert.Equal(t, "ok", response)

	configs := configRepository.FindAll()

	statusCode, _ = strconv.Atoi(defaultTestConfigInput.StatusCode)
	assert.True(t, len(configs) == 1)
	assert.Equal(t, configs[0].Path, defaultTestConfigInput.Path)
	assert.Equal(t, configs[0].Target, defaultTestConfigInput.Target)
	assert.Equal(t, configs[0].Response.StatusCode, statusCode)
	assert.Equal(t, configs[0].Response.Body, defaultTestConfigInput.Body)
}

func Test_ConfigAPI_DELETE(t *testing.T) {
	before()
	putDefaultTestConfig()

	assert.True(t, len(configRepository.FindAll()) == 1)

	req, _ := http.NewRequest("DELETE", "/api/configs"+defaultTestConfigInput.Path, nil)
	_, _ = runTestServe(req, configAPI)

	assert.True(t, len(configRepository.FindAll()) == 0)
}
