package httpHandler

import (
	"bytes"
	"code-fabrik.com/bend/application"
	"code-fabrik.com/bend/domain/request"
	"code-fabrik.com/bend/infrastructure/boltDB"
	"code-fabrik.com/bend/infrastructure/jwt/keycloak"
	"encoding/json"
	"github.com/boltdb/bolt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
)

var (
	db, _             = bolt.Open(os.TempDir()+"/test.db", 0600, nil)
	keycloakService   = keycloak.New()
	configRepository  = boltDB.ConfigRepository{DB: db}
	requestRepository = boltDB.RequestRepository{DB: db}
	testTransport     = TestTransport{}

	configService = application.ConfigService{
		ConfigRepository: configRepository,
	}

	requestService = application.RequestService{
		RequestRepository: requestRepository,
		ConfigRepository:  configRepository,
		Transport:         testTransport,
	}

	dashboardService = application.DashboardService{
		RequestRepository: requestRepository,
	}

	configPage = ConfigPage{KeyCloakService: keycloakService, ConfigService: configService}
	configAPI  = ConfigAPI{KeyCloakService: keycloakService, ConfigService: configService}

	dashboardPage = DashboardPage{KeyCloakService: keycloakService, DashboardService: dashboardService}

	requestAPI = RequestAPI{RequestService: requestService}

	tracker = TrackRequest{RequestService: requestService}
)

var (
	defaultTestConfigInput = ConfigInput{
		OriginalPath: "/originalPath",
		Path:         "/path",
		Target:       "",
		StatusCode:   "200",
		Body:         "testConfigBody",
	}

	defaultTestSendRequestInput = SendRequestInput{
		TargetUrl: "sendRequestTargetUrl",
	}
)

var (
	lastTransportRequest   request.Request
	lastTransportTargetUrl string
)

type TestTransport struct {
}

func (t TestTransport) SendRequestToTarget(rr request.Request, targetUrl string) request.Response {
	lastTransportRequest = rr
	lastTransportTargetUrl = targetUrl
	return request.Response{
		Error:              "",
		Target:             targetUrl,
		ResponseBody:       "transportResponse",
		ResponseStatusCode: 200,
	}
}

func before() {
	//set workdir to root, but only once
	workdir, _ := os.Getwd()
	if strings.Index(workdir, "httpHandler") > 0 {
		//workdir is in httpHandler folder
		//set to application root folder
		_ = os.Chdir("../..")
	}

	cleanDB()
}

func cleanDB() {
	_ = db.Update(func(txn *bolt.Tx) error {
		return txn.ForEach(func(name []byte, b *bolt.Bucket) error {
			_ = txn.DeleteBucket(name)
			return nil
		})
	})
}

func runTestServe(req *http.Request, handler http.Handler) (statusCode int, response string) {
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	responseData, _ := ioutil.ReadAll(rr.Body)
	return rr.Result().StatusCode, string(responseData)
}

func putDefaultTestConfig() {
	requestBody, _ := json.Marshal(defaultTestConfigInput)
	req, _ := http.NewRequest("PUT", "/api/configs", bytes.NewBuffer(requestBody))
	_, _ = runTestServe(req, configAPI)
}

func simulateDefaultTrackRequest() {
	simulateTrackRequest("PUT", defaultTestConfigInput.Path, "defaultRequestBody")
}

func simulateTrackRequest(method, path, requestBody string) {
	req, _ := http.NewRequest(method, path, bytes.NewBuffer([]byte(requestBody)))
	_, _ = runTestServe(req, tracker)
}
