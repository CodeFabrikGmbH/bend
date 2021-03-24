package httpHandler

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"strings"
	"testing"
)

func Test_ConfigPage(t *testing.T) {
	before()

	request, _ := http.NewRequest("GET", "/configs", nil)
	statusCode, response := runTestServe(request, configPage)

	assert.Equal(t, 200, statusCode)
	assert.True(t, strings.Index(response, "<") == 0)
}
