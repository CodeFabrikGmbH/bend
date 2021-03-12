package httpHandler

import (
	"code-fabrik.com/bend/application"
	"code-fabrik.com/bend/domain/request"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Request struct {
	RequestService application.RequestService
}

func (rs Request) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if rec := recover(); rec != nil {
			fmt.Println(rec)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}()

	defer func() {
		_ = r.Body.Close()
	}()
	body, _ := ioutil.ReadAll(r.Body)

	req := request.Request{
		Timestamp: time.Now().UnixNano(),
		Path:      r.URL.Path,
		Method:    r.Method,
		Body:      string(body),
		Header:    r.Header,
		Host:      r.Host,
		Uri:       r.RequestURI,
	}

	response := rs.RequestService.HandleTrackableRequest(req)

	w.WriteHeader(response.ResponseStatusCode)
	w.Write([]byte(response.ResponseBody))
}
