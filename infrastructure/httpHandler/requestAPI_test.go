package httpHandler

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func Test_RequestAPI_DELETE_Requests(t *testing.T) {
	before()
	simulateDefaultTrackRequest()
	simulateDefaultTrackRequest()

	requestsForPath := requestRepository.GetRequestsForPath(defaultTestConfigInput.Path)
	assert.Equal(t, 2, len(requestsForPath))

	req, _ := http.NewRequest("DELETE", "/api/requests"+defaultTestConfigInput.Path+"/"+requestsForPath[0].ID, nil)
	statusCode, _ := runTestServe(req, requestAPI)

	assert.Equal(t, 200, statusCode)
	assert.Equal(t, 1, len(requestRepository.GetRequestsForPath(defaultTestConfigInput.Path)))
}

func Test_RequestAPI_DELETE_All_Requests_For_Path(t *testing.T) {
	before()
	simulateDefaultTrackRequest()
	simulateDefaultTrackRequest()

	assert.Equal(t, 2, len(requestRepository.GetRequestsForPath(defaultTestConfigInput.Path)))

	req, _ := http.NewRequest("DELETE", "/api/requests"+defaultTestConfigInput.Path+"/*", nil)
	statusCode, _ := runTestServe(req, requestAPI)

	assert.Equal(t, 200, statusCode)
	assert.Equal(t, 0, len(requestRepository.GetRequestsForPath(defaultTestConfigInput.Path)))
}

func Test_RequestAPI_POST_Send_Request(t *testing.T) {
	before()
	simulateTrackRequest("PUT", defaultTestConfigInput.Path, "myTestBody")

	requestBody, _ := json.Marshal(defaultTestSendRequestInput)

	requestsForPath := requestRepository.GetRequestsForPath(defaultTestConfigInput.Path)
	req, _ := http.NewRequest("POST", "/api/requests"+defaultTestConfigInput.Path+"/"+requestsForPath[0].ID, bytes.NewBuffer(requestBody))
	statusCode, _ := runTestServe(req, requestAPI)

	assert.Equal(t, 200, statusCode)
	assert.Equal(t, defaultTestSendRequestInput.TargetUrl, lastTransportTargetUrl)
	assert.Equal(t, "PUT", lastTransportRequest.Method)
	assert.Equal(t, defaultTestConfigInput.Path, lastTransportRequest.Path)
	assert.Equal(t, "myTestBody", lastTransportRequest.Body)
}
