package application

import (
	"code-fabrik.com/bend/domain/request"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type SendRequestService struct {
	RequestRepository request.Repository
	TransportService  request.Transport
}

func (rs SendRequestService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if rec := recover(); rec != nil {
			fmt.Println(rec)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}()

	path := strings.TrimPrefix(r.URL.Path, "/sendRequest")

	requestId := getQueryValueOrNil(r.URL.Query(), "requestId")
	if requestId == nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("requestId missing"))
		return
	}

	targetUrl := getQueryValueOrNil(r.URL.Query(), "targetUrl")
	if targetUrl == nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("targetUrl missing"))
		return
	}

	req := rs.RequestRepository.GetRequest(path, *requestId)
	response := rs.TransportService.SendRequestToTarget(req, *targetUrl)

	marshal, _ := json.Marshal(response)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(marshal)
}
