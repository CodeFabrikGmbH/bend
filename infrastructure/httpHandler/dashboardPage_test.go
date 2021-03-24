package httpHandler

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"strings"
	"testing"
)

func Test_DashboardPage(t *testing.T) {
	before()

	request, _ := http.NewRequest("GET", "/dashboard", nil)
	statusCode, response := runTestServe(request, dashboardPage)

	assert.Equal(t, 200, statusCode)
	assert.True(t, strings.Index(response, "<") == 0)
}
