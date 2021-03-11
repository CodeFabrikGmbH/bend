package application

import (
	"code-fabrik.com/bend/domain/environment"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type SendRequestService struct {
	Env environment.Environment
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

	req := rs.Env.RequestRepository.GetRequest(path, *requestId)
	response := rs.Env.Transport.SendRequestToTarget(req, *targetUrl)

	marshal, _ := json.Marshal(response)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(marshal)
}
