package httpHandler

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func Test_TrackRequest(t *testing.T) {
	before()
	putDefaultTestConfig()

	req, _ := http.NewRequest("PUT", defaultTestConfigInput.Path, bytes.NewBuffer([]byte("requestBody")))
	statusCode, response := runTestServe(req, tracker)

	assert.Equal(t, 200, statusCode)
	assert.Equal(t, defaultTestConfigInput.Body, response)

	requests := requestRepository.GetRequestsForPath(defaultTestConfigInput.Path)
	assert.Equal(t, 1, len(requests))

	assert.Equal(t, "PUT", requests[0].Method)
	assert.Equal(t, "requestBody", requests[0].Body)
	assert.Equal(t, defaultTestConfigInput.Path, requests[0].Path)
}

func Test_TrackRequest_Multiple(t *testing.T) {
	before()

	simulateDefaultTrackRequest()
	assert.Equal(t, 1, len(requestRepository.GetRequestsForPath(defaultTestConfigInput.Path)))

	simulateDefaultTrackRequest()
	assert.Equal(t, 2, len(requestRepository.GetRequestsForPath(defaultTestConfigInput.Path)))
}
