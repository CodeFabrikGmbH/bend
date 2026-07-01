package httpHandler

import (
	"code-fabrik.com/bend/application"
	"code-fabrik.com/bend/infrastructure/jwt/keycloak"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
)

type SendRequestInput struct {
	TargetUrl string `json:"targetUrl"`
}

type RequestAPI struct {
	KeyCloakService *keycloak.Service
	RequestService  application.RequestService
}

func (rs RequestAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if rec := recover(); rec != nil {
			slog.Error("panic in RequestAPI", "recover", rec)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}()

	if _, err := rs.KeyCloakService.Authenticate(w, r); err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/requests")

	requestPath, requestId := getRequestPathAndId(path)

	var err error
	var response string

	switch r.Method {
	case http.MethodPost:
		body, err := readRequestBody(r)
		if err == nil {
			requestBody := SendRequestInput{}
			err = json.Unmarshal(body, &requestBody)

			if err == nil {
				requestResponse := rs.RequestService.SendRequestToTarget(requestPath, requestId, requestBody.TargetUrl)

				marshal, _ := json.Marshal(requestResponse)
				response = string(marshal)
			}
		}
	case http.MethodDelete:
		if requestId == "*" {
			err = rs.RequestService.DeleteAllRequestsForPath(requestPath)
		} else {
			err = rs.RequestService.DeleteRequest(requestPath, requestId)
		}
	default:
		err = fmt.Errorf("method %s not implemented", r.Method)
	}
	writeResponse(w, response, err)
}

func getRequestPathAndId(path string) (requestPath string, requestId string) {
	i := strings.LastIndex(path, "/")
	if i == -1 {
		return
	}

	requestPath = path[0:i]
	requestId = path[i+1:]
	return
}
